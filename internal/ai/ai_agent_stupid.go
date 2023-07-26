package ai

import (
	"acquire/internal/acquire_2"
	"errors"
	"git.sr.ht/~bonbon/gmcts"
	"math/rand"
)

type StupidAgent struct {
}

func NewStupidAgent() *StupidAgent {
	return &StupidAgent{}
}

func (agent StupidAgent) SelectAction(game *acquire_2.Game, actions []gmcts.Action) (gmcts.Action, error) {
	if len(actions) == 0 {
		return nil, errors.New("no actions to select")
	}

	n := rand.Intn(len(actions))
	return actions[n], nil
}
