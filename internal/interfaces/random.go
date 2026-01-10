package interfaces

// RandomSource abstracts randomness for testability
// This allows us to inject deterministic random values in tests
type RandomSource interface {
	IntN(n int) int
}
