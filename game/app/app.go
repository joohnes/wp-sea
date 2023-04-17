package app

import (
	"fmt"
	gui "github.com/grupawp/warships-lightgui"
)

const (
	NICK = "Nerien"
	DESC = "test"
)

type client interface {
	InitGame(coords []string, desc, nick, targetOpponent string, wpbot bool) error
	PrintToken()
	Board() ([]string, error)
	Status() (*StatusResponse, error)
	Shoot(coord string) (string, error)
	Resign() error
}

type App struct {
	client       client
	oppShots     []string
	playerHits   int
	opponentHits int
}

func New(c client) *App {
	return &App{
		c,
		[]string{},
		0,
		0,
	}
}

func (a *App) Run() error {
	fmt.Println("Starting application...")

	err := a.client.InitGame(nil, DESC, NICK, "", true)
	if err != nil {
		return err
	}

	boardCoords, err := a.client.Board()
	if err != nil {
		return err
	}
	status, err := a.WaitForStart()
	if err != nil {
		return err
	}

	board := gui.New(
		gui.NewConfig(),
	)
	board.Import(boardCoords)
	show(board, status)

	// MAIN GAME LOOP
	for {
		a.CheckIfWon()
		err = a.Play(board, status)
		if err != nil {
			return err
		}
		err = a.OpponentShots(board)
		if err != nil {
			return err
		}
	}
}
