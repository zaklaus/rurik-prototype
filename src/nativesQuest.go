package main

import "github.com/zaklaus/rurik/src/core"

func registerQuestNatives() {
	core.RegisterNative("quest", func(jsData core.InvokeData) interface{} {
		var data struct {
			ID        int64
			EventName string
			Args      []float64
		}
		data.ID = -1

		core.DecodeInvokeData(&data, jsData)

		currentGameMode.quests.callEvent(data.ID, data.EventName, data.Args)
		return nil
	})

	core.RegisterNative("addQuest", func(jsData core.InvokeData) interface{} {
		var data struct {
			Name string
		}
		core.DecodeInvokeData(&data, jsData)

		_, _, id := currentGameMode.quests.addQuest(data.Name, nil)
		return id
	})
}
