package domain

// Fibonacci is the datatype for numbers that belong to the Fibonacci Series
type Fibonacci int

// FibonacciPair is the pair of the latest pair of known Fibonacci Numbers
type FibonacciPair struct {
	IA, IB int       // Indexes
	A, B   Fibonacci // Values
}

// FibonacciRepository defines a backing storage for Fibonacci Numbers
type FibonacciRepository interface {
	Get(nth int) (Fibonacci, error)  // Retrieve the Nth Fibonacci if available
	Save(nth int, x Fibonacci) error // Save the Nth Fibonacci only if (N-1) is known
	LatestPair() FibonacciPair       // Retrieves the latest know pair of Fibonacci
}

// Next produces the Fibonacci Number and the index that comes right after f
func (f FibonacciPair) Next() (int, Fibonacci) {
	return f.IB + 1, f.A + f.B
}
