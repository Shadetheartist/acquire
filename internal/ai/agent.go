package ai

import (
	"acquire/internal/acquire"
	"git.sr.ht/~bonbon/gmcts"
)

type IAgent interface {
	SelectAction(game *acquire.Game, actions []gmcts.Action) (gmcts.Action, error)
}
