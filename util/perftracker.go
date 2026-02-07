package util

import (
	"fmt"
	"snakehem/model"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/rs/zerolog/log"
)

const (
	// WindowSize is the number of TPS/FPS samples to keep for rolling averages
	windowSize = 5 * model.Tps // 5 seconds at 60 TPS

	// Histogram parameters
	minMicros = 1       // 1 microsecond minimum
	maxMicros = 1000000 // 1 second maximum
	sigFigs   = 3       // 3 significant figures

	// Growth detection parameters
	percentileHistorySize = 1 * model.Tps // 1 second of percentile snapshots
	growthThreshold       = 1.5           // 50% increase triggers warning
)

type percentileSnapshot struct {
	p50 int64
	p90 int64
	p95 int64
	p99 int64
}

type PerfTracker struct {
	updateHist           *hdrhistogram.Histogram
	drawHist             *hdrhistogram.Histogram
	tpsSamples           []float64
	fpsSamples           []float64
	updateBaselineP      percentileSnapshot
	drawBaselineP        percentileSnapshot
	ticksSinceLastUpdate int
}

func NewPerfTracker() *PerfTracker {
	return &PerfTracker{
		updateHist: hdrhistogram.New(minMicros, maxMicros, sigFigs),
		drawHist:   hdrhistogram.New(minMicros, maxMicros, sigFigs),
		tpsSamples: make([]float64, 0, windowSize),
		fpsSamples: make([]float64, 0, windowSize),
	}
}

func (t *PerfTracker) RecordUpdate(duration time.Duration) {
	if err := t.updateHist.RecordValue(duration.Microseconds()); err != nil {
		log.Warn().Dur("duration", duration).Msg("Update took incredibly long")
	}
}

func (t *PerfTracker) RecordDraw(duration time.Duration) {
	if err := t.drawHist.RecordValue(duration.Microseconds()); err != nil {
		log.Warn().Dur("duration", duration).Msg("Draw took incredibly long")
	}
}

func (t *PerfTracker) RecordTPS(tps float64) {
	t.tpsSamples = append(t.tpsSamples, tps)
	if len(t.tpsSamples) > windowSize {
		t.tpsSamples = t.tpsSamples[1:]
	}
}

func (t *PerfTracker) RecordFPS(fps float64) {
	t.fpsSamples = append(t.fpsSamples, fps)
	if len(t.fpsSamples) > windowSize {
		t.fpsSamples = t.fpsSamples[1:]
	}
}

type PerfStats struct {
	UpdateP50     time.Duration
	UpdateP90     time.Duration
	UpdateP95     time.Duration
	UpdateP99     time.Duration
	DrawP50       time.Duration
	DrawP90       time.Duration
	DrawP95       time.Duration
	DrawP99       time.Duration
	TPSAvg        float64
	FPSAvg        float64
	SampleCount   int64
	UpdateWarning bool // True if Update percentiles are growing too fast
	DrawWarning   bool // True if Draw percentiles are growing too fast
}

//goland:noinspection DuplicatedCode
func (t *PerfTracker) GetStats() PerfStats {
	stats := PerfStats{
		SampleCount: t.updateHist.TotalCount(),
	}

	var currentUpdateSnapshot, currentDrawSnapshot percentileSnapshot

	if t.updateHist.TotalCount() > 0 {
		stats.UpdateP50 = time.Duration(t.updateHist.Mean()) * time.Microsecond
		stats.UpdateP90 = time.Duration(t.updateHist.ValueAtQuantile(90)) * time.Microsecond
		stats.UpdateP95 = time.Duration(t.updateHist.ValueAtQuantile(95)) * time.Microsecond
		stats.UpdateP99 = time.Duration(t.updateHist.ValueAtQuantile(99)) * time.Microsecond

		currentUpdateSnapshot = percentileSnapshot{
			p50: int64(t.updateHist.Mean()),
			p90: t.updateHist.ValueAtQuantile(90),
			p95: t.updateHist.ValueAtQuantile(95),
			p99: t.updateHist.ValueAtQuantile(99),
		}
	}

	if t.drawHist.TotalCount() > 0 {
		stats.DrawP50 = time.Duration(t.drawHist.Mean()) * time.Microsecond
		stats.DrawP90 = time.Duration(t.drawHist.ValueAtQuantile(90)) * time.Microsecond
		stats.DrawP95 = time.Duration(t.drawHist.ValueAtQuantile(95)) * time.Microsecond
		stats.DrawP99 = time.Duration(t.drawHist.ValueAtQuantile(99)) * time.Microsecond

		currentDrawSnapshot = percentileSnapshot{
			p50: int64(t.drawHist.Mean()),
			p90: t.drawHist.ValueAtQuantile(90),
			p95: t.drawHist.ValueAtQuantile(95),
			p99: t.drawHist.ValueAtQuantile(99),
		}
	}

	if len(t.tpsSamples) > 0 {
		stats.TPSAvg = avgFloat(t.tpsSamples)
	}

	if len(t.fpsSamples) > 0 {
		stats.FPSAvg = avgFloat(t.fpsSamples)
	}

	// Detect rapid growth in percentiles (compare against baseline from 1 second ago)
	stats.UpdateWarning = t.detectGrowth(t.updateBaselineP, currentUpdateSnapshot)
	stats.DrawWarning = t.detectGrowth(t.drawBaselineP, currentDrawSnapshot)

	// Update baseline snapshot every percentileHistorySize ticks (1 second)
	t.ticksSinceLastUpdate++
	if t.ticksSinceLastUpdate >= percentileHistorySize {
		t.updateBaselineP = currentUpdateSnapshot
		t.drawBaselineP = currentDrawSnapshot
		t.ticksSinceLastUpdate = 0
	}

	return stats
}

func (s PerfStats) AsString() []string {
	return []string{
		fmt.Sprintf("TPS - avg: %.0f, FPS - avg: %.0f", s.TPSAvg, s.FPSAvg),
		fmt.Sprintf("Update - P50: %v, P90: %v, P95: %v, P99: %v",
			formatDuration(s.UpdateP50),
			formatDuration(s.UpdateP90),
			formatDuration(s.UpdateP95),
			formatDuration(s.UpdateP99)),
		fmt.Sprintf("Draw - P50: %v, P90: %v, P95: %v, P99: %v",
			formatDuration(s.DrawP50),
			formatDuration(s.DrawP90),
			formatDuration(s.DrawP95),
			formatDuration(s.DrawP99)),
	}
}

func (t *PerfTracker) detectGrowth(baseline, current percentileSnapshot) bool {
	// Need a valid baseline (all zeros means we haven't collected 1 second of data yet)
	if baseline.p50 == 0 && baseline.p90 == 0 && baseline.p95 == 0 && baseline.p99 == 0 {
		return false
	}

	// Check if any percentile has grown too fast
	return hasExcessiveGrowth(baseline.p50, current.p50) ||
		hasExcessiveGrowth(baseline.p90, current.p90) ||
		hasExcessiveGrowth(baseline.p95, current.p95) ||
		hasExcessiveGrowth(baseline.p99, current.p99)
}

func hasExcessiveGrowth(baseline, current int64) bool {
	if baseline == 0 {
		return false
	}
	ratio := float64(current) / float64(baseline)
	return ratio > growthThreshold
}

func avgFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0"
	}
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fμs", float64(d.Nanoseconds())/1000.0)
	}
	return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1000000.0)
}
