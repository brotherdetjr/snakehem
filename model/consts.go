package model

const (
	GameSpeedFps               = 10
	TpsMultiplier              = 6
	Tps                        = GameSpeedFps * TpsMultiplier
	GridSize                   = 63
	ControllerRepeatIntervalHz = 5
	ControllerRepeatPeriod     = Tps / ControllerRepeatIntervalHz

	MaxNameLength                 = 9
	CountdownSeconds              = 4
	SnakeTargetLength             = 50
	HealthReductionPerBite        = 10
	NippedTailLinkBonusMultiplier = 2
	BitLinkScore                  = 1
	AppleScore                    = 45
	TargetScore                   = 999
	MaxSnakes                     = 9
	ApproachingTargetScoreGap     = SnakeTargetLength*NippedTailLinkBonusMultiplier - 1
	GridFadeCountdown             = TpsMultiplier * 15
	NewAppleProbabilityParam      = Tps * 3
)
