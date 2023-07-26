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

func (game *Game) ApplyAction(gmctsAction gmcts.Action) (gmcts.Game, error) {

	gameStruct := *game
	clone := gameStruct

	action, ok := gmctsAction.(IAction)
	if !ok {
		panic("action type was not convertable to IAction")
	}
	switch action.Type() {
	case ActionType_PlaceTile:
		clone.applyPlaceTileAction(action.(Action_PlaceTile))
		break
	case ActionType_PickHotelToFound:
		clone.applyPickHotelToFoundAction(action.(Action_PickHotelToFound))
		break
	case ActionType_PickHotelToMerge:
		clone.applyPickHotelToMergeAction(action.(Action_PickHotelToMerge))
		break
	case ActionType_Merge:
		clone.applyMergeHotel(action.(Action_Merge))
		break
	case ActionType_PurchaseStock:
		clone.applyPurchaseStockAction(action.(Action_PurchaseStock))
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
