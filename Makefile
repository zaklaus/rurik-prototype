all:
	go build -o build/game.exe src/*.go

play:
	./build/game.exe

perf:
	go tool pprof --pdf build/cpu.pprof > build/shit.pdf

bt: all play
