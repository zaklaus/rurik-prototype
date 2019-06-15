package main

import "log"

const (
	maxQuests = 5
)

var (
	stepCounter = 0
)

type questCommandTable func(qs *quest, qt *questTask, args []string) bool

type questManager struct {
	quests   []quest
	commands map[string]questCommandTable
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

func (q *questManager) callEvent(name string, args []int) {
	for i := range q.quests {
		q.quests[i].callEvent(q, name, args)
	}
}
