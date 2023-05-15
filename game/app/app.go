package app

import (
	"time"

	gui "github.com/grupawp/warships-gui/v2"
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
	client       client
	nick         string
	desc         string
	oppShots     []string
	oppNick      string
	oppDesc      string
	shotsCount   int
	shotsHit     int
	my_states    [10][10]gui.State
	enemy_states [10][10]gui.State
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
		[10][10]gui.State{},
		[10][10]gui.State{},
	}
}

func (a *App) Run() error {
	// Start:
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
	status, err := a.WaitForStart()
	if err != nil {
		return err
	}

	a.oppNick, a.oppDesc, err = a.client.GetOppDesc()
	if err != nil {
		return err
	}

	// MAIN GAME LOOP

	go func() error {
		for {
			if a.CheckIfWon() {
				time.Sleep(30 * time.Second)
				// goto Start
			}
			err = a.Play(status)
			if err != nil {
				return err
			}
			err = a.OpponentShots()
			if err != nil {
				return err
			}
		}
	}()
	a.ShowBoard()
	return nil
}
