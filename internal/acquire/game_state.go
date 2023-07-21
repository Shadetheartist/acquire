package acquire

import (
	"acquire/internal/util"
)

// MergerState
// as multiple actions occur during a merger, and the state of the board matters
// we have to the state of the merger over different turns and actions
// nil when not in use
type MergerState struct {
	RemainingChainsToMerge  []HotelChain
	RemainingPlayersToMerge map[HotelChain][]*Player
	Pos                     util.Point[int]
	NeighboringHotels       []PlacedHotel
	ChainsInNeighbors       []Hotel
	LargestChains           []Hotel
	AcquiringHotel          Hotel
	AcquiredHotel           Hotel
}

func (ms *MergerState) clone(playerMap map[*Player]*Player) *MergerState {
	if ms == nil {
		return nil
	}

	cloneRPTM := func() map[HotelChain][]*Player {
		clone := make(map[HotelChain][]*Player)

		for k, players := range ms.RemainingPlayersToMerge {
			clone[k] = make([]*Player, len(players))
			for i, p := range players {
				clone[k][i] = playerMap[p]
			}
		}

		return clone
	}

	clone := &MergerState{
		// will copy naturally
		Pos:            ms.Pos,
		AcquiringHotel: ms.AcquiringHotel,
		AcquiredHotel:  ms.AcquiredHotel,

		// needs to be cloned
		RemainingPlayersToMerge: cloneRPTM(),
		NeighboringHotels:       util.Clone(ms.NeighboringHotels),
		ChainsInNeighbors:       util.Clone(ms.ChainsInNeighbors),
		LargestChains:           util.Clone(ms.LargestChains),
		RemainingChainsToMerge:  util.Clone(ms.RemainingChainsToMerge),
	}

	return clone
}

// FoundState
// sub-state used when founding a new hotel chain
// nil when not in use
type FoundState struct {
	FoundingHotel Hotel
	Pos           util.Point[int]
}
