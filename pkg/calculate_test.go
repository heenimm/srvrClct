package pkg_test

import (
	"errors"
	"serverCalc/pkg"
	"testing"
)

func TestCalculate_ValidExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"simple addition", "2+2", 4},
		{"mixed operations", "2+2*2", 6},
		{"parentheses", "(2+2)*2", 8},
		{"division", "6/2", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkg.Calculate(tt.expression)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCalculate_InvalidExpressions(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectedErr string
	}{
		{"insufficient operands", "+2", "insufficient operands"},
		{"division by zero", "6/0", "division by zero"},
		{"invalid expression", "", "internal error"},
		{"missing operators", "123", "internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pkg.Calculate(tt.expression)
			if err == nil {
				t.Fatalf("expected error, got none")
			}
			if !errors.Is(err, pkg.ErrInternalError) && err.Error() != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestCalculate_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"large numbers", "1000000*1000000", 1e+12},
		{"nested parentheses", "((2+3)*2)/5", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkg.Calculate(tt.expression)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
