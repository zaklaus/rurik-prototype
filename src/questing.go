package main

import (
	"fmt"
	"log"
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
	ID               int64
	name             string
	runsInBackground bool // used by events, they don't count as an actual quest
	state            int
	timers           map[string]questTimer
	stages           map[int]questStage
	tasks            []questTask
	activeQuestTask  *questTask
	questDef
}

const (
	kindNumber = iota
	kindVector
)

type questVarData interface {
	str() string
}

type questVar struct {
	kind  int
	value questVarData
}

type questTask struct {
	variables map[string]questVar
	questTaskDef
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

func (qs *quest) getNumberOrVariable(arg string) (float64, bool) {
	val, err := strconv.ParseFloat(arg, 64)

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

		return res.(float64), true
	}

	return val, true
}

func (qs *quest) getRelevantVariables() (a map[string]questVar) {
	if qs.activeQuestTask == &qs.tasks[0] {
		return qs.tasks[0].variables
	}

	a = map[string]questVar{}

	for k, v := range qs.tasks[0].variables {
		a[k] = v
	}

	for k, v := range qs.activeQuestTask.variables {
		a[k] = v
	}

	return a
}

func (qs *quest) processText(content string) string {
	for k, v := range qs.getRelevantVariables() {
		content = strings.ReplaceAll(content, fmt.Sprintf("%%%s%%", k), v.value.str())
	}

	return content
}

func (qs *quest) resolveVariables(expr string) string {
	for k, v := range qs.getRelevantVariables() {
		expr = strings.ReplaceAll(expr, k, v.value.str())
	}

	return expr
}

func (qs *quest) getTaskOverride(name string) *questTask {
	aq := qs.activeQuestTask

	_, ok := qs.tasks[0].variables[name]

	if ok {
		aq = &qs.tasks[0]
	}

	return aq
}

func (qs *quest) setVariable(name string, val float64) {
	qs.getTaskOverride(name).variables[name] = questVar{
		kind:  kindNumber,
		value: &questVarNumber{value: val},
	}
}

func (qs *quest) setVector(name string, val rl.Vector2) {
	qs.getTaskOverride(name).variables[name] = questVar{
		kind:  kindVector,
		value: &questVarVector{value: val},
	}
}

func (qs *quest) getVariable(name string) (float64, bool) {
	vars := qs.getRelevantVariables()

	val, ok := vars[name]

	if !ok {
		return 0, false
	}

	return val.value.(*questVarNumber).value, true
}

func (qs *quest) getVector(name string) (rl.Vector2, bool) {
	vars := qs.getRelevantVariables()

	val, ok := vars[name]

	if !ok {
		return rl.Vector2{}, false
	}

	return val.value.(*questVarVector).value, true
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

		qs.setVariable(k, float64(core.RoundFloatToInt32(v.time)))
	}
}

func (qs *quest) processTask(q *questManager, qt *questTask) bool {
	if qt.pc >= len(qt.commands) {
		qt.isDone = true
		return false
	}

	qs.activeQuestTask = qt

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

		qs.setVariable(v.name, float64(state))
	}
}

func (qs *quest) callEvent(q *questManager, name string, args []float64) {
	for i := range qs.tasks {
		v := &qs.tasks[i]

		if v.name != name {
			continue
		}

		v.isDone = false
		v.eventArgs = args[:]

		for qs.processTask(q, v) {
			// task is being processed
		}
	}
}

func (qs *quest) processVariables() {
	qt := qs.activeQuestTask
	qs.activeQuestTask = &qs.tasks[0]
	qs.setVariable("$random", float64(rand.Int()))
	qs.setVariable("$frandom", rand.Float64())
	qs.setVariable("$step", float64(stepCounter))
	qs.setVariable("$time", float64(rl.GetTime()))

	// player
	qs.setVector("$pc.position", core.LocalPlayer.Position)

	// temp
	qs.setVariable("$pc.health", float64(barStats[barHealth].Value))
	qs.activeQuestTask = qt
}
