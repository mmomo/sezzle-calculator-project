package calculator

import (
	"errors"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name    string
		op      Operation
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "add", op: Add, a: 2, b: 3, want: 5},
		{name: "subtract", op: Subtract, a: 5, b: 3, want: 2},
		{name: "multiply", op: Multiply, a: 4, b: 5, want: 20},
		{name: "divide", op: Divide, a: 10, b: 2, want: 5},
		{name: "divide by zero", op: Divide, a: 10, b: 0, wantErr: ErrDivideByZero},
		{name: "exponentiate", op: Exponentiate, a: 2, b: 10, want: 1024},
		{name: "exponentiate negative exponent", op: Exponentiate, a: 2, b: -1, want: 0.5},
		{name: "exponentiate overflow", op: Exponentiate, a: 10, b: 400, wantErr: ErrNotFinite},
		{name: "sqrt", op: Sqrt, a: 16, want: 4},
		{name: "sqrt negative", op: Sqrt, a: -4, wantErr: ErrNegativeSqrt},
		{name: "percentage", op: Percentage, a: 50, b: 200, want: 100},
		{name: "percentage of zero", op: Percentage, a: 50, b: 0, want: 0},
		{name: "percentage negative operand", op: Percentage, a: -10, b: 50, want: -5},
		{name: "unknown operation", op: Operation("bogus"), a: 1, b: 1, wantErr: ErrUnknownOperation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.op, tt.a, tt.b)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Calculate() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Calculate() unexpected error = %v", err)
			}

			if got != tt.want {
				t.Fatalf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		op   Operation
		want bool
	}{
		{Add, true},
		{Subtract, true},
		{Multiply, true},
		{Divide, true},
		{Exponentiate, true},
		{Sqrt, true},
		{Percentage, true},
		{Operation("bogus"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.op), func(t *testing.T) {
			if got := IsValid(tt.op); got != tt.want {
				t.Fatalf("IsValid(%q) = %v, want %v", tt.op, got, tt.want)
			}
		})
	}
}

func TestIsUnary(t *testing.T) {
	if !IsUnary(Sqrt) {
		t.Fatal("IsUnary(Sqrt) = false, want true")
	}
	if IsUnary(Add) {
		t.Fatal("IsUnary(Add) = true, want false")
	}
}
