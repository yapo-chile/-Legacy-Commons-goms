package domain

// WayStair contains struct to be used as output delivering
// possible N ways and its combinations
type WayStair struct {
	Ways  int
	Combs string
}

// WayStairRepository defines function that do the maths
type WayStairRepository interface {
	Calculate(nth int) (WayStair, error)
}
