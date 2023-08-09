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

// AsMap
// returns the purchases associated with this action in a map.
// NoHotel is not included as a key in the map
func (a Action_PurchaseStock) AsMap() map[Hotel]int {
	m := make(map[Hotel]int)
	for _, p := range a.Purchases {
		if p.Hotel == NoHotel {
			continue
		}
		m[p.Hotel] += p.Amount
	}
	return m
}

func (a Action_PurchaseStock) Type() ActionType {
	return ActionType_PurchaseStock
}

func generateCombinations(activeHotels []Hotel) [][]Hotel {

	var combinations [][]Hotel

	l := len(activeHotels)
	for i := 0; i < l; i++ {
		for j := i; j < l; j++ {
			for k := j; k < l; k++ {
				combinations = append(combinations, []Hotel{
					activeHotels[i],
					activeHotels[j],
					activeHotels[k],
				})
			}
		}
	}

	return combinations
}

func (game *Game) getPurchaseStockActions() []gmcts.Action {

	options := append(game.Computed.ActiveChains, NoHotel)
	combinations := generateCombinations(options)

	actions := make([]gmcts.Action, 0)

	for _, combination := range combinations {

		remainingMoney := game.CurrentPlayer().Money
		remainingHotels := make(map[Hotel]int, 0)

		for _, h := range combination {
			if h != NoHotel {
				remainingHotels[h] = game.Stocks[h.Index()]
			} else {
				remainingHotels[h] = 0
			}
		}

		use := func(hotel Hotel) int {

			if hotel == NoHotel {
				return 0
			}

			value := shareValueCalc(hotel, game.ChainSize[hotel.Index()])
			if remainingMoney < value {
				return 0
			}

			v := util.Min(1, remainingHotels[hotel])

			remainingHotels[hotel] -= 1
			remainingMoney -= value
			return util.Max(v, 0)
		}

		action := Action_PurchaseStock{
			Purchases: [3]StockPurchase{
				{
					Hotel:  combination[0],
					Amount: use(combination[0]),
				},
				{
					Hotel:  combination[1],
					Amount: use(combination[1]),
				},
				{
					Hotel:  combination[2],
					Amount: use(combination[2]),
				},
			},
		}

		// scan through the created purchases array to eliminate the actions which don't do anything
		// created as a local function for code clarity
		isPointless := func() bool {
			hasHotel := false
			totalAmount := 0
			for _, p := range action.Purchases {
				totalAmount += p.Amount
				if p.Hotel != NoHotel {
					hasHotel = true
				}
			}

			if totalAmount < 1 {
				return true
			}

			if hasHotel == false {
				return true
			}

			return false
		}

		if isPointless() == false {
			actions = append(actions, action)
		}
	}

	actions = append([]gmcts.Action{
		Action_PurchaseStock{},
	}, actions...)

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

func (game *Game) end(_ string) {

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

// CanEnd
// returns true if it's possible for a player to 'declare' the game over
func (game *Game) CanEnd() (string, bool) {

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
