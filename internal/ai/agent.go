package ai

import (
	"acquire/internal/acquire_2"
	"git.sr.ht/~bonbon/gmcts"
)

type IAgent interface {
	SelectAction(game *acquire_2.Game, actions []gmcts.Action) (gmcts.Action, error)
}
