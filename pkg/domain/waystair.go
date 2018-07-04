package domain

// WayStair contains struct to be used as output delivering
// possible N ways and its combinations
type Ways int
type Stair int
type Combs string

type WayStair struct {
	Stair Stair
	Ways  Ways
	Combs Combs
}

// WayStairRepository defines function that do the maths
type WayStairRepository interface {
	Get(nth int) (WayStair, error)
	Save(WayStair) error
}
