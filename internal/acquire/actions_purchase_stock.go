package acquire

import (
	"acquire/internal/util"
	"fmt"
	"git.sr.ht/~bonbon/gmcts"
	"os"
)

type StockPurchase struct {
	Hotel  Hotel
	Amount int
}

type Action_PurchaseStock struct {
	Purchases [3]StockPurchase
}

func (a Action_PurchaseStock) Type() ActionType {
	return ActionType_PurchaseStock
}

func (game *Game) getPurchaseStockActions() []gmcts.Action {
	// keeping this simple for now
	actions := util.Map(game.Computed.ActiveChains, func(val Hotel) gmcts.Action {
		return Action_PurchaseStock{
			Purchases: [3]StockPurchase{
				{
					Hotel:  val,
					Amount: 0,
				},
			},
		}
	})

	// default action to buy nothing
	actions = append(actions, Action_PurchaseStock{
		Purchases: [3]StockPurchase{
			{
				Hotel:  NoHotel,
				Amount: 0,
			},
		},
	})

	return actions
}

func (game *Game) applyPurchaseStockAction(action Action_PurchaseStock) {

	for _, purchase := range action.Purchases {
		if purchase.Hotel == NoHotel || purchase.Amount == 0 {
			continue
		}

		err := game.CurrentPlayer().buyStock(game, purchase.Hotel, purchase.Amount)
		if err != nil {
			panic(err)
		}
	}

	// going to merge this and the draw tile 'action' for better efficiency and it's easier

	// take a new tile from the bank
	// ignoring the error, if there weren't any tiles left to take
	err := game.CurrentPlayer().takeTileFromBank(game)
	if err != nil {
		// panic(err)
	}

	// game always ends at the end of the player's turn
	if game.WillEnd {
		game.end("a player has declared the game is over")
		return
	}

	game.NextActionType = ActionType_PlaceTile
	game.Turn++
}

func (game *Game) Abort(reason string) {
	fmt.Println("Game aborted: " + reason)
	os.Exit(1)
}

func (game *Game) end(reason string) {

	game.IsOver = true

	// payout shareholder bonuses
	for _, hotel := range HotelChainList {
		game.payShareholderBonuses(hotel)
	}

	// sell all stocks
	for _, p := range game.Players {
		for _, hotel := range HotelChainList {
			stock := Stock(hotel)
			err := p.sellStock(game, stock, p.Stocks[hotel.Index()])
			if err != nil {
				panic(err)
			}
		}
	}
}

func (game *Game) NumRemainingTiles() int {
	c := 0
	for _, t := range game.Tiles {
		if t != NoTile {
			c++
		}
	}
	return c
}

func (game *Game) hasRemainingTiles() bool {
	for _, t := range game.Tiles {
		if t != NoTile {
			return true
		}
	}
	return false
}

// canEnd
// returns true if it's possible for a player to 'declare' the game over
func (game *Game) canEnd() (string, bool) {

	if game.hasRemainingTiles() {
		return "no tiles left somehow", true
	}

	// if there are any chains larger than 40, the game can end
	for _, size := range game.Computed.LargestChains {
		if size > 40 {
			return "there is a chain with length 41", true
		}
	}

	// get any unsafe active chains
	hasUnsafe := func() bool {
		for _, size := range game.Computed.ActiveChains {
			if size <= 10 {
				return true
			}
		}
		return false
	}()

	// if there are no remaining unsafe (from merger) hotels, the game can end
	if !hasUnsafe {
		return "all chains on the board are safe", true
	}

	return "", false
}
