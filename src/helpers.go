package main

import (
	"math"
	"strconv"

	"github.com/zaklaus/rurik/src/core"
)

func roundFloat(x float32) float32 {
	return float32(math.Round(float64(x)))
}

func roundInt(x int) float32 {
	return core.RoundInt32ToFloat(int32(x))
}

func roundInt32(x int) float32 {
	return core.RoundInt32ToFloat(int32(x))
}

func float64to32(x float64) float32 {
	return float32(x)
}

func absFloat(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func clamp(x, a, b float32) float32 {
	if x < a {
		return a
	} else if x > b {
		return b
	}

	return x
}

func atoiUnsafe(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
