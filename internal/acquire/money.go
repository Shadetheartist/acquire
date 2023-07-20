package acquire

import "fmt"

var hotelTierMap = map[Hotel]int{
	WorldwideHotel:   0,
	SacksonHotel:     0,
	FestivalHotel:    1,
	ImperialHotel:    1,
	AmericanHotel:    1,
	ContinentalHotel: 2,
	TowerHotel:       2,
}

// sharesCalc
// Returns the total value of a number of shares of a hotel chain of some size.
// Based on the 'Acquire reference chart'.
func sharesCalc(hotel Hotel, size int, numShares int) int {
	return shareValueCalc(hotel, size) * numShares
}

// shareValueCalc
// Returns the value per share of a hotel chain of some size.
// Based on the 'Acquire reference chart'.
func shareValueCalc(hotel Hotel, size int) int {
	tier := hotelTierMap[hotel]
	return tierCalc(size, tier)
}

func tierCalc(size int, tier int) int {
	if tier < 0 || tier > 2 {
		panic(fmt.Sprintf("tier must be within 0-2, was (%d)", tier))
	}
	return sizeCalc(size) + (tier * 100)
}

func sizeCalc(size int) int {
	if size <= 1 {
		return 0
	}

	if size <= 2 {
		return 200
	}

	if size <= 3 {
		return 300
	}

	if size <= 4 {
		return 400
	}

	if size <= 5 {
		return 500
	}

	if size <= 10 {
		return 600
	}

	if size <= 20 {
		return 700
	}

	if size <= 30 {
		return 800
	}

	if size <= 40 {
		return 900
	}

	return 1000
}

// shareholderBonusCalc
// takes in the largest and second-largest hotel chain sizes
// returns the major and minor shareholder bonuses respectively
func shareholderBonusCalc(hotel Hotel, size int, majorSH *Player, primaryNumShares int, minorSh *Player, secondaryNumShares int) (int, int) {

	if primaryNumShares < secondaryNumShares {
		panic("primaryNumShares must be larger than or equal to secondaryNumShares")
	}

	if primaryNumShares == 0 {
		return 0, 0
	}

	tier := hotelTierMap[hotel]
	major := majorShareholderBonusCalc(size, tier)
	minor := minorShareholderBonusCalc(size, tier)

	if majorSH != minorSh && primaryNumShares == secondaryNumShares {
		split := (major + minor) / 2
		split = roundToNearestHundred(split)
		return split, split
	}

	return major, minor
}

func majorShareholderBonusCalc(size int, tier int) int {
	return tierCalc(size, tier) * 10
}

func minorShareholderBonusCalc(size int, tier int) int {
	return majorShareholderBonusCalc(size, tier) / 2
}

func roundToNearestHundred(num int) int {
	return ((num + 99) / 100) * 100
}
