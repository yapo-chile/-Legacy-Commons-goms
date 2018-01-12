package usecases

import (
	"fmt"
)

type Adder interface {
	Add(a, b int) int
}
type Calculator interface {
	Execute(op string, a, b int) (int, error)
}

type ModularCalculator struct {
	Adder Adder `inject:""`
}

func (c *ModularCalculator) Execute(op string, a, b int) (int, error) {
	if op == "add" {
		return c.Adder.Add(a, b), nil
	}
	return -1, fmt.Errorf("Unsupported operation: %s", op)
}
