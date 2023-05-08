package app

import (
	"time"
)

type client interface {
	InitGame(coords []string, desc, nick, targetOpponent string, wpbot bool) error
	PrintToken()
	Board() ([]string, error)
	Status() (*StatusResponse, error)
	Shoot(coord string) (string, error)
	Resign() error
	GetOppDesc() (string, string, error)
	Refresh() error
	PlayerList() ([]map[string]string, error)
	Stats() (map[string][]int, error)
	StatsPlayer(nick string) ([]int, error)
}

type App struct {
	client     client
	nick       string
	desc       string
	oppShots   []string
	oppNick    string
	oppDesc    string
	shotsCount int
	shotsHit   int
}

func New(c client) *App {
	return &App{
		c,
		"",
		"",
		[]string{},
		"",
		"",
		0,
		0,
	}
}

func (a *App) Run() error {
Start:
	err := a.getName()
	if err != nil {
		return err
	}
	err = a.getDesc()
	if err != nil {
		return err
	}
	err = a.ChooseOption()
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

	a.oppNick, a.oppDesc, err = a.client.GetOppDesc()
	if err != nil {
		return err
	}

	board := a.Create()
	board.Import(boardCoords)
	a.show(board, status)

	// MAIN GAME LOOP
	for {
		if a.CheckIfWon() {
			time.Sleep(30 * time.Second)
			goto Start
		}
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
