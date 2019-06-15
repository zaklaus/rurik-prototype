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
	qsInProgress = iota
	qsFinished
	qsFailed
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

type questStage struct {
	step  string
	state int
}

func (qs *quest) printf(qt *questTask, format string, args ...interface{}) {
	log.Printf("Quest '%s':'%s'(%d): %s", qs.name, qt.name, qt.pc, fmt.Sprintf(format, args...))
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

func (qs *quest) processText(content string) string {
	for k, v := range qs.variables {
		content = strings.ReplaceAll(content, fmt.Sprintf("%%%s%%", k), strconv.Itoa(v))
	}

	return content
}

func (qs *quest) resolveVariables(expr string) string {
	for k, v := range qs.variables {
		expr = strings.ReplaceAll(expr, k, strconv.Itoa(v))
	}

	return expr
}

func (qs *quest) setVariable(name string, val int) {
	qs.variables[name] = val
}

func (qs *quest) processTimers() {
	for k, v := range qs.timers {
		if v.time >= 0 {
			v.time -= system.FrameTime

			if v.time < 0 {
				v.time = 0
			}

			qs.timers[k] = v
		}

		qs.setVariable(k, int(core.RoundFloatToInt32(v.time)))
	}
}

func (qs *quest) processTask(q *questManager, qt *questTask) bool {
	if qt.pc >= len(qt.commands) {
		qt.isDone = true
		return false
	}

	qs.processVariables()

	cmd := qt.commands[qt.pc]
	ok, err := q.dispatchCommand(qs, qt, cmd.name, cmd.args)

	if err {
		qt.isDone = true
		return false
	}

	if !ok {
		return false
	}

	qt.pc++
	return true
}

func (qs *quest) processTasks(q *questManager) {
	for i := range qs.tasks {
		v := &qs.tasks[i]

		if v.isDone || v.isEvent {
			continue
		}

		for qs.processTask(q, v) {
			// task is being processed
		}

		state := 0

		if v.isDone {
			state = 1
		}

		qs.setVariable(v.name, state)
	}
}

func (qs *quest) callEvent(q *questManager, name string, args []int) {
	for i := range qs.tasks {
		v := &qs.tasks[i]

		v.isDone = false
		v.eventArgs = args[:]

		for qs.processTask(q, v) {
			// task is being processed
		}
	}
}

func (qs *quest) processVariables() {
	qs.setVariable("$random", rand.Int())
	qs.setVariable("$step", stepCounter)
	qs.setVariable("$time", int(core.RoundFloatToInt32(rl.GetTime())))

	// temp
	qs.setVariable("$pc.health", int(barStats[barHealth].Value))
}
