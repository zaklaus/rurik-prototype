all:
	go build -o build/game.exe src/*.go

win:
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o build/game.exe src/*.go

play:
	./build/game.exe

perf:
	go tool pprof --pdf build/cpu.pprof > build/shit.pdf

bt: all play
wt: win play
