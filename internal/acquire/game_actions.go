package acquire

import (
	"acquire/internal/util"
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

	// ActionType_EndGame this doesn't indicate that the game is going to end,
	//  it just gives the player the opportunity to end the game, since it's optional
	ActionType_EndGame
)

type IAction interface {
	Type() ActionType
}

// PlaceTile

type Action_PlaceTile struct {
	Tile Tile
}

func (a Action_PlaceTile) Type() ActionType {
	return ActionType_PlaceTile
}

// PickHotelToFound

type Action_PickHotelToFound struct {
	Hotel Hotel
}

func (a Action_PickHotelToFound) Type() ActionType {
	return ActionType_PickHotelToFound
}

// PickHotelToMerge

type Action_PickHotelToMerge struct {
	Hotel Hotel
}

func (a Action_PickHotelToMerge) Type() ActionType {
	return ActionType_PickHotelToMerge
}

// Merge

type MergeSubAction struct {
	MergeType MergerAction
	Amount    int
}

type Action_Merge struct {
	Actions [3]MergeSubAction
}

func (a Action_Merge) Type() ActionType {
	return ActionType_Merge
}

// PurchaseStock

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

// EndGame

type Action_EndGame struct {
	end bool
}

func (a Action_EndGame) Type() ActionType {
	return ActionType_EndGame
}

func (game *Game) applyPlaceTileAction(action Action_PlaceTile) {

	goNext := func() {
		game.NextActionType = ActionType_EndGame
	}

	// player wants to skip their turn
	if action.Tile == NoTile {
		goNext()
		return
	}

	player := game.CurrentPlayer()
	tile := action.Tile

	player.Inventory.Tiles.remove(tile)
	pos := tile.Pos()
	matrix := player.game.Board.Matrix
	matrix.Set(pos.X, pos.Y, PlacedHotel{Pos: pos, Hotel: UndefinedHotel})
	game.LastPlacedTile = tile

	// what's the next state?
	// depends on the players next action
	neighboringHotels := matrix.GetNeighbors(pos)

	// no neighbors - no effects, go to next player's turn
	if !hasNeighboringHotel(neighboringHotels) {
		goNext()
		return
	}

	// growing a chain - occurs when, of all neighbors, there is only one type of hotel
	chainsInNeighbors := getChainsInNeighbors(neighboringHotels)
	if len(chainsInNeighbors) == 1 {
		hotel := chainsInNeighbors[0]
		newPlacedHotel := PlacedHotel{Pos: pos, Hotel: hotel}
		matrix.Set(pos.X, pos.Y, newPlacedHotel)

		propagateHotelChain(player.game, newPlacedHotel)

		goNext()
		return
	}

	// merger - if there are more than two chains in the neighboring tiles, a merger must take place
	if len(chainsInNeighbors) > 1 {

		largestChains, _ := getLargestChainsOf(player.game, chainsInNeighbors)

		// count the hotel chains *before* the tile is placed for accurate stock purchasing
		chainSizeMap := countHotelChains(player.game, chainsInNeighbors)

		// in order of largest to smallest
		sortedHotelChains := sortChainSizeMap(chainSizeMap)

		game.MergerState = &MergerState{
			Pos:               pos,
			LargestChains:     largestChains,
			NeighboringHotels: neighboringHotels,
			ChainsInNeighbors: chainsInNeighbors,
			AcquiringHotel:    largestChains[0],
			// leave all chains in for now, we're not sure which one the player will select, so we remove it later
			RemainingChainsToMerge:  sortedHotelChains,
			RemainingPlayersToMerge: map[HotelChain][]*Player{},
		}

		if len(largestChains) > 1 {
			game.NextActionType = ActionType_PickHotelToMerge
			return
		}

		// remove the selected chain which is the largest, when there's no conflict
		game.MergerState.RemainingChainsToMerge = util.RemoveAt(game.MergerState.RemainingChainsToMerge, 0)

		// the players will proceed in turn order for each chain to merge, in order of largest to smallest chain
		for _, c := range game.MergerState.RemainingChainsToMerge {
			game.MergerState.RemainingPlayersToMerge[c] = game.playersInTurnOrder()
		}

		game.NextActionType = ActionType_Merge

		return
	}

	// found a new chain - occurs when a tile has one or more neighbors which are all still undefined
	undefinedNeighbors := getUndefinedNeighbors(neighboringHotels)
	if len(undefinedNeighbors) > 0 {
		game.NextActionType = ActionType_PickHotelToFound
		game.FoundState = FoundState{
			FoundingHotel: NoHotel,
			Pos:           pos,
		}
		return
	}

	panic("unexpectedly got here")
}

