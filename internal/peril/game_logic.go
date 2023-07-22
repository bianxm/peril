package peril

import "sync"

type GameStatus int

const (
	Pregame GameStatus = iota
	Choosing
	Answering
)

type Game struct {
	Board        map[int]map[int]string
	Players      []string
	PlayerStates map[string]PlayerState
	GameStatus   GameStatus
	Mutex        sync.Mutex
	// Scores    map[string]int
}

func NewGame() *Game {
	return &Game{
		Board:        make(map[int]map[int]string),
		GameStatus:   Pregame,
		PlayerStates: make(map[string]PlayerState),
	}
}

type PlayerRole int

const (
	Waiter PlayerRole = iota
	Answerer
	Chooser
	Asker
)

type PlayerState struct {
	Score int
	Role  PlayerRole
	ID    int
}

// possible messages to player:
// player choosing q: - board status
// player answering
// 		enabled buzzer
// 		disabled buzzer
// player is asker
// 		turn timer on
// 		check, cross
