package ai

import (
	"acquire/internal/acquire_2"
	"git.sr.ht/~bonbon/gmcts"
)

type SmartAgent struct {
}

func NewSmartAgent() *SmartAgent {
	return &SmartAgent{}
}

func (agent SmartAgent) SelectAction(game *acquire_2.Game, actions []gmcts.Action) (gmcts.Action, error) {
	mcts := gmcts.NewMCTS(game)

	//Spawn a new tree and play 1000 game simulations
	tree := mcts.SpawnTree()
	tree.SearchRounds(5)

	//Add the searched tree into the mcts tree collection
	mcts.AddTree(tree)

	//Get the best action based off of the trees collected from mcts.AddTree()
	bestAction := mcts.BestAction()

	return bestAction, nil
}
