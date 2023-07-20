package acquire

import (
	"acquire/internal/util"
	"sort"
)

// getLargestChainsOf
// returns a slice of each chain tied for the largest hotel chain length,
// normally consisting of a single element (largest chain, no tie)
func getLargestChainsOf(game *Game, hotels []Hotel) ([]Hotel, int) {
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
			count: countHotelChain(game, h),
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

// countHotelChains
// returns a map of each hotel chain and the current size of the chain
func countHotelChains(game *Game, hotels []Hotel) map[Hotel]int {
	counts := make(map[Hotel]int)

	for _, h := range hotels {
		counts[h] = countHotelChain(game, h)
	}

	return counts
}

// countHotelChain
// this returns the number of tiles comprising a chain of hotels, based on a specified hotel chain type
// if there are no hotels of that chain on the board, the returned value has hotel = NoHotel
func countHotelChain(game *Game, hotel Hotel) int {
	placedHotel := findPlacedHotelOnBoard(game, hotel)
	if placedHotel.Hotel == NoHotel {
		return 0
	}

	return countConnectedPlacedHotels(game, placedHotel)
}

// findPlacedHotelOnBoard
// this returns the first PlacedHotel it encounters on the board of the chain specified
func findPlacedHotelOnBoard(game *Game, hotel Hotel) PlacedHotel {
	return game.Board.Matrix.Find(func(rt PlacedHotel, x int, y int, idx int) bool {
		return rt.Hotel == hotel
	})
}

// countConnectedPlacedHotels
// this breadth-first searches the matrix for connected hotels of the same chain and returns
// the size of the hotel chain
func countConnectedPlacedHotels(game *Game, rootHotel PlacedHotel) int {
	count := 0

	visited := make([]bool, game.Board.Matrix.Area())
	stack := make([]PlacedHotel, 0, 8)

	// root stack push
	stack = append(stack, rootHotel)
	var hotel PlacedHotel

	for len(stack) > 0 {
		// pop off stack
		hotel, stack = stack[len(stack)-1], stack[:len(stack)-1]
		idx := game.Board.Matrix.IndexPt(hotel.Pos)

		if hotel.Hotel == NoHotel {
			continue
		}

		//avoid revisiting
		if visited[idx] {
			continue
		}

		//keep track of what has been visited already
		visited[idx] = true

		// only track the hotel we're propagating
		if hotel.Hotel != rootHotel.Hotel {
			continue
		}

		count++

		// add neighbors to stack
		neighborPts := hotel.Pos.OrthogonalNeighbours()
		neighbors := util.Map(neighborPts, func(val util.Point[int]) PlacedHotel {
			h, err := game.Board.Matrix.GetPt(val)
			if err != nil {
				// on out of bounds we can map to the current element and it will get skipped
				return PlacedHotel{}
			}
			return h
		})

		stack = append(stack, neighbors...)

	}

	return count
}

// propagateHotelChain
// this breadth-first searches the matrix for connected hotels of
// any chain and converts them to the type of the rootHotel
func propagateHotelChain(game *Game, rootHotel PlacedHotel) {
	visited := make([]bool, game.Board.Matrix.Area())
	stack := make([]PlacedHotel, 0, 8)

	// root stack push
	stack = append(stack, rootHotel)
	var hotel PlacedHotel

	for len(stack) > 0 {
		// pop off stack
		hotel, stack = stack[len(stack)-1], stack[:len(stack)-1]
		idx := game.Board.Matrix.IndexPt(hotel.Pos)

		if hotel.Hotel == NoHotel {
			continue
		}

		//avoid revisiting
		if visited[idx] {
			continue
		}

		//keep track of what has been visited already
		visited[idx] = true

		// the hotel we're already propagating
		if hotel.Hotel == rootHotel.Hotel && hotel != rootHotel {
			continue
		}

		game.Board.Matrix.Set(hotel.Pos.X, hotel.Pos.Y, PlacedHotel{
			Hotel: rootHotel.Hotel,
			Pos:   hotel.Pos,
		})

		// add neighbors to stack
		neighborPts := hotel.Pos.OrthogonalNeighbours()
		neighbors := util.Map(neighborPts, func(val util.Point[int]) PlacedHotel {
			h, err := game.Board.Matrix.GetPt(val)
			if err != nil {
				// on out of bounds we can map to the current element and it will get skipped
				return PlacedHotel{}
			}
			return h
		})

		stack = append(stack, neighbors...)

	}
}
