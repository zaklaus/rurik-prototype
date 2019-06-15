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
	collisionPawn uint32 = core.FirstCollisionType
)

var (
	currentGameMode *gameMode
)

func main() {
	currentGameMode = &gameMode{}

	rl.SetTraceLog(0)

	core.InitUserEvents = registerEvents
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

func registerEvents() {
	core.RegisterNative("quest", func(jsData core.InvokeData) interface{} {
		var data struct {
			Name string
			Args []int64
		}
		core.DecodeInvokeData(&data, jsData)

		args := []int{}

		for _, x := range data.Args {
			args = append(args, int(x))
		}

		currentGameMode.quests.callEvent(data.Name, args)
		return nil
	})

	core.RegisterNative("addQuest", func(jsData core.InvokeData) interface{} {
		var data struct {
			Name string
		}
		core.DecodeInvokeData(&data, jsData)

		currentGameMode.quests.addQuest(data.Name, nil)
		return nil
	})
}

func registerClasses() {
	core.RegisterClass("player", NewPlayer)
	core.RegisterClass("water", NewWater)
	core.RegisterClass("ladder", NewLadder)
	core.RegisterClass("ball", NewBall)
}
