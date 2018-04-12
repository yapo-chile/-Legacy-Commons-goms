package domain

type Fibonacci int

type FibonacciPair struct {
	IA, IB int
	A, B   Fibonacci
}

type FibonacciRepository interface {
	Get(nth int) (Fibonacci, error)
	Save(nth int, x Fibonacci) error
	LatestPair() FibonacciPair
}

func (f FibonacciPair) Next() (int, Fibonacci) {
	return f.IB + 1, f.A + f.B
}
