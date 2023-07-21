package main

import (
	"acquire/internal/acquire"
	"acquire/internal/console_interface"
	"acquire/internal/util"
	"fmt"
	"git.sr.ht/~bonbon/gmcts"
	"math/rand"
)

func main() {

	rand.Seed(int64(2))
	inputInterface := &console_interface.ConsoleInputInterface{}

	game := acquire.NewGame(inputInterface)

	for !game.IsTerminal() {
		mcts := gmcts.NewMCTS(game)

		//Spawn a new tree and play 1000 game simulations
		tree := mcts.SpawnTree()
		tree.SearchRounds(100)

		//Add the searched tree into the mcts tree collection
		mcts.AddTree(tree)

		//Get the best action based off of the trees collected from mcts.AddTree()
		bestAction := mcts.BestAction()

		if bestAction == nil {
			continue
		}

		//Update the game state using the tree's best action
		newGame, _ := game.ApplyAction(bestAction)
		game = newGame.(*acquire.Game)
		console_interface.Render(game)

	}

	console_interface.Render(game)

	for i, p := range game.Winners() {
		_, ap := util.Find(game.Players, func(val *acquire.Player) bool {
			return int(p) == val.Id
		})

		fmt.Printf("%d: %s with $%d\n", i+1, ap.Name(), ap.Inventory.Money)
	}
}
