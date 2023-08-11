package acquire

import "testing"

func TestPropagation(t *testing.T) {
	game := NewGame()
	// make a big block of W
	game.placeTileOnBoard(Tile1A, WorldwideHotel)
	game.placeTileOnBoard(Tile2A, WorldwideHotel)
	game.placeTileOnBoard(Tile3A, WorldwideHotel)
	game.placeTileOnBoard(Tile4A, WorldwideHotel)

	game.placeTileOnBoard(Tile1B, WorldwideHotel)
	game.placeTileOnBoard(Tile2B, WorldwideHotel)
	game.placeTileOnBoard(Tile3B, WorldwideHotel)
	game.placeTileOnBoard(Tile4B, WorldwideHotel)

	game.placeTileOnBoard(Tile1C, WorldwideHotel)
	game.placeTileOnBoard(Tile2C, WorldwideHotel)
	game.placeTileOnBoard(Tile3C, WorldwideHotel)
	game.placeTileOnBoard(Tile4C, WorldwideHotel)

	// a gap, then some C

	game.placeTileOnBoard(Tile1E, ContinentalHotel)
	game.placeTileOnBoard(Tile2E, ContinentalHotel)
	game.placeTileOnBoard(Tile3E, ContinentalHotel)
	game.placeTileOnBoard(Tile4E, ContinentalHotel)

	game.placeTileOnBoard(Tile1F, ContinentalHotel)
	game.placeTileOnBoard(Tile2F, ContinentalHotel)
	game.placeTileOnBoard(Tile3F, ContinentalHotel)
	game.placeTileOnBoard(Tile4F, ContinentalHotel)

	game.placeTileOnBoard(Tile1G, ContinentalHotel)
	game.placeTileOnBoard(Tile2G, ContinentalHotel)
	game.placeTileOnBoard(Tile3G, ContinentalHotel)
	game.placeTileOnBoard(Tile4G, ContinentalHotel)

	// calculate the stuff
	game.Computed = NewComputed(game)

	if game.ChainSize[WorldwideHotel.Index()] != 12 {
		t.Fatal("wrong chain size, pre-computation err")
	}

	if game.ChainSize[ContinentalHotel.Index()] != 12 {
		t.Fatal("wrong chain size, pre-computation err")
	}

	Render(game)

	// place a new hotel which would merge
	game.placeTileOnBoard(Tile1D, WorldwideHotel)

	// just placing the tile should increment the chain size
	if game.ChainSize[WorldwideHotel.Index()] != 13 {
		t.Fatal("wrong chain size, placing tile of a chain should increment chain size")
	}

	Render(game)

	// test the propagation
	propagateHotelChain(game, PlacedHotel{
		Hotel: WorldwideHotel,
		Tile:  Tile1D,
	})

	Render(game)

	// after propagating into the new chain, the chain size should be the sum of the two, plus 1 for the merging piece
	if game.ChainSize[WorldwideHotel.Index()] != 25 {
		t.Fatal("wrong chain size, propagating should increment chain size")
	}

	// and the merged chain should be at zero, since there are none left
	if game.ChainSize[ContinentalHotel.Index()] != 0 {
		t.Fatal("wrong chain size, propagating (merging) should end up at zero")
	}

}
