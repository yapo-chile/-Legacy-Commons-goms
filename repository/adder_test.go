package repository

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"testing"
)

func TestRegularAdderAdd(t *testing.T) {
	a, b := 5, 7
	adder := RegularAdder{}
	assert.Equal(t, a+b, adder.Add(a, b))
}

func TestModuloAdderAdd(t *testing.T) {
	a, b, m := 5, 7, 5
	adder := ModuloAdder{M: m}
	assert.Equal(t, (a+b)%m, adder.Add(a, b))
}
