package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"

	rl "github.com/zaklaus/raylib-go/raylib"
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

const (
	maxQuests = 5
)

const (
	qsInProgress = iota
	qsFinished
	qsFailed
)

var (
	stepCounter = 0
)

type quest struct {
	name      string
	state     int
	variables map[string]int
	timers    map[string]questTimer
	stages    map[int]questStage
	questDef
}

type questTimer struct {
	time     float32
	duration float32
}

type questCommandTable func(qs *quest, qt *questTask, args []string) bool

type questManager struct {
	quests   []quest
	commands map[string]questCommandTable
}

type questStage struct {
	step  string
	state int
}

func newQuestManager() questManager {
	res := questManager{
		quests:   []quest{},
		commands: map[string]questCommandTable{},
	}

	questInitBaseCommands(&res)

	return res
}

func (qs *quest) printf(qt *questTask, format string, args ...interface{}) {
	log.Printf("Quest '%s':'%s'(%d): %s", qs.name, qt.name, qt.pc, fmt.Sprintf(format, args...))

}

func questCommandErrorBase(cmd string, qs *quest, qt *questTask) string {
	return fmt.Sprintf("Command '%s' failed at quest '%s':'%s'(%d): ", cmd, qs.name, qt.name, qt.pc)
}

func questCommandErrorArgCount(cmd string, qs *quest, qt *questTask, has, need int) bool {
	log.Printf("%s needs '%d' arguments, got: '%d'", questCommandErrorBase(cmd, qs, qt), need, has)
	return false
}

func questCommandErrorThing(cmd, thing string, qs *quest, qt *questTask, resName string) bool {
	log.Printf("%s %s '%s' could not be found", questCommandErrorBase(cmd, qs, qt), thing, resName)
	return false
}

func questCommandErrorArgType(cmd string, qs *quest, qt *questTask, argName, has, need string) bool {
	log.Printf("%s argument '%s' has to be '%s', got: '%s'", questCommandErrorBase(cmd, qs, qt), argName, need, has)
	return false
}

func questCommandErrorArgComp(cmd string, qs *quest, qt *questTask, argName string) bool {
	log.Printf("%s argument has to be either 'above,below,equals,!equals', got: '%s'", questCommandErrorBase(cmd, qs, qt), argName)
	return false
}

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
		if len(args) != 3 {
			return questCommandErrorArgCount("timer", qs, qt, len(args), 3)
		}

		duration, ok := qs.getNumberOrVariable(args[1])

		if !ok {
			return questCommandErrorArgType("timer", qs, qt, args[1], "string", "integer")
		}

		startTime, ok2 := qs.getNumberOrVariable(args[2])

		if !ok2 {
			return questCommandErrorArgType("timer", qs, qt, args[2], "string", "integer")
		}

		qs.timers[args[0]] = questTimer{
			time:     float32(startTime),
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

	// temp
	q.registerCommand("say", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 1 {
			return questCommandErrorArgCount("say", qs, qt, len(args), 1)
		}

		res, ok := qs.getResource(args[0])

		if !ok {
			return questCommandErrorThing("say", "message", qs, qt, args[0])
		}

		qs.printf(qt, "temp saying[%s]: %s", args[0], res.content)

		return true
	})

	q.registerCommand("play", func(qs *quest, qt *questTask, args []string) bool {
		qs.printf(qt, "playing something")
		return true
	})

	q.registerCommand("give", func(qs *quest, qt *questTask, args []string) bool {
		qs.printf(qt, "giving something")
		return true
	})
}

func (qs *quest) getResource(id string) (*questResource, bool) {
	val, err := strconv.Atoi(id)

	if err != nil {
		return nil, false
	}

	res, ok := qs.resources[val]

	if !ok {
		return nil, false
	}

	return &res, true
}

func (qs *quest) getNumberOrVariable(arg string) (int, bool) {
	val, err := strconv.Atoi(arg)

	if err != nil {
		exprStr := qs.resolveVariables(arg)
		expr, err := govaluate.NewEvaluableExpression(exprStr)

		if err != nil {
			return 0, false
		}

		res, err := expr.Evaluate(nil)

		if err != nil {
			return 0, false
		}

		return int(math.Floor(res.(float64))), true
	}

	return val, true
}

func (qs *quest) resolveVariables(expr string) string {
	for k, v := range qs.variables {
		expr = strings.ReplaceAll(expr, k, strconv.Itoa(v))
	}

	return expr
}

func (q *questManager) getActiveQuests() []*quest {
	qs := []*quest{}

	for _, v := range q.quests {
		if v.state == qsInProgress {
			qs = append(qs, &v)
		}
	}

	return qs
}

func (q *questManager) addQuest(tplName string, details map[string]int) (bool, string) {
	if len(q.getActiveQuests()) >= maxQuests {
		return false, "Maximum number of quests has been reached!"
	}

	qd := parseQuest(tplName)

	if qd == nil {
		return false, "Quest template could not be found!"
	}

	if details == nil {
		details = map[string]int{}
	}

	q.quests = append(q.quests, quest{
		name:      tplName,
		questDef:  *qd,
		state:     qsInProgress,
		variables: details,
		timers:    map[string]questTimer{},
		stages:    map[int]questStage{},
	})

	log.Printf("Quest '%s' with title '%s' has been added!", tplName, qd.title)

	return true, ""
}

func (qs *quest) setVariable(name string, val int) {
	qs.variables[name] = val
}

func (qs *quest) processTimers() {
	for k, v := range qs.timers {
		if v.time < 0 {
			continue
		}

		v.time -= system.FrameTime

		if v.time < 0 {
			v.time = 0
		}

		qs.timers[k] = v

		qs.setVariable(k, int(core.RoundFloatToInt32(v.time)))
	}
}

func (qs *quest) processTasks(q *questManager) {
	for i := range qs.tasks {
		v := &qs.tasks[i]

		if v.isDone {
			continue
		}

		ok := true
		var err bool

		for {
			if v.pc >= len(v.commands) {
				v.isDone = true
				break
			}

			qs.processVariables()

			cmd := v.commands[v.pc]
			ok, err = q.dispatchCommand(qs, v, cmd.name, cmd.args)

			if err {
				v.isDone = true
				break
			}

			if !ok {
				break
			}

			v.pc++
		}

		state := 0

		if v.isDone {
			state = 1
		}

		qs.setVariable(v.name, state)
	}
}

func (qs *quest) processVariables() {
	qs.setVariable("$random", rand.Int())
	qs.setVariable("$step", stepCounter)
	qs.setVariable("$time", int(core.RoundFloatToInt32(rl.GetTime())))

	// temp
	qs.setVariable("$pc.health", int(barStats[barHealth].Value))
}

func (q *questManager) registerCommand(name string, cb questCommandTable) {
	q.commands[name] = cb
}

func (q *questManager) dispatchCommand(qs *quest, qt *questTask, name string, args []string) (bool, bool) {
	cmd, ok := q.commands[name]

	if ok {
		return cmd(qs, qt, args), false
	}

	log.Printf("Quest '%s' has unrecognized command: '%s'!\n", qs.name, name)
	return false, true
}

func (q *questManager) processQuests() {
	for i := range q.quests {
		qs := &q.quests[i]

		if qs.state != qsInProgress {
			continue
		}

		qs.processTimers()
		qs.processTasks(q)
	}

	stepCounter++
}
