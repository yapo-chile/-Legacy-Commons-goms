package usecases

import (
	"fmt"
)

// Adder is an interface for a component that knows how to perform integer
// addition
type Adder interface {
	Add(a, b int) int
}

// Calculator is an interface for an entity that can execute math operations
type Calculator interface {
	Execute(op string, a, b int) (int, error)
}

// ModularCalculator is a Calculator implementation that can execute math
// operations based on the provided components
type ModularCalculator struct {
	Adder Adder `inject:""`
}

// Execute runs a math operation on integer domain
func (c *ModularCalculator) Execute(op string, a, b int) (int, error) {
	if op == "add" {
		return c.Adder.Add(a, b), nil
	}
	return -1, fmt.Errorf("Unsupported operation: %s", op)
}
