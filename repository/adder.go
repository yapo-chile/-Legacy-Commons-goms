package repository

type RegularAdder struct{}
type ModuloAdder struct{ M int }

func (*RegularAdder) Add(a, b int) int {
	return a + b
}

func (m *ModuloAdder) Add(a, b int) int {
	return ((a % m.M) + (b % m.M)) % m.M
}
