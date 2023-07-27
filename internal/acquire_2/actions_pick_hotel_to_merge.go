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
	pos := game.LastPlacedTile.Pos()
	neighboringHotels := getNeighbors(game, pos)
	chainsInNeighbors := getChainsInNeighbors(neighboringHotels)
	acquiredChains := util.Filter(chainsInNeighbors, func(val Hotel) bool {
		return val != action.Hotel
	})
	largestAcquiredChains, _ := game.getLargestChainsOf(acquiredChains)

	// remove the acquiring hotel chain from the list to merge (by setting it to zero)
	game.MergerState.ChainsToMerge[action.Hotel.Index()] = 0

	// prepare the 'chains to merge' array
	for _, h := range HotelChainList {
		// ok = this hotel is in the 'largest chains' slice, but isn't the largest chain
		_, ok := util.IndexOf(largestAcquiredChains, h)
		if ok {
			game.MergerState.ChainsToMerge[h.Index()] = len(game.Players)
		}
	}

	game.MergerState.AcquiringHotel = action.Hotel
	game.NextActionType = ActionType_Merge
}
