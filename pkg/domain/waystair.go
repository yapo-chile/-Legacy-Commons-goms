package domain

// Stair ...
type Stair int

// WayStair ...
type WayStair struct {
	Ways  int
	Combs string
}

// WayStairRepository ...
type WayStairRepository interface {
	WayStairs(nth int) (WayStair, error)
}
