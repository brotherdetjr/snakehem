package game

import (
	"errors"
	"fmt"
	"math/rand"

	"snakehem/internal/config"
	"snakehem/internal/engine"
	"snakehem/internal/entities"
	"snakehem/internal/gamestate"
	"snakehem/internal/input"
	"snakehem/internal/interfaces"
	"snakehem/internal/rendering"
)

// GameBuilder builds a fully-wired Game instance using dependency injection
type GameBuilder struct {
	config        *config.GameConfig
	inputProvider interfaces.InputProvider
	renderer      *rendering.CompositeRenderer
	randomSource  interfaces.RandomSource
}

// NewGameBuilder creates a new game builder with default dependencies
func NewGameBuilder() *GameBuilder {
	return &GameBuilder{
		config:        config.DefaultConfig(),
		inputProvider: input.NewEbitenInputProvider(),
		randomSource:  &DefaultRandomSource{},
	}
}

// WithConfig sets a custom game configuration
func (b *GameBuilder) WithConfig(cfg *config.GameConfig) *GameBuilder {
	b.config = cfg
	return b
}

// WithInputProvider sets a custom input provider
func (b *GameBuilder) WithInputProvider(provider interfaces.InputProvider) *GameBuilder {
	b.inputProvider = provider
	return b
}

// WithRandomSource sets a custom random source (useful for testing)
func (b *GameBuilder) WithRandomSource(random interfaces.RandomSource) *GameBuilder {
	b.randomSource = random
	return b
}

// Build creates and returns a fully-wired Game instance
// Returns error if any dependencies are missing or invalid
func (b *GameBuilder) Build() (*Game, error) {
	// Validate dependencies
	if b.config == nil {
		return nil, errors.New("game config is required")
	}
	if b.inputProvider == nil {
		return nil, errors.New("input provider is required")
	}
	if b.randomSource == nil {
		return nil, errors.New("random source is required")
	}

	// Build rendering components
	fontManager := rendering.NewFontManager()
	snakeRenderer := rendering.NewSnakeRenderer(b.config)
	uiRenderer := rendering.NewUIRenderer(b.config, fontManager)

	postProcessor, err := rendering.NewPostProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to create post processor: %w", err)
	}

	renderer := rendering.NewCompositeRenderer(
		b.config,
		fontManager,
		snakeRenderer,
		uiRenderer,
		postProcessor,
	)

	// Build game logic engines
	physics := engine.NewPhysicsEngine(b.config)
	scoring := engine.NewScoringEngine(b.config)

	// Create the game instance
	game := &Game{
		config:        b.config,
		inputProvider: b.inputProvider,
		renderer:      renderer,
		random:        b.randomSource,
		physics:       physics,
		scoring:       scoring,

		// Initialize mutable state
		grid:          entities.NewGameGrid(b.config.GridSize),
		snakes:        make([]*entities.Snake, 0),
		countdown:     b.config.Tps() * b.config.CountdownSeconds,
		elapsedFrames: 0,
		fadeCountdown: 0,
		applePresent:  false,
		currentState:  &gamestate.LobbyState{},
	}

	return game, nil
}

// DefaultRandomSource implements interfaces.RandomSource using math/rand
type DefaultRandomSource struct{}

// IntN returns a random integer in [0, n)
func (r *DefaultRandomSource) IntN(n int) int {
	return rand.Intn(n)
}
