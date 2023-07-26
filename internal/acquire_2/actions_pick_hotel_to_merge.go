package acquire_2

import (
	"acquire/internal/util"
	"git.sr.ht/~bonbon/gmcts"
)

type Action_PickHotelToMerge struct {
	Hotel Hotel
}

func (a Action_PickHotelToMerge) Type() ActionType {
	return ActionType_PickHotelToMerge
}

func (game *Game) getPickHotelToMergeActions() []gmcts.Action {
	return util.Map(game.Computed.LargestChains, func(val Hotel) gmcts.Action {
		return Action_PickHotelToMerge{Hotel: val}
	})
}

func (game *Game) applyPickHotelToMergeAction(action Action_PickHotelToMerge) {
	game.MergerState.AcquiringHotel = action.Hotel

	// remove the acquiring hotel chain from the list to merge (by setting it to zero)
	game.MergerState.ChainsToMerge[action.Hotel] = 0

	game.NextActionType = ActionType_Merge
}
