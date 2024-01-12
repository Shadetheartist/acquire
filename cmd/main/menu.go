package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PlayerType int

const (
	Human PlayerType = iota
	AI
)

type GameConfig struct {
	NumPlayers        int
	PlayerTypes       []PlayerType
	AIPlayerStrengths []int
}

var DefaultGameConfig = GameConfig{
	NumPlayers:        4,
	PlayerTypes:       []PlayerType{Human, AI, AI, AI},
	AIPlayerStrengths: []int{0, 250, 500, 750},
}

func menu() *GameConfig {
	fmt.Println("Acquire - Shell Version")
	fmt.Println()

	var config GameConfig

	fmt.Print("Num Players? [2-6]: ")
	config.NumPlayers = getBoundedInput("Num Players? [2-6]: ", 2, 6)
	config.PlayerTypes = make([]PlayerType, config.NumPlayers)
	config.AIPlayerStrengths = make([]int, config.NumPlayers)

	for i := 0; i < config.NumPlayers; i++ {
		prompt := fmt.Sprintf("Player %d Type? [1 = Human, 2 = AI]: ", i+1)
		config.PlayerTypes[i] = PlayerType(getBoundedInput(prompt, 1, 2) - 1)
		if config.PlayerTypes[i] == AI {
			prompt := "AI Player Strength? [1 = Easy, 2 = Med, 3 = Hard]: "
			config.AIPlayerStrengths[i] = getBoundedInput(prompt, 1, 3) * 250
		}
	}

	return &config
}

func getBoundedInput(prompt string, min int, max int) int {
	tries := 3
	for t := 0; t < tries; t++ {
		fmt.Print(prompt)
		inputInt, err := getInputInt()

		if err != nil {
			fmt.Println("\nInvalid Input, Try Again.\n")
			continue
		}

		if inputInt < min || inputInt > max {
			fmt.Println("\nInvalid Input, Try Again.\n")
			continue
		}

		return inputInt
	}

	os.Exit(1)

	return 0
}

func getInputInt() (int, error) {
	input, err := getInput()
	if err != nil {
		return 0, err
	}

	inputInt, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return inputInt, nil
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
