package peril

import "sync"

type GameStatus int

const (
	Pregame GameStatus = iota
	Choosing
	Asking
	Buzzing
	Answering
)

type Game struct {
	Board        map[int]map[int]string
	Players      []string
	PlayerStates map[int]*PlayerState
	GameStatus   GameStatus
	Mutex        sync.Mutex
	currQ        [2]int
	lastChooser  int
	// add: current question
	// Scores    map[string]int
}

func NewGame() *Game {
	return &Game{
		Board:        make(map[int]map[int]string),
		GameStatus:   Pregame,
		PlayerStates: make(map[int]*PlayerState),
		currQ:        [2]int{-1, -1},
	}
}

type PlayerRole int

const (
	Waiter PlayerRole = iota
	Answerer
	Chooser
	Buzzer
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

// given game, client who sent message, and message they sent

func (g *Game) UpdateGame(playerID int, m []byte) {
	g.Mutex.Lock()
	// game state mutation here
	switch g.GameStatus {
	case Pregame:
		// call a different fxn here to change username + input q's on board
	case Choosing:
		if g.PlayerStates[playerID].Role == Chooser { // and m is as expected
			// decode m for choice
			g.lastChooser = playerID
			g.GameStatus = Asking
			for _, pl := range g.PlayerStates {
				// check if they're supposed to be the asker, then mark them as asker
				pl.Role = Buzzer
			}
		}
	case Asking:
		if g.PlayerStates[playerID].Role == Asker { // and m is as expected
			g.GameStatus = Buzzing
		}
	case Buzzing:
		if g.PlayerStates[playerID].Role == Buzzer { // and m is as expected
			g.GameStatus = Answering
			g.PlayerStates[playerID].Role = Answerer
		}
	case Answering:
		if g.PlayerStates[playerID].Role == Asker {
			if string(m) == "check" {
				g.GameStatus = Choosing
				for _, pl := range g.PlayerStates {
					if pl.Role == Answerer {
						pl.Score += g.currQ[1] * 100
						pl.Role = Chooser
					} else {
						pl.Role = Waiter
					}
				}
			}
			if string(m) == "x" {
				count_buzzers := 0
				g.GameStatus = Buzzing
				for _, pl := range g.PlayerStates {
					if pl.Role == Answerer {
						pl.Score -= g.currQ[1] * 100
						pl.Role = Waiter
					}
					if pl.Role == Buzzer {
						count_buzzers += 1
					}
				}
				if count_buzzers == 0 {
					g.GameStatus = Choosing
					for _, pl := range g.PlayerStates {
						if pl.ID == g.lastChooser {
							pl.Role = Chooser
						} else {
							pl.Role = Waiter
						}
					}
				}
			}
		}
	}
	g.Mutex.Unlock()
}
