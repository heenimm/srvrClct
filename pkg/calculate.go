package pkg

import (
	"strconv"
	"strings"
	"unicode"
)

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func operation(a, b float64, op rune) (float64, error) {
	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	}
	return 0, ErrInvalidOperation
}

func Calculate(expression string) (float64, error) {
	var values []float64
	var ops []rune
	i := 0

	// Возвращаем ошибку 500
	if !strings.ContainsAny(expression, "-*/+") {
		return 0, ErrInternalError
	}

	for i < len(expression) {
		char := expression[i]

		if char == ' ' {
			i++
			continue
		}

		if unicode.IsDigit(rune(char)) {
			var sb strings.Builder
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				sb.WriteByte(expression[i])
				i++
			}
			value, err := strconv.ParseFloat(sb.String(), 64)
			if err != nil {
				return 0, ErrInvalidNumber
			}
			values = append(values, value)
			continue
		}

		if char == '(' {
			ops = append(ops, '(')
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if len(values) < 2 {
					return 0, ErrInsufficientOperands
				}
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				res, err := operation(val1, val2, op)
				if err != nil {
					return 0, err
				}
				values = append(values, res)
			}

			if len(ops) == 0 {
				return 0, ErrMismatchedParentheses
			}
			ops = ops[:len(ops)-1]
		} else {
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(rune(char)) {
				if len(values) < 2 {
					return 0, ErrInsufficientOperands
				}
				op := ops[len(ops)-1]
				ops = ops[:len(ops)-1]

				val2 := values[len(values)-1]
				values = values[:len(values)-1]
				val1 := values[len(values)-1]
				values = values[:len(values)-1]

				res, err := operation(val1, val2, op)
				if err != nil {
					return 0, err
				}
				values = append(values, res)
			}
			ops = append(ops, rune(char))
		}
		i++
	}

	for len(ops) > 0 {
		if len(values) < 2 {
			return 0, ErrInsufficientOperands
		}
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]

		val2 := values[len(values)-1]
		values = values[:len(values)-1]
		val1 := values[len(values)-1]
		values = values[:len(values)-1]

		res, err := operation(val1, val2, op)
		if err != nil {
			return 0, err
		}
		values = append(values, res)
	}

	if len(values) != 1 {
		return 0, ErrInvalidExpression
	}
	return values[0], nil
}
