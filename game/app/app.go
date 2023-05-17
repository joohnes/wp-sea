package app

import (
	"context"
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
	myStates     [10][10]gui.State
	enemyStates  [10][10]gui.State
	gameState    Gamestate
	actualStatus StatusResponse
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
		StateStart,
		StatusResponse{},
	}
}

func (a *App) Run() error {
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
	_, err = a.WaitForStart()
	if err != nil {
		return err
	}

	a.oppNick, a.oppDesc, err = a.client.GetOppDesc()
	if err != nil {
		return err
	}
	err = a.GetBoard()
	if err != nil {
		return err
	}

	// SETUP CHANNELS
	coordchan := make(chan string)
	textchan := make(chan string)
	errorchan := make(chan error)
	timeLeftchan := make(chan int)
	resetTimerchan := make(chan int)
	//

	// SETUP CONTEXTS
	ctx, cancel := context.WithCancel(context.Background())
	//

	// SETUP GOROUTINES
	go a.CheckStatus(ctx, cancel, textchan)
	go a.OpponentShots(ctx, errorchan)
	go a.Play(ctx, coordchan, textchan, errorchan, resetTimerchan)
	go a.Timer(ctx, timeLeftchan, resetTimerchan)
	//

	// SHOW BOARD
	a.ShowBoard(coordchan, textchan, errorchan, timeLeftchan)
	return nil
}
