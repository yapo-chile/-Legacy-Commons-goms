package usecases

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/mock"
	"testing"
)

type MockAdder struct {
	mock.Mock
}

func (m *MockAdder) Add(a, b int) int {
	args := m.Called(a, b)
	return args.Int(0)
}

func TestCalculatorExecuteAdd(t *testing.T) {
	m := MockAdder{}
	m.On("Add", 5, 7).Return(42, nil)
	c := ModularCalculator{Adder: &m}
	r, err := c.Execute("add", 5, 7)
	assert.Equal(t, r, 42)
	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestCalculatorExecuteOther(t *testing.T) {
	m := MockAdder{}
	c := ModularCalculator{Adder: &m}
	r, err := c.Execute("other", 5, 7)
	assert.Equal(t, r, -1)
	assert.Error(t, err)
	m.AssertExpectations(t)
}
