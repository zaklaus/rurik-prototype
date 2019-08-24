package main

import (
	"fmt"
	"math"

	rl "github.com/zaklaus/raylib-go/raylib"
)

type questVarNumber struct {
	value float64
}

func (v *questVarNumber) str() string {
	if math.Floor(v.value) == v.value {
		return fmt.Sprintf("%d", int64(v.value))
	}

	return fmt.Sprintf("%f", v.value)
}

type questVarVector struct {
	value rl.Vector2
}

func (v *questVarVector) str() string {
	return fmt.Sprintf("[%f, %f]", v.value)
}
