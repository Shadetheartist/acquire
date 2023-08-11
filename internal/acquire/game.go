package acquire

import (
	"acquire/internal/util"
	"git.sr.ht/~bonbon/gmcts"
	"sort"
)

const BOARD_MAX_X = 12
const BOARD_MAX_Y = 9

const MAX_PLAYERS = 2
const MAX_TILES_IN_HAND = 6
const NUM_CHAINS = 7

type PlacedHotel struct {
	Hotel Hotel
	Tile  Tile
}

// MergerState
// as multiple actions occur during a merger, and the state of the board matters
// we have to the state of the merger over different turns and actions
// nil when not in use
type MergerState struct {
	// mapped to hotel
	// the value is the number of players remaining to merge (3, 2, 1, 0 = done)
	ChainsToMerge    [NUM_CHAINS]int
	MergingPlayerIdx int
	AcquiringHotel   Hotel
}

type Game struct {
	Players        [MAX_PLAYERS]Player
	NextActionType ActionType

	Turn               int
	SkippedTurnsInARow int
	IsOver             bool
	WillEnd            bool

	LastPlacedTile Tile

	Board [BOARD_MAX_X * BOARD_MAX_Y]PlacedHotel
	Tiles [BOARD_MAX_X * BOARD_MAX_Y]Tile

	// index by hotel
	ChainSize [NUM_CHAINS]int
	Stocks    [NUM_CHAINS]int

	// used when founding a hotel
	FoundingHotel Hotel
	Pos           util.Point[int]

	MergerState MergerState

	Computed *Computed
}

func NewGame() *Game {

	game := &Game{}
	game.Tiles = randomizedTiles()

	game.Players = [MAX_PLAYERS]Player{}
	for i := 1; i <= MAX_PLAYERS; i++ {
		game.Players[i-1] = Player{
			Id:     i,
			Money:  6000,
			Tiles:  [MAX_TILES_IN_HAND]Tile{},
			Stocks: [NUM_CHAINS]int{},
		}
	}

	for i := 0; i < NUM_CHAINS; i++ {
		game.Stocks[i] = 25
	}

	for idx := range game.Players {
		// this cannot fail yet
		for i := 0; i < MAX_TILES_IN_HAND; i++ {
			err := game.Players[idx].takeTileFromBank(game)
			if err != nil {
				panic(err)
			}
		}
	}

	game.Computed = NewComputed(game)

	return game
}

func (game *Game) Player() gmcts.Player {
	return gmcts.Player(game.CurrentPlayer().Id)
}

func (game *Game) IsTerminal() bool {
	return game.IsOver
}

func (game *Game) Winners() []gmcts.Player {
	playerSlice := game.PlayerSlice()

	sort.SliceStable(playerSlice, func(i, j int) bool {
		return playerSlice[i].NetWorth(game) > playerSlice[j].NetWorth(game)
	})

	outSlice := make([]gmcts.Player, 0)

	highestMoney := playerSlice[0].NetWorth(game)
	for _, p := range playerSlice {
		if p.NetWorth(game) == highestMoney {
			outSlice = append(outSlice, gmcts.Player(p.Id))
		} else {
			break
		}
	}

	return outSlice
}

// playerTurn
// returns the index of the player whose turn it is (if the offset is zero).
// supply a non-zero offset to return the index of the play whose turn it will be after 'offset' turns.
func (game *Game) playerTurn(offset int) int {
	return (game.Turn + offset) % len(game.Players)
}

// CurrentPlayer
// the player whose turn it is (does not account for merger, use ActivePlayer forthat)
func (game *Game) CurrentPlayer() *Player {
	return &game.Players[game.playerTurn(0)]
}

// NextPlayer
// the player whose turn it is after the current turn
func (game *Game) NextPlayer() *Player {
	return &game.Players[game.playerTurn(1)]
}

// ActivePlayer
// this is just the current player, but if there is a merger happening, it is the player taking the merge action
func (game *Game) ActivePlayer() *Player {
	if game.NextActionType == ActionType_Merge {
		return &game.Players[game.MergerState.MergingPlayerIdx]
	}

	return game.CurrentPlayer()
}

// PlayerSlice
// a slice containing each player in the game
func (game *Game) PlayerSlice() []Player {
	playerSlice := make([]Player, len(game.Players))

	for i := 0; i < len(game.Players); i++ {
		playerSlice[i] = game.Players[i]
	}

	return playerSlice
}

// payShareholderBonuses
// calculates and pays out shareholder bonuses to each player
// the majority shareholder gets the major bonus, the minority shareholder gets the minor bonus
// if there is no minor shareholder the major shareholder gets both bonuses combined
func (game *Game) payShareholderBonuses(hotel Hotel) {
	size := game.ChainSize[hotel.Index()]
	majShareholder, majShareholderShares, minorShareholder, minorShareholderShares := game.getShareholders(Stock(hotel))

	// if majShareholder is nil, no player held shares in this hotel
	if majShareholder == nil {
		return
	}

	majBonus, minBonus := shareholderBonusCalc(
		hotel,
		size,
		majShareholder,
		majShareholderShares,
		minorShareholder,
		minorShareholderShares,
	)

	majShareholder.Money += majBonus
	minorShareholder.Money += minBonus
}

// getShareholders
// returns the major and minor shareholders, and how many shares they have
// if there is no shareholder it will return nil for both players
func (game *Game) getShareholders(s Stock) (*Player, int, *Player, int) {

	var majorShareholder *Player
	var majorShareholderShares int

	var minorShareholder *Player
	var minorShareholderShares int

	for i := range game.Players {
		h := Hotel(s)
		numShares := game.Players[i].Stocks[h.Index()]
		if numShares > majorShareholderShares {
			majorShareholder = &game.Players[i]
			majorShareholderShares = numShares
		} else if numShares > minorShareholderShares {
			minorShareholder = &game.Players[i]
			minorShareholderShares = numShares
		}
	}

	// if major shareholder is still null at this point, then there weren't any players holding shares in this chain
	if majorShareholder == nil {
		return nil, 0, nil, 0
	}

	// if there is no minor shareholder, the major shareholder becomes the major AND minor shareholder
	if minorShareholder == nil {
		minorShareholder = majorShareholder
		minorShareholderShares = majorShareholderShares
	}

	return majorShareholder, majorShareholderShares, minorShareholder, minorShareholderShares
}

func (game *Game) PlacementAtPos(pt util.Point[int]) PlacedHotel {
	if !isInBounds(pt.X, pt.Y) {
		return PlacedHotel{
			Hotel: NoHotel,
			Tile:  NoTile,
		}
	}

	idx := index(pt.X, pt.Y)

	return game.Board[idx]
}

func (game *Game) placeTileOnBoard(tile Tile, hotel Hotel) PlacedHotel {
	newPlacedHotel := PlacedHotel{Tile: tile, Hotel: hotel}
	game.Board[tile.Index()] = newPlacedHotel
	game.LastPlacedTile = tile

	if hotel != NoHotel && hotel != UndefinedHotel {
		game.modifyChainSize(hotel, 1)
	}

	return newPlacedHotel
}

func (game *Game) GetPlayerById(id int) *Player {
	for _, p := range game.Players {
		if p.Id == id {
			return &p
		}
	}

	return nil
}

// modifyChainSize
// just want to keep track of these in a func, easier to track down usage
func (game *Game) modifyChainSize(hotel Hotel, amount int) {
	game.ChainSize[hotel.Index()] += amount
}