func (game *Game) applyPickHotelToFoundAction(action Action_PickHotelToFound) {
	pos := game.FoundState.Pos

	newPlacedHotel := PlacedHotel{Pos: pos, Hotel: action.Hotel}
	game.Board.Matrix.Set(pos.X, pos.Y, newPlacedHotel)

	propagateHotelChain(game, newPlacedHotel)

	// receive a free stock in the chain for founding it
	_ = game.CurrentPlayer().Inventory.Stocks[action.Hotel].take(game.Inventory.Stocks[action.Hotel])

	game.NextActionType = ActionType_EndGame
	// game.FoundState = nil

}

func (game *Game) applyPickHotelToMergeAction(action Action_PickHotelToMerge) {
	game.MergerState.AcquiringHotel = action.Hotel

	// remove the acquiring hotel chain from the list to merge
	idx, _ := util.Find(game.MergerState.RemainingChainsToMerge, func(val HotelChain) bool {
		return val.Hotel == action.Hotel
	})

	game.MergerState.RemainingChainsToMerge = util.RemoveAt(game.MergerState.RemainingChainsToMerge, idx)

	// the players will proceed in turn order for each chain to merge, in order of largest to smallest chain
	for _, c := range game.MergerState.RemainingChainsToMerge {
		game.MergerState.RemainingPlayersToMerge[c] = game.playersInTurnOrder()
	}

	game.NextActionType = ActionType_Merge
}

func (game *Game) applyMerge(action Action_Merge) {

	state := game.MergerState

	chainToMerge := state.RemainingChainsToMerge[0]
	player := state.RemainingPlayersToMerge[chainToMerge][0]
	matrix := game.Board.Matrix
	pos := state.Pos

	game.payShareholderBonuses(chainToMerge.Hotel)

	goNext := func() {
		// pop out the player we processed
		_, state.RemainingPlayersToMerge[chainToMerge] = util.Pop(state.RemainingPlayersToMerge[chainToMerge])

		// if there's no more players to process for this merge, pop the chain we finished merging
		if len(state.RemainingPlayersToMerge[chainToMerge]) < 1 {
			_, state.RemainingChainsToMerge = util.Pop(state.RemainingChainsToMerge)
		}

		// if there's no more chains to process, we're done merging
		l := len(state.RemainingChainsToMerge)
		if l < 1 {
			// finally place the piece and propagate the chain
			newPlacedHotel := PlacedHotel{Pos: pos, Hotel: state.AcquiringHotel}
			matrix.Set(pos.X, pos.Y, newPlacedHotel)
			propagateHotelChain(game, newPlacedHotel)

			state = nil

			game.NextActionType = ActionType_EndGame
		}
	}

	// i've set this up so that all the information for a player's merger is provided at once
	// since they can choose multiple things to do, this loop represents those things.
	// for instance, you could have 6 stocks, trade in 4, and sell two - much depends on the action options provided
	for _, subAction := range action.Actions {
		switch subAction.MergeType {

		case Hold:
			goNext()
			return

		case Trade:
			err := player.tradeIn(Stock(state.AcquiredHotel), Stock(state.AcquiringHotel), subAction.Amount)
			if err != nil {
				panic(err)
			}
			break

		case Sell:
			err := player.sellStock(Stock(state.AcquiredHotel), subAction.Amount)
			if err != nil {
				panic(err)
			}
			break
		}
	}

	goNext()

}

func (game *Game) applyEndGameAction(action Action_EndGame) {
	if action.end {
		game.WillEnd = true
	}

	game.NextActionType = ActionType_PurchaseStock
}

func (game *Game) applyPurchaseStockAction(action Action_PurchaseStock) {

	for _, purchase := range action.Purchases {
		if purchase.Hotel == NoHotel || purchase.Amount == 0 {
			continue
		}

		err := game.CurrentPlayer().buyStock(purchase.Hotel, purchase.Amount)
		if err != nil {
			panic(err)
		}
	}

	// going to merge this and the draw tile 'action' for better efficiency and it's easier

	// take a new tile from the bank
	// ignoring the error, if there weren't any tiles left to take
	_ = game.CurrentPlayer().Inventory.Tiles.take(game.Inventory.Tiles)

	// game always ends at the end of the player's turn
	if game.WillEnd {
		game.end("a player has declared the game is over")
		return
	}

	game.NextActionType = ActionType_PlaceTile
	game.Turn++
}
