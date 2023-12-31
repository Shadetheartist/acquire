package acquire

import (
	"acquire/internal/util"
	"fmt"
	"git.sr.ht/~bonbon/gmcts"
)

type Action_PickHotelToFound struct {
	Hotel Hotel
}

func (a Action_PickHotelToFound) String(game *Game) string {
	return fmt.Sprintf("Player %s chooses to found %s.", game.CurrentPlayer().Name(), a.Hotel.String())
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

	remainingStock := game.Stocks[action.Hotel.Index()]
	err := game.CurrentPlayer().takeStockFromBank(game, action.Hotel, util.Min(1, remainingStock))
	if err != nil {
		// this indicates that an action was generated which was invalid
		panic(err)
	}

	game.NextActionType = ActionType_PurchaseStock
}
