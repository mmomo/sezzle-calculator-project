// Package calculator implements the calculator's arithmetic. It has no
// knowledge of HTTP; callers are responsible for structural validation
// (unknown operations, missing operands) before calling Calculate.
package calculator

import (
	"errors"
	"math"
)

// Operation identifies a supported calculation.
type Operation string

const (
	Add          Operation = "add"
	Subtract     Operation = "subtract"
	Multiply     Operation = "multiply"
	Divide       Operation = "divide"
	Exponentiate Operation = "exponentiate"
	Sqrt         Operation = "sqrt"
	Percentage   Operation = "percentage"
)

var (
	ErrDivideByZero     = errors.New("division by zero")
	ErrNegativeSqrt     = errors.New("cannot take square root of negative number")
	ErrNotFinite        = errors.New("result is not a finite number")
	ErrUnknownOperation = errors.New("unknown operation")
)

// IsValid reports whether op is one of the supported operations.
func IsValid(op Operation) bool {
	switch op {
	case Add, Subtract, Multiply, Divide, Exponentiate, Sqrt, Percentage:
		return true
	default:
		return false
	}
}

// IsUnary reports whether op takes a single operand (a only).
func IsUnary(op Operation) bool {
	return op == Sqrt
}

// Calculate applies op to a and b. b is ignored for unary operations.
//
// Percentage is defined as "a percent of b": (a / 100) * b.
func Calculate(op Operation, a, b float64) (float64, error) {
	var result float64

	switch op {
	case Add:
		result = a + b
	case Subtract:
		result = a - b
	case Multiply:
		result = a * b
	case Divide:
		if b == 0 {
			return 0, ErrDivideByZero
		}
		result = a / b
	case Exponentiate:
		result = math.Pow(a, b)
	case Sqrt:
		if a < 0 {
			return 0, ErrNegativeSqrt
		}
		result = math.Sqrt(a)
	case Percentage:
		result = (a / 100) * b
	default:
		return 0, ErrUnknownOperation
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrNotFinite
	}

	return result, nil
}
