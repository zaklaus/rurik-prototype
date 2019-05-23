all:
	go build -o build/darkorbia.exe src/*.go

play:
	./build/darkorbia.exe

bt: all play
