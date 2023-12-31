package ai

import (
	"acquire/internal/acquire"
	"acquire/internal/util"
	"errors"
	"fmt"
	"git.sr.ht/~bonbon/gmcts"
	"strconv"
	"strings"
)

type HumanAgent struct {
}

func NewHumanAgent() *HumanAgent {
	return &HumanAgent{}
}

func (agent HumanAgent) SelectAction(game *acquire.Game, actions []gmcts.Action) (gmcts.Action, error) {

	if len(actions) == 0 {
		panic("there are no actions")
	}

	if len(actions) == 1 {
		return actions[0], nil
	}

	// ask for input until something valid is selected
	for {
		var action gmcts.Action
		var err error
		switch game.NextActionType {
		case acquire.ActionType_PlaceTile:
			action, err = handlePlaceTileActions(game, actions)
			break
		case acquire.ActionType_PickHotelToFound:
			action, err = handlePickHotelToFoundActions(game, actions)
			break
		case acquire.ActionType_PickHotelToMerge:
			action, err = handlePickHotelToMergeActions(game, actions)
			break
		case acquire.ActionType_Merge:
			action, err = handleMergeActions(game, actions)
			break
		case acquire.ActionType_PurchaseStock:
			action, err = handlePurchaseStockActions(game, actions)
			break
		default:
			panic(fmt.Sprintf("action %s is not handled", game.NextActionType))
		}

		if err != nil {
			fmt.Println("input err, " + err.Error())
			continue // the ol' try again strat
		}

		return action, err
	}

}

func handlePlaceTileActions(game *acquire.Game, actions []gmcts.Action) (gmcts.Action, error) {
	fmt.Println("Select a Tile to Place")

	for i, a := range actions {
		action := util.AsType[acquire.Action_PlaceTile](a)
		fmt.Printf("%d: %s\n", i, action.Tile.String())
	}

	return getTileSelection(actions)
}

func handlePurchaseStockActions(game *acquire.Game, actions []gmcts.Action) (gmcts.Action, error) {

	chains := ""
	for _, h := range game.Computed.ActiveChains {
		chains += fmt.Sprintf("%s: $%d - [%s]\n", h.String(), h.Value(game, 1), h.Initial())
	}

	fmt.Println("Available Stock:\n" + chains)
	fmt.Println("Select a Set of Stocks To Buy (Default=None):")

	input, err := getInput()
	if err != nil {
		return nil, err
	}

	// map string of initials to slice of hotels
	inputHotelMap := make(map[acquire.Hotel]int)
	for _, c := range input {
		hotel, err := acquire.ChainFromInitial(string(c))
		if err != nil {
			return nil, err
		}
		inputHotelMap[hotel]++
	}

	// local function comparing the input hotel map to some input map
	isInputMapEqual := func(m map[acquire.Hotel]int) bool {

		if len(m) != len(inputHotelMap) {
			return false
		}

		for k, v := range m {
			iv, ok := inputHotelMap[k]
			if !ok {
				return false
			}

			if iv != v {
				return false
			}
		}

		return true
	}

	// validate the hotels against the actions in a way where order doesn't matter
	for _, a := range actions {
		action := util.AsType[acquire.Action_PurchaseStock](a)
		hotelMap := action.AsMap()
		if isInputMapEqual(hotelMap) {
			return action, nil
		}
	}

	return nil, errors.New("not a valid input")

}

func handlePickHotelToFoundActions(game *acquire.Game, actions []gmcts.Action) (gmcts.Action, error) {
	fmt.Println("Pick a Hotel to Found (Default=0):")

	for i, a := range actions {
		action := util.AsType[acquire.Action_PickHotelToFound](a)
		fmt.Printf("%d: %s\n", i, action.Hotel.String())
	}

	return getSelection(actions)
}

func handlePickHotelToMergeActions(game *acquire.Game, actions []gmcts.Action) (gmcts.Action, error) {
	fmt.Println("Pick a Hotel to Merge (Default=0):")

	for i, a := range actions {
		action := util.AsType[acquire.Action_PickHotelToMerge](a)
		fmt.Printf("%d: %s\n", i, action.Hotel.String())
	}

	return getSelection(actions)
}

func handleMergeActions(game *acquire.Game, actions []gmcts.Action) (gmcts.Action, error) {
	fmt.Println("Merge (Default=0):")

	for i, a := range actions {
		action := util.AsType[acquire.Action_Merge](a)
		str := ""
		for _, sub := range action.Actions {
			str += fmt.Sprintf("%s %d | ", sub.MergeType.String(), sub.Amount)
		}

		fmt.Printf("%d: %s\n", i, str)
	}

	return getSelection(actions)
}

func getTileSelection(actions []gmcts.Action) (gmcts.Action, error) {
	fmt.Printf("Select one: ")

	input, err := getInput()
	if err != nil {

	}

	tile, err := parseTileStr(input)
	if err != nil {
		inputInt, err := strconv.Atoi(input)

		if err != nil {
			return nil, errors.New("invalid tile selection")
		}

		if inputInt >= len(actions) {
			return nil, errors.New("invalid tile selection")
		}

		return actions[inputInt], nil
	}

	for _, action := range actions {
		_action := util.AsType[acquire.Action_PlaceTile](action)
		if _action.Tile == tile {
			return action, nil
		}
	}

	return nil, errors.New("invalid tile selection")
}

func getSelection(actions []gmcts.Action) (gmcts.Action, error) {
	fmt.Printf("Select one: ")

	input, err := getInput()
	if err != nil {
		return 0, err
	}

	if input == "" {
		return actions[0], nil
	}

	inputInt, err := strconv.Atoi(input)
	if err != nil {
		return nil, errors.New("select one of the numbered actions")
	}

	if inputInt >= len(actions) {
		return nil, errors.New("not a valid action to take")
	}

	return actions[inputInt], nil
}

func getInput() (string, error) {
	var input string
	_, err := fmt.Scanln(&input)

	if err != nil {
		if err.Error() == "unexpected newline" {
			return "", nil
		}

		return "", err
	}

	input = strings.Trim(input, " ")
	input = strings.ToUpper(input)

	return input, nil
}

var chars = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'}

func isTileChar(b byte) bool {
	for _, c := range chars {
		if b == c {
			return true
		}
	}
	return false
}

func parseTileStr(str string) (acquire.Tile, error) {
	str = strings.Trim(str, " ")
	str = strings.ToUpper(str)

	var col int
	var row string
	var err error
	if isTileChar(str[len(str)-1]) {
		// if the end of the string is an A through I char, parse the format 12A
		row = string(str[len(str)-1])
		col, err = strconv.Atoi(str[:len(str)-1])
		if err != nil {
			return acquire.NoTile, err
		}
	} else if isTileChar(str[0]) {
		// if the start of the string is an A through I char, parse the format A12
		row = string(str[0])
		col, err = strconv.Atoi(str[1:])
		if err != nil {
			return acquire.NoTile, err
		}
	} else {
		return acquire.NoTile, errors.New("not a parsable tile string")
	}

	tile := tileFromParts(col, row)

	return tile, nil
}

func tileFromParts(col int, row string) acquire.Tile {
	tileStr := strconv.Itoa(col) + row
	for tileInt, ts := range acquire.TileStringMap {
		if ts == tileStr {
			return acquire.Tile(tileInt)
		}
	}

	return acquire.NoTile
}
