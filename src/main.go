package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	screenW = 640
	screenH = 360
	windowW = 1600
	windowH = 900
)

const (
	collisionPawn uint32 = core.FirstCollisionType
)

func main() {
	gm := &gameMode{}

	registerClasses()

	core.InitCore("Darkorbia", windowW, windowH, screenW, screenH)
	registerClasses()
	registerInputActions()
	registerCollisionTypes()
	core.Run(gm, true)
}

func registerInputActions() {
	system.BindInputAction("jump", system.InputAction{
		AllKeys:    []int32{rl.KeySpace},
		JoyButtons: []int32{rl.GamepadXboxButtonA},
	})
}

func registerCollisionTypes() {
	core.AddCollisionType("pawn", collisionPawn)
}

func registerClasses() {
	core.RegisterClass("player", NewPlayer)
	core.RegisterClass("water", NewWater)
	core.RegisterClass("ladder", NewLadder)
	core.RegisterClass("ball", NewBall)
}
