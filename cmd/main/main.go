package main

import (
	"acquire/internal/acquire"
	"acquire/internal/console_interface"
	"math/rand"
)

func main() {

	rand.Seed(int64(2))
	inputInterface := &console_interface.ConsoleInputInterface{}

	game := acquire.NewGame(inputInterface)
	console_interface.Render(game)

	for {
		game.Step()
		if game.IsOver {
			break
		}
		console_interface.Render(game)
	}

}
