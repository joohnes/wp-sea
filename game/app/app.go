package app

import (
	"context"
	"fmt"
	"time"

	"github.com/joohnes/wp-sea/game/logger"

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
	client        client
	nick          string
	desc          string
	oppShots      []string
	oppNick       string
	oppDesc       string
	shotsCount    int
	shotsHit      int
	myStates      [10][10]gui.State
	playerStates  [10][10]gui.State
	enemyStates   [10][10]gui.State
	gameState     int
	actualStatus  StatusResponse
	enemyShips    map[int]int
	placeShips    map[int]int
	playerShots   map[string]string
	statistics    map[string]int
	algorithm     Algorithm
	LastPlayerHit string
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
		NewAlgorithm(),
		"",
	}
}

func (a *App) Run() error {
	logger.GetLoggerInstance().Info.Println("Started app")
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
		logger.GetLoggerInstance().Error.Println("couldn't load statistics")
	} else {
		logger.GetLoggerInstance().Info.Println("Loaded statistics")
	}

	screen.Clear()
	screen.MoveTopLeft()
	for {
		a.nick, err = helpers.GetName()
		if err == nil {
			break
		}
		logger.GetLoggerInstance().Error.Println(err)
	}
	logger.GetLoggerInstance().Info.Printf("Setup name: %s\n", a.nick)

	for {
		a.desc, err = helpers.GetDesc()
		if err == nil {
			break
		}
		logger.GetLoggerInstance().Error.Println(err)
	}
	logger.GetLoggerInstance().Info.Printf("Setup desc: %s\n", a.desc)

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
			logger.GetLoggerInstance().Error.Println(err)

			if a.gameState != StateStart {
				break
			}
			if err.Error() == "not found" {
				fmt.Println("Player not found.")
				time.Sleep(2 * time.Second)
				continue
			}

			time.Sleep(2 * time.Second)
			fmt.Println("Server error occurred. Please try again")
		}
		shipCancel()

		err = helpers.ServerErrorWrapper(ShowErrors, a.WaitForStart)
		if err != nil {
			logger.GetLoggerInstance().Error.Println(err)
		} else {
			logger.GetLoggerInstance().Info.Println("Game started")
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
			logger.GetLoggerInstance().Error.Println(err)
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
			logger.GetLoggerInstance().Error.Println(err)
		}

		// SETUP GOROUTINES
		go a.CheckStatus(playingCtx, playingCancel, textchan)
		go a.OpponentShots(playingCtx, errorchan)
		if !a.algorithm.enabled {
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
					if a.actualStatus.GameStatus != "ended" {
						return a.client.Resign()
					} else {
						return nil
					}
				})
				if err == nil {
					break
				}
				logger.GetLoggerInstance().Error.Println(err)
			}
		}
		var won string
		if a.actualStatus.LastGameStatus == "win" {
			won = a.nick
		} else {
			won = a.oppNick
		}
		logger.GetLoggerInstance().Info.Printf("Game ended, %s won", won)
		a.SaveStatistics()
		a.Reset()
	}
}
