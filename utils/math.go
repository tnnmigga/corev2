package utils

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func IfElse[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}
