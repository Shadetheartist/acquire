package acquire_2

import (
	"fmt"
	"git.sr.ht/~bonbon/gmcts"
)

// this file contains the MCTS related functions

type ActionType int

// need to rework the game so that it is a state machine which is updated by these action types
const (
	ActionType_PlaceTile ActionType = iota
	ActionType_PickHotelToFound
	ActionType_PickHotelToMerge
	ActionType_Merge
	ActionType_PurchaseStock
)

type IAction interface {
	Type() ActionType
}

func asType[T any](action any) T {
	_action, ok := action.(T)
	if !ok {
		panic("action not of correct type")
	}
	return _action
}

func (game *Game) ApplyAction(gmctsAction gmcts.Action) (gmcts.Game, error) {

	gameStruct := *game
	clone := gameStruct

	action, ok := gmctsAction.(IAction)
	if !ok {
		panic("action type was not convertable to IAction")
	}

	switch action.Type() {
	case ActionType_PlaceTile:
		clone.applyPlaceTileAction(asType[Action_PlaceTile](action))
	case ActionType_PickHotelToFound:
		clone.applyPickHotelToFoundAction(asType[Action_PickHotelToFound](action))
		break
	case ActionType_PickHotelToMerge:
		clone.applyPickHotelToMergeAction(asType[Action_PickHotelToMerge](action))
		break
	case ActionType_Merge:
		clone.applyMergeHotel(asType[Action_Merge](action))
		break
	case ActionType_PurchaseStock:
		clone.applyPurchaseStockAction(asType[Action_PurchaseStock](action))
		break
	default:
		panic(fmt.Sprintf("action %d is not handled", action))
	}

	clone.Computed = NewComputed(&clone)

	return &clone, nil
}

func (game *Game) GetActions() []gmcts.Action {
	switch game.NextActionType {

	case ActionType_PlaceTile:
		return game.getPlaceTileActions()

	case ActionType_PickHotelToFound:
		return game.getFoundHotelActions()

	case ActionType_PickHotelToMerge:
		return game.getPickHotelToMergeActions()

	case ActionType_Merge:
		return game.getMergeHotelActions()

	case ActionType_PurchaseStock:
		return game.getPurchaseStockActions()

	default:
		panic("action type not implemented here")
	}
}
