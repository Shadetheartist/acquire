package acquire_2

import (
	"acquire/internal/util"
	"git.sr.ht/~bonbon/gmcts"
)

type Action_PickHotelToFound struct {
	Hotel Hotel
}

func (a Action_PickHotelToFound) Type() ActionType {
	return ActionType_PickHotelToFound
}

func (game *Game) getFoundHotelActions() []gmcts.Action {
	return util.Map(game.Computed.AvailableChains, func(val Hotel) gmcts.Action {
		return Action_PickHotelToFound{Hotel: val}
	})
}

func (game *Game) applyPickHotelToFoundAction(action Action_PickHotelToFound) {

	tile := game.LastPlacedTile

	newPlacedHotel := game.placeTileOnBoard(tile, action.Hotel)
	propagateHotelChain(game, newPlacedHotel)

	err := game.CurrentPlayer().takeStockFromBank(game, action.Hotel, 1)
	if err != nil {
		// this indicates that an action was generated which was invalid
		panic(err)
	}

	game.NextActionType = ActionType_PurchaseStock
}
