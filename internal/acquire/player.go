package acquire

import (
	"acquire/internal/util"
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

type Player struct {
	Id         int
	PlayerName string
	Game       *Game
	Inventory  *Inventory
}

func (p *Player) clone() *Player {
	clone := &Player{
		// will copy naturally
		Id:         p.Id,
		PlayerName: p.PlayerName,

		// needs to be cloned
		Inventory: p.Inventory.clone(),

		// these don't contribute to state
		Game: p.Game,
	}

	// update refs
	clone.Inventory.Owner = clone

	return clone
}

func (p *Player) Name() string {
	return p.PlayerName
}

func NewPlayer(game *Game, id int, name string) *Player {
	player := &Player{
		Id:         id,
		PlayerName: name,
		Game:       game,
	}

	player.Inventory = newInventory(player, 0)

	return player
}

func (p *Player) pay(amount int) error {
	if p.Inventory.Money < amount {
		return fmt.Errorf("player '%s' cannot afford to pay $%d", p.PlayerName, amount)
	}

	p.Game.Inventory.takeMoney(p.Inventory, amount)

	return nil
}

func (p *Player) buyStock(hotel Hotel, amount int) error {

	if !isActualHotelChain(hotel) {
		panic("trying to buy a non-chain hotel")
	}

	chainSize := countHotelChain(p.Game, hotel)

	if chainSize == 0 {
		return errors.New("chain does not exist yet")
	}

	cost := sharesCalc(hotel, chainSize, amount)

	err := p.pay(cost)
	if err != nil {
		return err
	}

	amountAvailable := len(p.Game.Inventory.Stocks[hotel].Items)
	if amountAvailable < amount {
		return fmt.Errorf("can't buy %d stocks in %s, there's only %d remaining", amount, hotel.String(), amountAvailable)
	}

	amountTaken := 0
	for i := 0; i < amount; i++ {
		err := p.Inventory.Stocks[hotel].take(p.Game.Inventory.Stocks[hotel])
		if err != nil {
			return err
		}
		amountTaken++
	}

	return nil
}

// anon func for the return syntax
func hasNeighboringHotel(neighbors []PlacedHotel) bool {
	for _, d := range util.Directions {
		if neighbors[d].Hotel != NoHotel {
			return true
		}
	}
	return false
}

func (p *Player) legalMoves() []Tile {
	return util.Filter(p.Inventory.Tiles.Items, func(val Tile) bool {
		return isLegalToPlace(p.Game, val)
	})
}

func (p *Player) tiles() []Tile {
	return p.Inventory.Tiles.Items
}

// returnTile
// aliasing function for returning a specific tile to the Game inventory
func (p *Player) returnTile(t Tile) {

	_, ok := p.Inventory.Tiles.indexOf(t)
	if !ok {
		return
	}

	_ = p.Game.Inventory.Tiles.take(p.Inventory.Tiles)
}

// takeTile
// aliasing function for taking a tile from the Game inventory
// this is normally a blind selection
func (p *Player) takeTile() {
	_ = p.Inventory.Tiles.take(p.Game.Inventory.Tiles)
}

func (p *Player) takeTiles(other *Inventory, amount int) error {
	for i := 0; i < amount; i++ {
		err := p.Inventory.Tiles.take(other.Tiles)
		if err != nil {
			return err
		}
	}

	return nil
}

// refreshTiles
// when a player has no legal moves left to play, they can refresh their hand with this func
// puts all tiles back in the Game inv, then takes 6 new ones
func (p *Player) refreshTiles() {

	// clone this so that it isn't getting fucked with as it removes items
	items := util.Clone(p.Inventory.Tiles.Items)
	for _, t := range items {
		p.returnTile(t)
	}

	// shuffle the tiles
	gameTiles := p.Game.Inventory.Tiles.Items
	rand.Shuffle(len(gameTiles), func(i, j int) {
		gameTiles[i], gameTiles[j] = gameTiles[j], gameTiles[i]
	})

	for i := 0; i < 6; i++ {
		p.takeTile()
	}
}

// refreshOrSkip
// this function will refresh the tiles of a player if they have no legal moves repeatedly n times
// if the player doesn't have a valid move after a refresh then their turn should be skipped
// (as indicated by true in the returned bool)
func (p *Player) refreshOrSkip(n int) bool {
	for i := 0; i < n; i++ {
		if len(p.legalMoves()) < 1 {
			p.refreshTiles()
		} else {
			return false
		}
	}

	return true
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

// getActiveHotelChains
// gets the active hotels, which are on the board
// this is pretty brutally inefficient but whatever
func getActiveHotelChains(game *Game) []Hotel {
	chainsMap := make(map[Hotel]struct{}, 0)

	game.Board.Matrix.Iterate(func(rt PlacedHotel, x int, y int, idx int) {
		if rt.Hotel != NoHotel && rt.Hotel != UndefinedHotel {
			chainsMap[rt.Hotel] = struct{}{}
		}
	})

	chainsSlice := util.Keys(chainsMap)

	sortHotels(chainsSlice)

	return chainsSlice
}

// GetAvailableHotelChains
// gets the available hotels, which are not on the board
// this is pretty brutally inefficient but whatever
func GetAvailableHotelChains(game *Game) []Hotel {
	chains := make(map[Hotel]struct{}, 0)

	for _, h := range HotelChainList {
		chains[h] = struct{}{}
	}

	game.Board.Matrix.Iterate(func(rt PlacedHotel, x int, y int, idx int) {
		delete(chains, rt.Hotel)
	})

	keys := util.Keys(chains)

	sortHotels(keys)

	return keys
}

func sortChainSizeMap(hotelSizeMap map[Hotel]int) []HotelChain {
	chainSizeSortedList := make([]HotelChain, 0)
	for hotel, size := range hotelSizeMap {
		chainSizeSortedList = append(chainSizeSortedList, HotelChain{
			Hotel: hotel,
			Size:  size,
		})
	}

	sort.Slice(chainSizeSortedList, func(i, j int) bool {
		return chainSizeSortedList[i].Size > chainSizeSortedList[j].Size
	})

	return chainSizeSortedList
}

func (p *Player) returnedAmount(tradeInAmount int) int {
	return tradeInAmount / 2
}

func (p *Player) canTradeIn(in Stock, out Stock, tradeInAmount int) error {

	if (tradeInAmount % 2) != 0 {
		return errors.New("trade in amount must be a multiple of two")
	}

	if p.remainingStock(in) < tradeInAmount {
		return errors.New("player does not have enough stock to trade in for")
	}

	returnAmount := p.returnedAmount(tradeInAmount)
	if returnAmount < 1 {
		return nil
	}

	outRemaining := p.Game.remainingStock(out)

	if outRemaining < returnAmount {
		return errors.New("there is not enough stock to trade-in for")
	}

	return nil
}

func (p *Player) tradeIn(in Stock, out Stock, tradeInAmount int) error {

	err := p.canTradeIn(in, out, tradeInAmount)
	if err != nil {
		return err
	}

	returnAmount := p.returnedAmount(tradeInAmount)

	// return the stock to the bank
	_ = p.Game.Inventory.takeHotelStock(Hotel(in), tradeInAmount, p.Inventory)

	// take half the amount of the new stock
	_ = p.Inventory.takeHotelStock(Hotel(out), returnAmount, p.Game.Inventory)

	return nil
}

func (p *Player) canSellStock(stock Stock, amount int) error {
	if p.remainingStock(stock) < amount {
		return errors.New("amount sold cannot be greater than the amount the player has")
	}

	return nil
}

func (p *Player) sellStock(stock Stock, amount int) error {
	err := p.canSellStock(stock, amount)
	if err != nil {
		return err
	}

	hotel := Hotel(stock)
	err = p.Game.Inventory.takeHotelStock(hotel, amount, p.Inventory)
	if err != nil {
		panic("shouldn't take more than there are")
	}

	chainSize := countHotelChain(p.Game, hotel)
	value := sharesCalc(hotel, chainSize, amount)

	p.Inventory.Money += value

	return nil
}

func (p *Player) remainingStock(stock Stock) int {
	return len(p.Inventory.Stocks[Hotel(stock)].Items)
}
