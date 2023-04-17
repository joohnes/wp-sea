package app

import (
	"fmt"
	"github.com/fatih/color"
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
}

type App struct {
	client client
}

func New(c client) *App {
	return &App{
		c,
	}
}

func (a *App) Run() error {

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

	for {
		result, err := a.Shoot(board)
		if err != nil {
			return err
		}
		show(board, status)
		if result == "miss" {
			break
		}
	}
	return nil
}

func show(board *gui.Board, status *StatusResponse) {
	red := color.New(color.FgBlack, color.BgRed).SprintFunc()
	green := color.New(color.FgBlack, color.BgGreen).SprintFunc()

	board.Display()
	fmt.Println("Your name: ", green(NICK))
	fmt.Println("Your description: ", green(DESC))
	fmt.Println("\nYour opponent's name: ", red(status.Opponent))
	fmt.Println("Your opponent's description: ", red(status.Opp_desc))
}
