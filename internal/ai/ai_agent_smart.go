package ai

import (
	"acquire/internal/acquire"
	"git.sr.ht/~bonbon/gmcts"
)

type SmartAgent struct {
	intelligence int
}

func NewSmartAgent(intelligence int) *SmartAgent {
	return &SmartAgent{
		intelligence: intelligence,
	}
}

func (agent SmartAgent) SelectAction(game *acquire.Game, _ []gmcts.Action) (gmcts.Action, error) {

	simGame := game
	simGame.Sim = true

	mcts := gmcts.NewMCTS(game)

	//Spawn a new tree and play some n number game simulations
	tree := mcts.SpawnTree()
	tree.SearchRounds(agent.intelligence)

	//Add the searched tree into the mcts tree collection
	mcts.AddTree(tree)

	//Get the best action based off of the trees collected from mcts.AddTree()
	bestAction := mcts.BestAction()

	return bestAction, nil
}
