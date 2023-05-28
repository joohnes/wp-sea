package app

import (
	"context"
	"fmt"
	"github.com/inancgumus/screen"
	"github.com/joohnes/wp-sea/game/logger"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
)

const showErrors = false

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
	enemyShips   map[int]int
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
		map[int]int{4: 1, 3: 2, 2: 3, 1: 4},
	}
}

func (a *App) Run() error {
	log := logger.GetLoggerInstance()
	screen.Clear()
	screen.MoveTopLeft()
	for {
		err := a.getName()
		if err == nil {
			break
		}
		log.Println(err)
		if showErrors {
			fmt.Println(err)
		}
	}

	for {
		err := a.getDesc()
		if err == nil {
			break
		}
		log.Println(err)
		if showErrors {
			fmt.Println(err)
		}
	}
	for {
		for {
			err := a.ChooseOption()
			if err == nil {
				break
			}
			log.Println(err)
			if showErrors {
				fmt.Println(err)
			}

			if a.gameState != StateStart {
				break
			}
			if err.Error() == "player not found" {
				fmt.Println("Player not found. Perhaps you did not play any games?")
				time.Sleep(2 * time.Second)
				continue
			}
			time.Sleep(2 * time.Second)
			fmt.Println("Server error occurred. Please try again")
		}

		err := ServerErrorWrapper(a.WaitForStart)
		if err != nil {
			log.Println(err)
			if showErrors {
				fmt.Println(err)
			}
		}

		for {
			var err error
			err = ServerErrorWrapper(func() error {
				a.oppDesc, a.oppNick, err = a.client.GetOppDesc()
				if err != nil {
					return err
				}
				return nil
			})
			if err == nil {
				break
			}
			log.Println(err)
			if showErrors {
				fmt.Println(err)
			}
		}

		for {
			err := a.GetBoard()
			if err == nil {
				break
			}
			log.Println(err)
			if showErrors {
				fmt.Println(err)
			}
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
		a.ShowBoard(ctx, coordchan, textchan, errorchan, timeLeftchan)
	}
}
