package repository

// RegularAdder implements a 32bit integer adder
type RegularAdder struct{}

// ModuloAdder implements modulo-m 32bit integer adder
type ModuloAdder struct{ M int }

// Add adds two 32 bit integers
func (*RegularAdder) Add(a, b int) int {
	return a + b
}

// Add ads two 32 bit integers using modulo m arithmetic
func (m *ModuloAdder) Add(a, b int) int {
	return ((a % m.M) + (b % m.M)) % m.M
}
