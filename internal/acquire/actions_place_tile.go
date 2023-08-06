package acquire

import (
	"acquire/internal/util"
	"git.sr.ht/~bonbon/gmcts"
)

type Action_PlaceTile struct {
	Tile Tile
	End  bool
}

func (a Action_PlaceTile) Type() ActionType {
	return ActionType_PlaceTile
}

func (game *Game) getPlaceTileActions() []gmcts.Action {
	moves := game.Computed.LegalMoves
	var skip bool
	if len(moves) < 1 {
		skip = refreshOrSkip(game, game.CurrentPlayer(), 1)
	}

	actions := util.Map(moves, func(val Tile) gmcts.Action {
		return Action_PlaceTile{Tile: val}
	})

	if skip {
		actions = append(actions, Action_PlaceTile{Tile: NoTile})
	}

	return actions
}

func (game *Game) applyPlaceTileAction(action Action_PlaceTile) {

	goNext := func() {
		game.NextActionType = ActionType_PurchaseStock
	}

	// player wants to skip their turn
	if action.Tile == NoTile {
		game.SkippedTurnsInARow++

		if game.SkippedTurnsInARow > len(game.Players) {
			game.end("no one had any moves left to play")
			return
		}

		goNext()
		return
	}

	game.SkippedTurnsInARow = 0

	player := game.CurrentPlayer()
	tile := action.Tile

	err := player.removeTileFromHand(tile)
	if err != nil {
		panic(err)
	}

	pos := tile.Pos()

	neighboringHotels := getNeighbors(game, pos)

	game.placeTileOnBoard(tile, UndefinedHotel)

	// no neighbors - no effects, go to next player's turn
	if !hasNeighboringHotel(neighboringHotels) {
		goNext()
		return
	}

	// growing a chain - occurs when, of all neighbors, there is only one type of hotel
	chainsInNeighbors := getChainsInNeighbors(neighboringHotels)
	if len(chainsInNeighbors) == 1 {
		hotel := chainsInNeighbors[0]
		game.placeTileOnBoard(tile, hotel)

		goNext()
		return
	}

	// merger - if there are more than two chains in the neighboring tiles, a merger must take place
	if len(chainsInNeighbors) > 1 {

		largestChains, _ := game.getLargestChainsOf(chainsInNeighbors)

		game.MergerState = MergerState{
			ChainsToMerge:    [7]int{},
			MergingPlayerIdx: game.playerTurn(0),
			AcquiringHotel:   largestChains[0], //select the largest chain by default
		}

		// more than one chain is tied for largest, player needs to decide which chain is acquired
		if len(largestChains) > 1 {
			game.NextActionType = ActionType_PickHotelToMerge
			return
		}

		// prepare the 'chains to merge' array
		for _, h := range HotelChainList {
			// ok = this hotel is in the 'largest chains' slice, but isn't the largest chain
			_, ok := util.IndexOf(largestChains[:1], h)
			if ok {
				game.MergerState.ChainsToMerge[h.Index()] = len(game.Players)
			}
		}

		// otherwise...

		game.NextActionType = ActionType_Merge

		return
	}

	// found a new chain - occurs when a tile has one or more neighbors which are all still undefined
	undefinedNeighbors := getUndefinedNeighbors(neighboringHotels)
	if len(undefinedNeighbors) > 0 {
		game.NextActionType = ActionType_PickHotelToFound
		game.FoundingHotel = NoHotel
		return
	}

	panic("unexpectedly got here")
}

// getChainNeighbors
// returns a slice of the hotels around a position which are not NoHotel nor UndefinedHotel
func getChainNeighbors(neighbors []PlacedHotel) []PlacedHotel {
	chainNeighbors := make([]PlacedHotel, 0)

	for _, d := range util.Directions {
		if neighbors[d].Hotel != NoHotel && neighbors[d].Hotel != UndefinedHotel {
			chainNeighbors = append(chainNeighbors, neighbors[d])
		}
	}

	return chainNeighbors
}

// getChainsInNeighbors
// returns a slice of unique hotel chains which are in the neighbors slice
func getChainsInNeighbors(neighbors []PlacedHotel) []Hotel {
	chainNeighbors := getChainNeighbors(neighbors)
	hotels := util.Map(chainNeighbors, func(val PlacedHotel) Hotel {
		return val.Hotel
	})
	return util.UniqueElements[Hotel](hotels)
}

// hasNeighboringHotel
// returns true if all the hotels in the slice are 'NoHotel'
func hasNeighboringHotel(neighbors []PlacedHotel) bool {
	for _, d := range util.Directions {
		if neighbors[d].Hotel != NoHotel {
			return true
		}
	}
	return false
}

func getUndefinedNeighbors(neighbors []PlacedHotel) []PlacedHotel {
	undefinedNeighbors := make([]PlacedHotel, 0)

	for _, d := range util.Directions {
		if neighbors[d].Hotel == UndefinedHotel {
			undefinedNeighbors = append(undefinedNeighbors, neighbors[d])
		}
	}

	return undefinedNeighbors
}

// refreshOrSkip
// this function will refresh the tiles of a player if they have no legal moves repeatedly n times
// if the player doesn't have a valid move after a refresh then their turn should be skipped
// (as indicated by true in the returned bool)
func refreshOrSkip(game *Game, p *Player, n int) bool {
	for i := 0; i < n; i++ {
		if len(game.Computed.LegalMoves) < 1 {
			p.refreshTiles(game)
		} else {
			return false
		}
	}

	return true
}
