all:
	go build -o build/game.exe src/*.go

play:
	./build/game.exe

bt: all play
