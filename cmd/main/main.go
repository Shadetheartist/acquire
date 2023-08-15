package main

import (
	"acquire/internal/acquire"
	"acquire/internal/ai"
	"fmt"
)

func main() {
	config := menu()
	runGame(config)
}

func runGame(config *GameConfig) *acquire.Game {
	game := acquire.NewGame()

	agents := make(map[int]ai.IAgent)

	// unset unused players
	for i := config.NumPlayers; i < len(game.Players); i++ {
		game.Players[i].Id = 0
	}

	// enabled from config
	for i := range config.PlayerTypes {
		if config.PlayerTypes[i] == Human {
			agents[game.Players[i].Id] = ai.NewHumanAgent()
		}
		if config.PlayerTypes[i] == AI {
			strength := config.AIPlayerStrengths[i]
			agents[game.Players[i].Id] = ai.NewSmartAgent(strength)
		}
	}

	for !game.IsTerminal() {

		agent := agents[game.CurrentPlayer().Id]
		actions := game.GetActions()

		// render before play
		acquire.Render(game)
		action, err := agent.SelectAction(game, actions)
		if err != nil {
			panic(err)
		}

		// describe play
		if _action, ok := action.(acquire.IAction); ok {
			fmt.Println(_action.String(game))
		}

		newGame, _ := game.ApplyAction(action)
		game = newGame.(*acquire.Game)
	}

	// render final board state
	acquire.Render(game)

	fmt.Println()
	reason, end := game.CanEnd()
	game2 := game
	game2.Computed = acquire.NewComputed(game2)
	if !end {
		reason = "Game was forced to end."
	}
	fmt.Println("End Reason: " + reason)

	return game
}
