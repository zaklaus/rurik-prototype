package main

import "math"

func roundFloat(x float32) float32 {
	return float32(math.Round(float64(x)))
}

func absFloat(x float32) float32 {
	return float32(math.Abs(float64(x)))
}
