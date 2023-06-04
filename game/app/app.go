package app

import (
	"context"
	"fmt"
	"github.com/joohnes/wp-sea/game/logger"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
	"github.com/inancgumus/screen"
	"github.com/joohnes/wp-sea/game/helpers"
)

const ShowErrors = true

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
	ResetToken()
}

type App struct {
	client         client
	nick           string
	desc           string
	oppShots       []string
	oppNick        string
	oppDesc        string
	shotsCount     int
	shotsHit       int
	myStates       [10][10]gui.State
	playerStates   [10][10]gui.State
	enemyStates    [10][10]gui.State
	gameState      int
	actualStatus   StatusResponse
	enemyShips     map[int]int
	placeShips     map[int]int
	playerShots    map[string]string
	statistics     map[string]int
	algorithm      bool
	mode           Mode
	LastPlayerHit  string
	algorithmTried []string
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
		[10][10]gui.State{},
		StateStart,
		StatusResponse{},
		map[int]int{4: 1, 3: 2, 2: 3, 1: 4},
		map[int]int{4: 1, 3: 2, 2: 3, 1: 4},
		map[string]string{},
		make(map[string]int),
		false,
		TargetState,
		"",
		[]string{},
	}
}

func (a *App) Run() error {
	// SETUP CHANNELS
	coordchan := make(chan string)
	textchan := make(chan string)
	errorchan := make(chan error)
	timeLeftchan := make(chan int)
	resetTimerchan := make(chan int)
	shipchannel := make(chan string)
	//

	var err error
	err = a.LoadStatistics()
	if err != nil {
		logger.GetLoggerInstance().Println("couldn't load statistics")
	}

	screen.Clear()
	screen.MoveTopLeft()
	for {
		a.nick, err = helpers.GetName()
		if err == nil {
			break
		}
		logger.GetLoggerInstance().Println(err)
	}

	for {
		a.desc, err = helpers.GetDesc()
		if err == nil {
			break
		}
		logger.GetLoggerInstance().Println(err)
	}
	for {
		// SETUP CONTEXTS
		playingCtx, playingCancel := context.WithCancel(context.Background())
		shipSetupCtx, shipCancel := context.WithCancel(context.Background())
		//

		for {
			err = a.ChooseOption(shipSetupCtx, shipchannel, errorchan)
			if err == nil {
				break
			}
			logger.GetLoggerInstance().Println(err)

			if a.gameState != StateStart {
				break
			}
			if err.Error() == "player not found" {
				fmt.Println("Player not found. Perhaps you did not play any games?")
				time.Sleep(2 * time.Second)
				continue
			}
			time.Sleep(2 * time.Second)
			fmt.Println("Server error occurred1. Please try again")
		}
		shipCancel()

		err = helpers.ServerErrorWrapper(ShowErrors, a.WaitForStart)
		if err != nil {
			logger.GetLoggerInstance().Println(err)
		}

		for {
			err = helpers.ServerErrorWrapper(ShowErrors, func() error {
				a.oppNick, a.oppDesc, err = a.client.GetOppDesc()
				if err != nil {
					return err
				}
				return nil
			})
			if err == nil {
				break
			}
			logger.GetLoggerInstance().Println(err)
		}

		for {
			err = helpers.ServerErrorWrapper(ShowErrors, func() error {
				err := a.GetBoard()
				if err != nil {
					return err
				}
				return nil
			})
			if err == nil {
				break
			}
			logger.GetLoggerInstance().Println(err)
		}

		// SETUP GOROUTINES
		go a.CheckStatus(playingCtx, playingCancel, textchan)
		go a.OpponentShots(playingCtx, errorchan)
		if !a.algorithm {
			go a.Play(playingCtx, coordchan, textchan, errorchan, resetTimerchan)
		} else {
			go a.AlgorithmPlay(playingCtx, textchan, errorchan, resetTimerchan)
		}
		go a.Timer(playingCtx, timeLeftchan, resetTimerchan)
		//

		// SHOW BOARD
		a.ShowBoard(playingCtx, coordchan, textchan, errorchan, timeLeftchan)
		if a.gameState != StateEnded {
			for {
				if a.gameState == StateEnded {
					break
				}
				err = helpers.ServerErrorWrapper(true, func() error {
					if a.gameState != StateEnded {
						return a.client.Resign()
					} else {
						return nil
					}
				})
				if err == nil {
					break
				}
				logger.GetLoggerInstance().Println(err)
			}
		}
		//a.AddStatistics()
		a.SaveStatistics()
		a.Reset()
	}
}
