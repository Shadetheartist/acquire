package acquire_2

import (
	"errors"
	"git.sr.ht/~bonbon/gmcts"
)

type MergerAction int

const (
	Hold MergerAction = iota
	Trade
	Sell
)

type MergeSubAction struct {
	MergeType MergerAction
	Amount    int
}

type Action_Merge struct {
	Actions [3]MergeSubAction
}

func (a Action_Merge) Type() ActionType {
	return ActionType_Merge
}

func (game *Game) getMergeHotelActions() []gmcts.Action {
	mergeActions := make([]gmcts.Action, 0)

	mergeActions = append(mergeActions, Action_Merge{
		Actions: [3]MergeSubAction{
			{
				MergeType: Hold, // keeping it simple for now
				Amount:    0,
			},
		},
	})

	return mergeActions
}

func getNextChainToMerge(mergerState *MergerState) (Hotel, error) {
	for idx, playersRemaining := range mergerState.ChainsToMerge {
		if playersRemaining > 0 {
			return ChainFromIdx(idx), nil
		}
	}
	return NoHotel, errors.New("no valid hotel found")
}

func (game *Game) applyMergeHotel(action Action_Merge) {

	hotelToMerge, err := getNextChainToMerge(&game.MergerState)
	if err != nil {
		panic(err)
	}

	player := game.Players[game.MergerState.MergingPlayerIdx]

	game.payShareholderBonuses(hotelToMerge)

	goNext := func() {
		game.MergerState.MergingPlayerIdx += 1
		game.MergerState.MergingPlayerIdx = game.MergerState.MergingPlayerIdx % len(game.Players)

		game.MergerState.ChainsToMerge[hotelToMerge.Index()] -= 1

		// err wil be set if there are no more chains to merge
		_, err := getNextChainToMerge(&game.MergerState)

		// if there's no more chains to process, we're done merging
		if err != nil {
			// finally place the piece and propagate the chain
			newPlacedHotel := game.placeTileOnBoard(game.LastPlacedTile, game.MergerState.AcquiringHotel)
			propagateHotelChain(game, newPlacedHotel)

			game.NextActionType = ActionType_PurchaseStock
		}
	}

	// i've set this up so that all the information for a player's merger is provided at once
	// since they can choose multiple things to do, this loop represents those things.
	// for instance, you could have 6 stocks, trade in 4, and sell two - much depends on the action options provided
	for _, subAction := range action.Actions {
		switch subAction.MergeType {

		case Hold:
			goNext()
			return

		case Trade:
			err := player.tradeIn(game, hotelToMerge, game.MergerState.AcquiringHotel, subAction.Amount)
			if err != nil {
				panic(err)
			}
			break

		case Sell:
			err := player.sellStock(game, Stock(hotelToMerge), subAction.Amount)
			if err != nil {
				panic(err)
			}
			break
		}
	}

	goNext()
}
