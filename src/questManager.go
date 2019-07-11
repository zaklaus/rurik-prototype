package main

import (
	"log"
)

const (
	maxQuests = 5
)

var (
	stepCounter = 0
)

type questCommandTable func(qs *quest, qt *questTask, args []string) bool

type questManager struct {
	commands map[string]questCommandTable
	quests   []quest
}

func makeQuestManager() questManager {
	res := questManager{
		commands: map[string]questCommandTable{},
		quests:   []quest{},
	}

	questInitBaseCommands(&res)

	return res
}

func (q *questManager) getActiveQuests() []*quest {
	qs := []*quest{}

	for _, v := range q.quests {
		if v.state == qsInProgress && !v.runsInBackground {
			qs = append(qs, &v)
		}
	}

	return qs
}

func (q *questManager) addQuest(tplName string, details map[string]float64) (bool, string, int64) {
	qd := parseQuest(tplName)

	if !qd.runsInBackground && len(q.getActiveQuests()) >= maxQuests {
		return false, "Maximum number of quests has been reached!", -1
	}

	if qd == nil {
		return false, "Quest template could not be found!", -1
	}

	if details == nil {
		details = map[string]float64{}
	}

	tasks := []questTask{}

	for _, v := range qd.taskDef {
		tasks = append(tasks, questTask{
			questTaskDef: v,
			variables:    map[string]questVar{},
		})
	}

	processedDetails := map[string]questVar{}

	for k, v := range details {
		processedDetails[k] = questVar{
			kind:  kindNumber,
			value: &questVarNumber{value: v},
		}
	}

	tasks[0].variables = processedDetails

	qn := quest{
		ID:       getNewID(),
		name:     tplName,
		questDef: *qd,
		state:    qsInProgress,
		timers:   map[string]questTimer{},
		stages:   map[int]questStage{},
		tasks:    tasks,
	}

	qn.activeQuestTask = &qn.tasks[0]

	for _, v := range qn.tasks {
		qn.setVariable(v.name, 0)
	}

	for qn.processTask(q, &qn.tasks[0]) {
		// process the whole entry point
	}

	q.quests = append(q.quests, qn)

	log.Printf("Quest '%s' with title '%s' has been added!", tplName, qd.title)

	return true, "", qn.ID
}

func (q *questManager) reset() {
	q.quests = []quest{}
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

func (q *questManager) callEvent(id int64, eventName string, args []float64) {
	for i := range q.quests {
		v := &q.quests[i]

		if id != -1 && id != v.ID {
			continue
		}

		v.callEvent(q, eventName, args)
	}
}
