package main

import (
	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	screenW = 960
	screenH = 540
	windowW = 1600
	windowH = 900
)

const (
	collisionMeta uint32 = core.FirstCollisionType
	collisionPawn
)

var (
	currentGameMode *gameMode
)

func main() {
	currentGameMode = &gameMode{}

	rl.SetTraceLog(0)
	rl.SetExitKey(0)

	core.InitUserEvents = registerNatives
	core.InitCore("Darkorbia", windowW, windowH, screenW, screenH)
	registerClasses()
	registerInputActions()
	registerCollisionTypes()
	core.Run(currentGameMode, true)
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

func quitGame() {
	core.CloseGame()
}
