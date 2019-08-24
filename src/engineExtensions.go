package main

var (
	globalIDCounter int64
)

func getNewID() int64 {
	v := globalIDCounter
	globalIDCounter++
	return v
}
