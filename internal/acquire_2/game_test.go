package acquire_2

import (
	"git.sr.ht/~bonbon/gmcts"
	"testing"
)

func TestPlaceTileOnBoard(t *testing.T) {
	tileToPlace := Tile4D
	hotelToPlace := WorldwideHotel
	game := NewGame()
	game.placeTileOnBoard(tileToPlace, hotelToPlace)
	game.Computed = NewComputed(game)

	if game.LastPlacedTile != tileToPlace {
		t.Fatal()
	}

	if game.Board[tileToPlace.Index()].Tile != tileToPlace {
		t.Fatal()
	}

	if game.Board[tileToPlace.Index()].Hotel != hotelToPlace {
		t.Fatal()
	}

	if game.ChainSize[hotelToPlace.Index()] != 1 {
		t.Fatal()
	}
}

func TestPlaceTileActions(t *testing.T) {
	tileA := Tile4D
	tileB := Tile5D

	game := NewGame()

	doAction := func(action gmcts.Action) {
		newGame, err := game.ApplyAction(action)

		if err != nil {
			t.Fatal(err)
		}

		game = newGame.(*Game)
	}

	game.placeTileOnBoard(tileA, UndefinedHotel)

	game.CurrentPlayer().Tiles[0] = tileB

	game.Computed = NewComputed(game)

	doAction(Action_PlaceTile{
		Tile: tileB,
		End:  false,
	})

	// purchase stock action
	actions := game.GetActions()
	doAction(actions[0])

	actions = game.GetActions()
	doAction(actions[0])

	print(actions)

}
