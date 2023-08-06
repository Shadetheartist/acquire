package acquire

import (
	"testing"
)

func TestComparison(t *testing.T) {
	game1 := NewGame()
	game2 := game1

	if game1 != game2 {
		t.Fatal("games are not equal wtf")
	}

	game2.Players[0].Money += 1

	if game1 != game2 {
		t.Fatal("games are still equal wtf")
	}

	game2 = game1

	if game1 != game2 {
		t.Fatal("games are not equal wtf")
	}

	game2.Board[0] = PlacedHotel{
		Hotel: AmericanHotel,
		Tile:  Tile1A,
	}

	if game1 != game2 {
		t.Fatal("games are still equal wtf")
	}

}
