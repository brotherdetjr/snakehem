package scoreboard

import (
	"image/color"
	"slices"
)

type Entry struct {
	Name       string
	Score      int
	ColourFunc func() color.Color
}

type Scoreboard struct {
	entries []Entry
}

func NewScoreboard(entries []Entry) *Scoreboard {
	sortedEntries := make([]Entry, len(entries))
	copy(sortedEntries, entries)
	slices.SortFunc(sortedEntries, func(a, b Entry) int {
		return b.Score - a.Score
	})
	return &Scoreboard{
		entries: sortedEntries,
	}
}
