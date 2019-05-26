package main

import "math"

func roundFloat(x float32) float32 {
	return float32(math.Round(float64(x)))
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
