package main

import (
	"acquire/internal/acquire_2"
	"acquire/internal/ai"
	"acquire/internal/console_interface"
	"math/rand"
)

func main() {

	rand.Seed(int64(2))

	game := acquire_2.NewGame()

	agents := make(map[int]ai.IAgent)
	for _, player := range game.Players {
		agents[player.Id] = ai.NewStupidAgent()
	}
	agents[game.Players[0].Id] = ai.NewSmartAgent()

	for !game.IsTerminal() {

		agent := agents[game.CurrentPlayer().Id]
		actions := game.GetActions()
		action, err := agent.SelectAction(game, actions)
		if err != nil {
			panic(err)
		}
		newGame, _ := game.ApplyAction(action)
		game = newGame.(*acquire_2.Game)
		console_interface.Render(game)
	}

}
