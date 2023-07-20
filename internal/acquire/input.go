package acquire

import (
	"acquire/internal/util"
	"errors"
	"fmt"
	"strconv"
)

// this input api is too coupled to the console style input
// should be refactored for better general purpose support

type InputRequest struct {
	InputType   string
	Instruction string
}

type InputResponse struct {
	Value string
}

type IInput interface {
	GetInput(InputRequest) (InputResponse, error)
}

func getBuyStockInput(game *Game) Hotel {
	availableStocks := hotelsWithPurchasableStock(game)

	if len(availableStocks) == 0 {
		return NoHotel
	}

	availableChainsInitials := hotelsAsInitials(availableStocks)
	request := InputRequest{
		InputType: "stock",
		Instruction: fmt.Sprintf(
			"Select a Stock to Buy %s: ",
			availableChainsInitials,
		),
	}

	return getValidInput[Hotel](game, request, func(response InputResponse) (Hotel, error) {
		if response.Value == "" {
			return NoHotel, nil
		}

		hotel := hotelFromInitial(response.Value)
		if !isActualHotelChain(hotel) {
			return NoHotel, fmt.Errorf("'%s' is not a valid hotel chain", hotel)
		}

		return hotel, nil
	})
}

func getNumStocksToBuy(game *Game, hotel Hotel, max int) int {

	value := sharesCalc(hotel, countHotelChain(game, hotel), 1)
	request := InputRequest{
		InputType: "number",
		Instruction: fmt.Sprintf(
			"How many? ($%d/share) [0-%d]:",
			value,
			max,
		),
	}

	return getValidInput[int](game, request, func(response InputResponse) (int, error) {
		if response.Value == "" {
			return 0, nil
		}

		number, err := strconv.Atoi(response.Value)
		if err != nil {
			return 0, err
		}

		if number < 0 {
			return 0, fmt.Errorf("must input a number from 0 to %d", max)
		}

		if number > max {
			return 0, fmt.Errorf("must input a number from 0 to %d", max)
		}

		return number, nil
	})
}

func getTradeInAmount(game *Game, player *Player, in Stock, out Stock) int {

	inRemaining := player.remainingStock(in)
	outRemaining := game.remainingStock(out)
	max := util.Min(inRemaining, outRemaining/2)

	request := InputRequest{
		InputType: "number",
		Instruction: fmt.Sprintf(
			"Trade in (0-%d): ",
			max,
		),
	}

	return getValidInput[int](game, request, func(response InputResponse) (int, error) {
		if response.Value == "" {
			return 0, nil
		}

		number, err := strconv.Atoi(response.Value)
		if err != nil {
			return 0, err
		}

		if number < 0 {
			return 0, fmt.Errorf("must input a number from 0 to %d", max)
		}

		if number > max {
			return 0, fmt.Errorf("must input a number from 0 to %d", max)
		}

		if (number % 2) != 0 {
			return 0, fmt.Errorf("must input a multiple of two")
		}

		return number, nil
	})
}

func getSellAmount(game *Game, player *Player, stock Stock) int {

	remaining := player.remainingStock(stock)
	value := sharesCalc(Hotel(stock), countHotelChain(game, Hotel(stock)), 1)
	max := remaining

	request := InputRequest{
		InputType: "number",
		Instruction: fmt.Sprintf(
			"Sell Stock in %s, $%d/share [0-%d]: ",
			stock.String(),
			value,
			max,
		),
	}

	return getValidInput[int](game, request, func(response InputResponse) (int, error) {
		if response.Value == "" {
			return 0, nil
		}

		number, err := strconv.Atoi(response.Value)
		if err != nil {
			return 0, err
		}

		if number < 0 {
			return 0, fmt.Errorf("must input a number from 0 to %d", max)
		}

		if number > max {
			return 0, fmt.Errorf("must input a number from 0 to %d", max)
		}

		return number, nil
	})
}

func getHotelInput(game *Game, reason string, availableChains []Hotel) Hotel {

	availableChainsInitials := hotelsAsInitials(availableChains)
	request := InputRequest{
		InputType: "hotel",
		Instruction: fmt.Sprintf(
			"Select a Hotel Chain (%s) %s: ",
			reason,
			availableChainsInitials,
		),
	}

	return getValidInput[Hotel](game, request, func(response InputResponse) (Hotel, error) {
		hotel := hotelFromInitial(response.Value)
		if !isActualHotelChain(hotel) {
			return NoHotel, fmt.Errorf("'%s' is not a valid hotel chain", hotel)
		}

		if _, ok := util.IndexOf(availableChains, hotel); !ok {
			return NoHotel, fmt.Errorf("'%s' is not available", hotel)
		}

		return hotel, nil
	})
}

func getTileInput(game *Game, player *Player) Tile {

	legalMoves := player.legalMoves()
	legalMoveStrings := util.Map(legalMoves, func(val Tile) string {
		return val.String()
	})

	request := InputRequest{
		InputType:   "tile",
		Instruction: fmt.Sprintf("Place a Tile %s: ", legalMoveStrings),
	}

	return getValidInput[Tile](game, request, func(response InputResponse) (Tile, error) {
		tile := Tile(response.Value)
		if !isActualTile(tile) {
			return NoTile, fmt.Errorf("'%s' is not a valid tile", tile)
		}

		// if player does not have the tile in their inventory
		if _, ok := util.IndexOf(legalMoves, tile); !ok {
			return NoTile, fmt.Errorf("'%s' is not a legal move", tile)
		}

		return tile, nil
	})
}

func getMergerActionTypeInput(game *Game, hotel Hotel) MergerAction {

	request := InputRequest{
		InputType:   "mergerAction",
		Instruction: fmt.Sprintf("%s has been acquired, choose to Hold, Trade, or Sell your stock [H, T, S]: ", hotel),
	}

	return getValidInput[MergerAction](game, request, func(response InputResponse) (MergerAction, error) {
		switch response.Value {
		case "H":
			return Hold, nil
		case "T":
			return Trade, nil
		case "S":
			return Sell, nil
		default:
			return Hold, errors.New("input not one of [H, T, S]")
		}
	})
}

func getInput(game *Game, request InputRequest) InputResponse {
	input, err := game.inputInterface.GetInput(request)
	if err != nil {
		panic(err)
	}
	return input
}

func getValidInput[T any](game *Game, request InputRequest, parser func(response InputResponse) (T, error)) T {
	var z T
	tile, err := attempt[T](10, func(attemptsLeft int, last error) (T, error) {

		if last != nil {
			sendMsg(game, fmt.Sprintf("%s (%d attempts left)", last.Error(), attemptsLeft))
		}

		response := getInput(game, request)
		val, err := parser(response)

		if err != nil {
			return z, err
		}

		return val, nil
	})

	if err != nil {
		game.abort("player wouldn't provide a valid input. exiting game.")
	}

	return tile
}

// sendMsg
// input requests of type "msg" are not to be responded to
// they simply instruct the interface to show a message
func sendMsg(game *Game, msg string) {
	getInput(game, InputRequest{
		InputType:   "msg",
		Instruction: msg,
	})
}

func attempt[T any](maxAttempts int, fn func(attemptsLeft int, last error) (T, error)) (T, error) {
	attempt := 0
	var err error
	var val T
	for attempt < maxAttempts {
		attempt++
		val, err = fn(maxAttempts-attempt+1, err)
		if err == nil {
			return val, nil
		}
	}

	return val, fmt.Errorf("failed to complete func after %d retries: %s", maxAttempts, err.Error())
}
