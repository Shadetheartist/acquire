package acquire_2

import (
	"sort"
)

type Computed struct {
	AvailableChains []Hotel
	ActiveChains    []Hotel
	LargestChains   []Hotel

	// the current player's legal moves
	LegalMoves []Tile
}

func NewComputed(game *Game) *Computed {
	computed := &Computed{}

	computed.computeChains(game)
	computed.computeLegalMoves(game)

	return computed
}

func (c *Computed) computeChains(game *Game) {
	availableChains := make([]Hotel, 0, NUM_CHAINS/2)
	activeChains := make([]Hotel, 0, NUM_CHAINS/2)
	largestChains := make([]Hotel, 0, len(game.ChainSize))

	for idx, size := range game.ChainSize {
		hotel := ChainFromIdx(idx)
		if size == 0 {
			availableChains = append(availableChains, hotel)
		} else {
			activeChains = append(activeChains, hotel)
		}

		largestChains = append(largestChains, hotel)
	}

	// sorts hotels largest to smallest, by chain size
	sort.Slice(largestChains, func(i, j int) bool {
		return game.ChainSize[largestChains[i].Index()] > game.ChainSize[largestChains[j].Index()]
	})

	c.AvailableChains = availableChains
	c.ActiveChains = activeChains
	c.LargestChains = largestChains

}

// getLargestChainsOf
// returns the hotel(s) with the largest size, and the size of the largest hotel(s)
// the slice is in descending chain size order
func (game *Game) getLargestChainsOf(hotels []Hotel) ([]Hotel, int) {
	if len(hotels) < 1 {
		return []Hotel{}, 0
	}

	type HC struct {
		hotel Hotel
		count int
	}

	sizes := make([]HC, 0)

	for _, h := range hotels {
		sizes = append(sizes, HC{
			hotel: h,
			count: game.ChainSize[h.Index()],
		})
	}

	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i].count > sizes[j].count
	})

	// take all the hotels with the same length as the largest chain
	largestChain := sizes[0]
	largestChains := make([]Hotel, 0)
	for _, hc := range sizes {
		if hc.count == largestChain.count {
			largestChains = append(largestChains, hc.hotel)
		}
	}

	return largestChains, sizes[0].count
}

func (c *Computed) computeLegalMoves(game *Game) {
	legalMoves := make([]Tile, 0, len(game.CurrentPlayer().Tiles))

	for _, t := range game.CurrentPlayer().Tiles {
		if c.isLegalToPlace(game, t) {
			legalMoves = append(legalMoves, t)
		}
	}

	c.LegalMoves = legalMoves
}

func (c *Computed) isLegalToPlace(game *Game, tile Tile) bool {

	if tile == NoTile {
		return false
	}

	pos := tile.Pos()
	neighboringHotels := getNeighbors(game, pos)
	chainsInNeighbors := getChainsInNeighbors(neighboringHotels)

	// this tile would start a merger if placed
	if len(chainsInNeighbors) > 1 {
		// if any two neighbors are size > 10, then the placement isn't legal
		numSafe := 0
		for _, hotel := range chainsInNeighbors {
			size := game.ChainSize[hotel.Index()]
			if size > 11 {
				numSafe += 1
			}

			if numSafe == 2 {
				return false
			}
		}
	}

	// this would found a new chain if placed
	undefinedNeighbors := getUndefinedNeighbors(neighboringHotels)
	if len(undefinedNeighbors) > 0 {
		// if there are no available chains left to create, this move is invalid
		if len(game.Computed.AvailableChains) < 1 {
			return false
		}
	}

	return true
}
