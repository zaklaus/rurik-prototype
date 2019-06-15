package main

import "github.com/zaklaus/rurik/src/core"

func questInitBaseCommands(q *questManager) {
	q.registerCommand("variable", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("variable", qs, qt, len(args), 1)
		}

		qs.setVariable(args[0], 0)

		qs.printf(qt, "variable '%s' was declared", args[0])

		return true
	})

	q.registerCommand("setvar", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 2 {
			return questCommandErrorArgCount("setvar", qs, qt, len(args), 2)
		}

		val, ok := qs.getNumberOrVariable(args[1])

		if !ok {
			return questCommandErrorArgType("setvar", qs, qt, args[1], "string", "integer")
		}

		qs.setVariable(args[0], val)

		qs.printf(qt, "variable '%s' was set to: %d", args[0], val)

		return true
	})

	q.registerCommand("timer", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 2 {
			return questCommandErrorArgCount("timer", qs, qt, len(args), 2)
		}

		duration, ok := qs.getNumberOrVariable(args[1])

		if !ok {
			return questCommandErrorArgType("timer", qs, qt, args[1], "string", "integer")
		}

		qs.timers[args[0]] = questTimer{
			time:     -1,
			duration: float32(duration),
		}

		qs.printf(qt, "timer '%s' was declared with duration: %d", args[0], duration)

		return true
	})

	q.registerCommand("stage", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("stage", qs, qt, len(args), 1)
		}

		res, ok := qs.getResource(args[0])

		if !ok {
			return questCommandErrorThing("stage", "resource", qs, qt, args[0])
		}

		stageID := atoiUnsafe(args[0])

		qs.stages[stageID] = questStage{
			step:  res.content,
			state: qsInProgress,
		}

		qs.printf(qt, "stage '%d' has been added!", stageID)

		return true
	})

	q.registerCommand("stdone", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("stdone", qs, qt, len(args), 1)
		}

		stageID := atoiUnsafe(args[0])
		sta, ok := qs.stages[stageID]

		if !ok {
			return questCommandErrorThing("stdone", "resource", qs, qt, args[0])
		}

		qs.printf(qt, "stage '%d' has succeeded!", stageID)

		sta.state = qsFinished
		qs.stages[stageID] = sta

		return true
	})

	q.registerCommand("stfail", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("stfail", qs, qt, len(args), 1)
		}

		stageID := atoiUnsafe(args[0])
		sta, ok := qs.stages[stageID]

		if !ok {
			return questCommandErrorThing("stfail", "resource", qs, qt, args[0])
		}

		qs.printf(qt, "stage '%d' has failed!", stageID)

		sta.state = qsFailed
		qs.stages[stageID] = sta

		return true
	})

	q.registerCommand("repeat", func(qs *quest, qt *questTask, args []string) bool {
		qt.pc = -1

		qs.printf(qt, "repeating task '%s'!", qt.name)

		return true
	})

	q.registerCommand("fire", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("fire", qs, qt, len(args), 1)
		}

		tm, ok := qs.timers[args[0]]

		if !ok {
			return questCommandErrorThing("fire", "timer", qs, qt, args[0])
		}

		qs.printf(qt, "timer '%s' was fired!", args[0])
		tm.time = tm.duration
		qs.timers[args[0]] = tm

		return true
	})

	q.registerCommand("stop", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("stop", qs, qt, len(args), 1)
		}

		tm, ok := qs.timers[args[0]]

		if !ok {
			return questCommandErrorThing("stop", "timer", qs, qt, args[0])
		}

		qs.printf(qt, "timer '%s' was stopped!", args[0])
		tm.time = -1
		qs.timers[args[0]] = tm

		return true
	})

	q.registerCommand("done", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("done", qs, qt, len(args), 1)
		}

		tm, ok := qs.timers[args[0]]

		if !ok {
			return questCommandErrorThing("done", "timer", qs, qt, args[0])
		}

		state := tm.time == 0

		if state {
			qs.printf(qt, "timer '%s' is done!", args[0])
		}

		return state
	})

	q.registerCommand("finish", func(qs *quest, qt *questTask, args []string) bool {
		qs.state = qsFinished

		qs.printf(qt, "quest '%s' has been finished!", qs.name)

		return true
	})

	q.registerCommand("fail", func(qs *quest, qt *questTask, args []string) bool {
		qs.state = qsFailed

		qs.printf(qt, "quest '%s' has been failed!", qs.name)

		return true
	})

	q.registerCommand("pop", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("pop", qs, qt, len(args), 1)
		}

		if len(qt.eventArgs) == 0 {
			return questCommandErrorEventArgsEmpty("pop", qs, qt)
		}

		val := qt.eventArgs[0]
		qt.eventArgs = qt.eventArgs[1:]

		qs.setVariable(args[0], val)

		qs.printf(qt, "event pop value '%d' for: '%s'", val, args[0])

		return true
	})

	q.registerCommand("when", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 1 {
			return questCommandErrorArgCount("when", qs, qt, len(args), 1)
		}

		lhs, ok := qs.getNumberOrVariable(args[0])

		if !ok {
			return questCommandErrorArgType("when", qs, qt, args[0], "string", "integer")
		}

		if len(args) == 1 {
			return lhs > 0
		}

		rhs, ok2 := qs.getNumberOrVariable(args[2])

		if !ok2 {
			return questCommandErrorArgType("when", qs, qt, args[2], "string", "integer")
		}

		switch args[1] {
		case kwBelow:
			return lhs < rhs
		case kwAbove:
			return lhs > rhs
		case kwEquals:
			return lhs == rhs
		case kwNotEquals:
			return lhs != rhs
		default:
			return questCommandErrorArgComp("when", qs, qt, args[2])
		}
	})

	q.registerCommand("invoke", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) < 1 {
			return questCommandErrorArgCount("invoke", qs, qt, len(args), 1)
		}

		core.FireEvent(args[0], args[1:])
		return true
	})

	// temp
	q.registerCommand("say", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("say", qs, qt, len(args), 1)
		}

		res, ok := qs.getResource(args[0])

		if !ok {
			return questCommandErrorThing("say", "message", qs, qt, args[0])
		}

		qs.printf(qt, "temp saying[%s]: %s", args[0], qs.processText(res.content))

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

		qs.printf(qt, "giving %d of %s", amount, args[0])
		return true
	})
}
