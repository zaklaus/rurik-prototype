package main

import (
	"strings"

	rl "github.com/zaklaus/raylib-go/raylib"
)

func questInitMiscCommands(q *questManager) {
	q.registerCommand("say", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("say", qs, qt, len(args), 1)
		}

		res, ok := qs.getResource(args[0])

		if !ok {
			return questCommandErrorThing("say", "message", qs, qt, args[0])
		}

		qs.printf(qt, "temp saying[%s]: %s", args[0], qs.processText(res.content))
		PushNotification(qs.processText(res.content), rl.RayWhite)

		return true
	})

	q.registerCommand("play", func(qs *quest, qt *questTask, args []string) bool {
		qs.printf(qt, "playing something")
		return true
	})

	q.registerCommand("give", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 2 {
			return questCommandErrorArgCount("give", qs, qt, len(args), 2)
		}

		amount, ok := qs.getNumberOrVariable(args[1])

		if !ok {
			return questCommandErrorArgType("give", qs, qt, args[1], "string", "integer")
		}

		qs.printf(qt, "giving %f of %s", amount, args[0])
		return true
	})

	q.registerCommand("log", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 2 {
			return questCommandErrorArgCount("log", qs, qt, len(args), 2)
		}

		logType := args[0]

		switch logType {
		case "str":
			qs.printf(qt, "%s", strings.Join(args[1:], " "))
		case "num":
			num, ok := qs.getNumberOrVariable(args[1])

			if ok {
				qs.printf(qt, "%f", num)
			} else {
				qs.printf(qt, "<unresolved>")
			}
		case "vec":
			vec, ok := qs.getVector(args[1])

			if ok {
				qs.printf(qt, "[%f, %f]", vec.X, vec.Y)
			} else {
				qs.printf(qt, "[<unresolved>]")
			}
		}

		return true
	})
}
