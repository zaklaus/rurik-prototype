package main

import (
	"fmt"
	"log"
)

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

func questCommandErrorEventArgsEmpty(cmd string, qs *quest, qt *questTask) bool {
	log.Printf("%s event's arg stack is already empty!", questCommandErrorBase(cmd, qs, qt))
	return false
}
