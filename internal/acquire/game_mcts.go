package acquire

// this file contains the MCTS related functions

import (
	"acquire/internal/util"
	"fmt"
	"git.sr.ht/~bonbon/gmcts"
	"sort"
)

func (game *Game) buildMergeActions() []gmcts.Action {
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

func (game *Game) GetActions() []gmcts.Action {
	switch game.NextActionType {

	case ActionType_PlaceTile:
		moves := game.CurrentPlayer().legalMoves()
		var skip bool
		if len(moves) < 1 {
			skip = game.CurrentPlayer().refreshOrSkip(1)
		}

		actions := util.Map(moves, func(val Tile) gmcts.Action {
			return Action_PlaceTile{Tile: val}
		})

		if skip {
			actions = append(actions, Action_PlaceTile{Tile: NoTile})
		}

		return actions

	case ActionType_PickHotelToFound:
		return util.Map(GetAvailableHotelChains(game), func(val Hotel) gmcts.Action {
			return Action_PickHotelToFound{Hotel: val}
		})

	case ActionType_PickHotelToMerge:
		return util.Map(game.MergerState.LargestChains, func(val Hotel) gmcts.Action {
			return Action_PickHotelToMerge{Hotel: val}
		})

	case ActionType_Merge:
		return game.buildMergeActions()

	case ActionType_PurchaseStock:
		// keeping this simple for now
		actions := util.Map(getActiveHotelChains(game), func(val Hotel) gmcts.Action {
			return Action_PurchaseStock{
				Purchases: [3]StockPurchase{
					{
						Hotel:  val,
						Amount: 0,
					},
				},
			}
		})

		// default action to buy nothing
		actions = append(actions, Action_PurchaseStock{
			Purchases: [3]StockPurchase{
				{
					Hotel:  NoHotel,
					Amount: 0,
				},
			},
		})

		return actions

	case ActionType_EndGame:
		var actions []gmcts.Action

		// AI wants to see what happens if nobody ever ends the game lol
		// need to say NO
		if len(game.Inventory.Tiles.Items) > 0 {
			actions = append(actions, Action_EndGame{false})
		}

		// if the player can end, provide the action
		if _, b := game.canEnd(); b {
			actions = append(actions, Action_EndGame{end: true})
		}

		return actions

	default:
		panic("action type not implemented here")
	}
}

func (game *Game) ApplyAction(gmctsAction gmcts.Action) (gmcts.Game, error) {

	clone := game.clone()

	action, ok := gmctsAction.(IAction)
	if !ok {
		panic("action type was not convertable to IAction")
	}
	switch action.Type() {
	case ActionType_PlaceTile:
		clone.applyPlaceTileAction(action.(Action_PlaceTile))
		break
	case ActionType_PickHotelToFound:
		clone.applyPickHotelToFoundAction(action.(Action_PickHotelToFound))
		break
	case ActionType_PickHotelToMerge:
		clone.applyPickHotelToMergeAction(action.(Action_PickHotelToMerge))
		break
	case ActionType_Merge:
		clone.applyMerge(action.(Action_Merge))
		break
	case ActionType_PurchaseStock:
		clone.applyPurchaseStockAction(action.(Action_PurchaseStock))
		break
	case ActionType_EndGame:
		clone.applyEndGameAction(action.(Action_EndGame))
		break
	default:
		panic(fmt.Sprintf("action %d is not handled", action))
	}

	return clone, nil
}

func (game *Game) clone() *Game {

	clone := &Game{
		// will copy naturally
		WillEnd:        game.WillEnd,
		NextActionType: game.NextActionType,
		LastPlacedTile: game.LastPlacedTile,
		IsOver:         game.IsOver,
		Turn:           game.Turn,
		FoundState:     game.FoundState,

		// must be cloned
		Players:     nil, //added later, needed ref
		Board:       game.Board.clone(),
		Inventory:   game.Inventory.clone(),
		MergerState: nil,
	}

	clone.Inventory.Owner = clone

	clonedPlayers := func() []*Player {
		clonedPlayers := make([]*Player, len(game.Players))
		for i, p := range game.Players {
			clonedPlayer := p.clone()
			clonedPlayer.Game = clone
			clonedPlayers[i] = clonedPlayer
		}
		return clonedPlayers
	}

	clone.Players = clonedPlayers()

	playerMap := make(map[*Player]*Player)
	for i, p := range game.Players {
		playerMap[p] = clone.Players[i]
	}

	clone.MergerState = game.MergerState.clone(playerMap)

	return clone
}

func (game *Game) Player() gmcts.Player {
	if game.MergerState != nil {
		if len(game.MergerState.RemainingChainsToMerge) > 0 {
			chainToMerge := game.MergerState.RemainingChainsToMerge[0]
			if len(game.MergerState.RemainingPlayersToMerge) > 0 {
				player := game.MergerState.RemainingPlayersToMerge[chainToMerge][0]
				return gmcts.Player(player.Id)
			}
		}
	}

	return gmcts.Player(game.CurrentPlayer().Id)
}

func (game *Game) IsTerminal() bool {

	return game.IsOver
}

func (game *Game) Winners() []gmcts.Player {
	sortedPlayers := make([]*Player, len(game.Players))

	copy(sortedPlayers, game.Players)
	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].Inventory.Money > sortedPlayers[j].Inventory.Money
	})

	winners := util.Filter(sortedPlayers, func(val *Player) bool {
		return val.Inventory.Money == sortedPlayers[0].Inventory.Money
	})

	return util.Map(winners, func(player *Player) gmcts.Player {
		return gmcts.Player(player.Id)
	})
}
