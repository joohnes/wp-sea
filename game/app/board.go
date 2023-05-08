package app

import (
	"fmt"

	"github.com/fatih/color"
	board "github.com/grupawp/warships-lightgui"
)

func (a *App) show(board *board.Board, status *StatusResponse) {
	red := color.New(color.FgBlack, color.BgRed).SprintFunc()
	green := color.New(color.FgBlack, color.BgGreen).SprintFunc()

	var score float64
	if a.shotsCount != 0 {
		score = float64(a.shotsHit) / float64(a.shotsCount)
	}

	board.Display()
	fmt.Println("Your name: ", green(a.nick))
	fmt.Println("Your description: ", green(a.desc))
	fmt.Println("\nYour opponent's name: ", red(a.oppNick))
	fmt.Println("Your opponent's description: ", red(a.oppDesc))
	fmt.Println("Hits: ", a.shotsHit, "Shoots: ", a.shotsCount, "Efficiency: ", score)
}

func (a *App) Create() *board.Board {
	return board.New(
		board.NewConfig(),
	)
}
