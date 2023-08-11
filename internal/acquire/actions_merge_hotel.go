package acquire

import (
	"errors"
	"git.sr.ht/~bonbon/gmcts"
)

type MergerAction int

func (ma MergerAction) String() string {
	switch ma {
	case Hold:
		return "Hold"
	case Trade:
		return "Trade"
	case Sell:
		return "Sell"
	default:
		panic("AAAH")
	}
}

const (
	Hold MergerAction = iota
	Trade
	Sell
)

const MAX_MERGE_SUB_ACTIONS = 2

type MergeSubAction struct {
	MergeType MergerAction
	Amount    int
}

type Action_Merge struct {
	Actions [MAX_MERGE_SUB_ACTIONS]MergeSubAction
}

func (a Action_Merge) Type() ActionType {
	return ActionType_Merge
}

func (game *Game) getMergeHotelActions() []gmcts.Action {
	mergeActions := make([]gmcts.Action, 0)

	// can always hold
	mergeActions = append(mergeActions, Action_Merge{
		// two 'hold's works the same as one
		Actions: [MAX_MERGE_SUB_ACTIONS]MergeSubAction{
			{
				MergeType: Hold,
				Amount:    0,
			},
		},
	})

	activePlayer := game.ActivePlayer()
	mergedHotel, err := game.getNextChainToMerge()
	acquiringHotel := game.MergerState.AcquiringHotel
	if err != nil {
		panic(err)
	}

	numStocks := activePlayer.Stocks[mergedHotel.Index()]

	// player can sell any number of their shares, so all options must be given
	for j := 0; j < numStocks; j++ {
		mergeActions = append(mergeActions, Action_Merge{
			Actions: [MAX_MERGE_SUB_ACTIONS]MergeSubAction{
				{
					MergeType: Sell,
					Amount:    j + 1,
				},
			},
		})
	}

	// if the player has enough shares to trade,
	// they can trade in any denomination mod 2
	for i := 0; i < numStocks/2; i++ {

		tradeInAmount := (i + 1) * 2

		if game.Stocks[acquiringHotel.Index()] < tradeInAmount {
			break
		}

		// default one with no selling
		mergeActions = append(mergeActions, Action_Merge{
			Actions: [MAX_MERGE_SUB_ACTIONS]MergeSubAction{
				{
					MergeType: Trade,
					Amount:    tradeInAmount,
				},
			},
		})

		// for each trade made, they can also sell any amount of the remaining shares
		// (this gets out of hand a bit in terms of combinations)
		for j := 0; j < numStocks-tradeInAmount; j++ {
			mergeActions = append(mergeActions, Action_Merge{
				Actions: [MAX_MERGE_SUB_ACTIONS]MergeSubAction{
					{
						MergeType: Trade,
						Amount:    tradeInAmount,
					},
					{
						MergeType: Sell,
						Amount:    j + 1,
					},
				},
			})
		}

	}

	return mergeActions
}

func (game *Game) getNextChainToMerge() (Hotel, error) {
	mergerState := game.MergerState
	for idx, playersRemaining := range mergerState.ChainsToMerge {
		if playersRemaining > 0 {
			return ChainFromIdx(idx), nil
		}
	}
	return NoHotel, errors.New("no valid hotel found")
}

func (game *Game) applyMergeHotel(action Action_Merge) {

	hotelToMerge, err := game.getNextChainToMerge()
	if err != nil {
		panic(err)
	}

	player := game.Players[game.MergerState.MergingPlayerIdx]

	goNext := func() {
		game.MergerState.MergingPlayerIdx += 1
		game.MergerState.MergingPlayerIdx = game.MergerState.MergingPlayerIdx % len(game.Players)
		game.MergerState.ChainsToMerge[hotelToMerge.Index()] -= 1

		// err wil be set if there are no more chains to merge
		_, err := game.getNextChainToMerge()

		// if there's no more chains to process, we're done merging
		if err != nil {
			// finally place the piece and propagate the chain
			newPlacedHotel := game.placeTileOnBoard(game.LastPlacedTile, game.MergerState.AcquiringHotel)
			propagateHotelChain(game, newPlacedHotel)

			game.NextActionType = ActionType_PurchaseStock
		}
	}

	// pay them boys
	game.payShareholderBonuses(hotelToMerge)

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
