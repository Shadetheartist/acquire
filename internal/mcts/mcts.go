package mcts

type Player int
type Action any

type Game interface {
	Value() float64

	ApplyAction(action Action)

	GetActions() []Action

	//Player returns the player that can take the next action
	Player() Player

	//IsTerminal returns true if this game state is a terminal state
	IsTerminal() bool

	//Winners returns a list of players that have won the game if
	//IsTerminal() returns true
	Winners() []Player
}

type MCTS struct {
}
