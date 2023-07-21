package acquire

import (
	"acquire/internal/util"
	"fmt"
	"os"
)

type Game struct {
	MergerState *MergerState
	FoundState  FoundState

	// once set, at the end of the turn, the game ends
	WillEnd bool

	NextActionType ActionType

	LastPlacedTile Tile

	IsOver         bool
	Turn           int
	Players        []*Player
	Board          *Board
	Inventory      *Inventory
	inputInterface IInput
}

func NewGame(inputInterface IInput) *Game {

	game := &Game{
		inputInterface: inputInterface,
		Board:          newBoard(),
	}

	game.Players = []*Player{
		NewPlayer(game, 0, "You", aiAgentStupidFactory),
		NewPlayer(game, 1, "Jef", aiAgentStupidFactory),
		NewPlayer(game, 2, "Jame", aiAgentStupidFactory),
		NewPlayer(game, 3, "Eric", aiAgentStupidFactory),
	}

	inventory := newInventory(game, 1e6)
	for _, hotel := range HotelChainList {
		inventory.Stocks[hotel].add(25, func() Stock {
			return Stock(hotel)
		})
	}

	// hard-set the tiles to this set of randomized tiles
	inventory.Tiles.Items = randomizedTiles()

	game.Inventory = inventory

	// setup

	for _, p := range game.Players {
		p.Inventory.takeMoney(game.Inventory, 6000)

		// this cannot fail yet
		_ = p.takeTiles(game.Inventory, 6)
	}

	return game
}

func (game *Game) Name() string {
	return "The Acquire Bank"
}

func (game *Game) CurrentPlayer() *Player {
	return game.Players[game.playerTurn(0)]
}

func (game *Game) abort(reason string) {
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
			err := p.sellStock(stock, p.remainingStock(stock))
			if err != nil {
				panic(err)
			}
		}
	}
}

// canEnd
// returns true if it's possible for a player to 'declare' the game over
func (game *Game) canEnd() (string, bool) {

	if len(game.Inventory.Tiles.Items) < 1 {
		return "no tiles left somehow", true
	}

	// if there are any chains larger than 40, the game can end
	sizes := countHotelChains(game, HotelChainList)
	anyLargeEnough := util.AnyInMap(sizes, func(key Hotel, val int) bool {
		return val > 40
	})

	if anyLargeEnough {
		return "there is a chain with length 41", true
	}

	// if there are no remaining unsafe (from merger) hotels, the game can end
	sizes = countHotelChains(game, getActiveHotelChains(game))

	if len(sizes) > 1 {
		anyUnsafe := util.AnyInMap(sizes, func(key Hotel, val int) bool {
			return val <= 10
		})

		if !anyUnsafe {
			return "all chains on the board are safe", true
		}
	}

	return "", false
}

// areThereAnyLegalMoves
// if there are no legal moves left in all tiles, bank and player inventories combined
// this is a special case for the stupid as fuck AI and probably would never happen in a real game,
// hence the lack of stated rule for this outcome in the rulebook
func (game *Game) areThereAnyLegalMoves() bool {
	everyRemainingTile := game.everyRemainingTile()
	for _, t := range everyRemainingTile {
		if isLegalToPlace(game, t) {
			return true
		}
	}
	return false
}

// playersInTurnOrder
// this returns a slice of players in the order which starts from
// the current player and includes each player, in turn order.
func (game *Game) playersInTurnOrder() []*Player {
	players := make([]*Player, 0, len(game.Players))
	for i := 0; i < len(game.Players); i++ {
		players = append(players, game.Players[game.playerTurn(i)])
	}
	return players
}

// playerTurn
// returns the index of the player whos turn it is (if the offset is zero).
// supply a non-zero offset to return the index of the play whos turn it will be after 'offset' turns.
func (game *Game) playerTurn(offset int) int {
	if offset < 0 {
		panic("offset cannot be < 0")
	}
	return (game.Turn + offset) % len(game.Players)
}

func (game *Game) remainingStock(stock Stock) int {
	return len(game.Inventory.Stocks[Hotel(stock)].Items)
}

func (game *Game) remainingStocks() map[Stock]int {
	stocks := make(map[Stock]int)

	for _, h := range HotelChainList {
		stocks[Stock(h)] = game.remainingStock(Stock(h))
	}

	return stocks
}

func (game *Game) purchasableStocks() map[Stock]int {
	availableChains := getActiveHotelChains(game)
	stocks := game.remainingStocks()

	// if a hotel is not yet placed on the board
	// it cannot be bought, so remove it from contention
	for _, hotel := range availableChains {
		delete(stocks, Stock(hotel))
	}

	return stocks
}

func (game *Game) payShareholderBonuses(hotel Hotel) {
	size := countHotelChain(game, hotel)
	majShareholder, majShareholderShares, minorShareholder, minorShareholderShares := Stock(hotel).getShareholders(game)

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

	majShareholder.Inventory.takeMoney(game.Inventory, majBonus)
	minorShareholder.Inventory.takeMoney(game.Inventory, minBonus)
}

// everyRemainingTile
// tiles in the bank + the tiles in each players inventory
func (game *Game) everyRemainingTile() []Tile {
	tiles := make([]Tile, 0)

	tiles = append(tiles, game.Inventory.Tiles.Items...)

	for _, p := range game.Players {
		tiles = append(tiles, p.Inventory.Tiles.Items...)
	}

	return tiles
}
