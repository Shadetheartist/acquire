package acquire

import (
	"acquire/internal/util"
	"fmt"
	"git.sr.ht/~bonbon/gmcts"
)

// this file contains the MCTS related functions

type ActionType int

func (at ActionType) String() string {
	switch at {
	case ActionType_PlaceTile:
		return "Place Tile"
	case ActionType_PickHotelToFound:
		return "Pick Hotel To Found"
	case ActionType_PickHotelToMerge:
		return "Pick Hotel To Merge"
	case ActionType_Merge:
		return "Merge"
	case ActionType_PurchaseStock:
		return "Purchase Stock"
	default:
		panic("wtf")
	}
}

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
	String(game *Game) string
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
		clone.applyPlaceTileAction(util.AsType[Action_PlaceTile](action))
	case ActionType_PickHotelToFound:
		clone.applyPickHotelToFoundAction(util.AsType[Action_PickHotelToFound](action))
		break
	case ActionType_PickHotelToMerge:
		clone.applyPickHotelToMergeAction(util.AsType[Action_PickHotelToMerge](action))
		break
	case ActionType_Merge:
		clone.applyMergeHotel(util.AsType[Action_Merge](action))
		break
	case ActionType_PurchaseStock:
		clone.applyPurchaseStockAction(util.AsType[Action_PurchaseStock](action))
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
