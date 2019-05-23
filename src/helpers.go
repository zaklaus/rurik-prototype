package main

import "math"

func roundFloat(x float32) float32 {
	return float32(math.Round(float64(x)))
}
