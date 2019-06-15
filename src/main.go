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

var (
	currentGameMode *gameMode
)

func main() {
	currentGameMode = &gameMode{}

	rl.SetTraceLog(0)

	core.InitCore("Darkorbia", windowW, windowH, screenW, screenH)
	registerClasses()
	registerInputActions()
	registerCollisionTypes()
	registerEvents()
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

func registerEvents() {
	core.RegisterNative("quest", func(jsData core.InvokeData) interface{} {
		var data struct {
			name string
			args []int
		}
		core.DecodeInvokeData(&data, jsData)

		currentGameMode.quests.callEvent(data.name, data.args)
		return nil
	})
}

func registerClasses() {
	core.RegisterClass("player", NewPlayer)
	core.RegisterClass("water", NewWater)
	core.RegisterClass("ladder", NewLadder)
	core.RegisterClass("ball", NewBall)
}
