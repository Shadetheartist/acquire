package acquire

import (
	"acquire/internal/util"
	"fmt"
	"os"
	"sort"
)

type Game struct {
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
		NewPlayer(game, "You", aiAgentFactory),
		NewPlayer(game, "Greg", aiAgentFactory),
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

func (game *Game) Step() {

	if game.IsOver {
		return
	}

	player := game.CurrentPlayer()

	// just place the first tile for each other player for testing
	tile, err := player.agent.DetermineTilePlacement()

	// if the tile is NoTile, skip their turn
	if tile == NoTile {
		if !game.areThereAnyLegalMoves() {
			game.end("there are no remaining legal moves in the game")
			return
		}
		game.Turn++
		return
	}

	if err != nil {
		game.abort(err.Error())
	}

	player.placeTile(tile)

	if reason, ok := game.canEnd(); ok {
		end, err := player.agent.DetermineGameEnd()
		if err != nil {
			game.abort(err.Error())
		}

		if end {
			game.end("game has been declared over, " + reason)
			return
		}
	}

	// buy stocks at this point
	hotel, n, err := player.agent.DetermineStockPurchase()
	if err != nil {
		panic(err)
	}

	if isActualHotelChain(hotel) {
		err = player.buyStock(hotel, n)
		if err != nil {
			panic(err)
		}
	}

	if len(player.Inventory.Tiles.Items) < 6 {
		// take a new tile from the bank
		// ignoring the error, if there weren't any tiles left to take
		_ = player.Inventory.Tiles.take(game.Inventory.Tiles)
	}

	game.Turn++
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

	sortedPlayers := make([]*Player, len(game.Players))
	copy(sortedPlayers, game.Players)
	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].Inventory.Money > sortedPlayers[j].Inventory.Money
	})

	sendMsg(game, "Game over: "+reason)

	// send winners
	sendMsg(game, "Result")
	for i, p := range sortedPlayers {
		sendMsg(game, fmt.Sprintf("%d: %s with $%d", i+1, p.Name(), p.Inventory.Money))
	}

}

// canEnd
// returns true if it's possible for a player to 'declare' the game over
func (game *Game) canEnd() (string, bool) {

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
	anyUnsafe := util.AnyInMap(sizes, func(key Hotel, val int) bool {
		return val <= 10
	})

	if !anyUnsafe {
		return "all chains on the board are safe", true
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
	availableChains := getAvailableHotelChains(game)
	stocks := game.remainingStocks()

	// if a hotel is not yet placed on the board
	// it cannot be bought, so the stock amount would be zero
	for _, hotel := range availableChains {
		stocks[Stock(hotel)] = 0
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

	sendMsg(game, fmt.Sprintf("%s got a major shareholder bonus of $%d", majShareholder.PlayerName, majBonus))
	sendMsg(game, fmt.Sprintf("%s got a minor shareholder bonus of $%d", minorShareholder.PlayerName, minBonus))
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
