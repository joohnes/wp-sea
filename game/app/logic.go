package app

import (
	"fmt"
	"github.com/fatih/color"
	board "github.com/grupawp/warships-lightgui"
	"time"
)

const (
	waitDuration = 3
)

func (a *App) WaitForStart() (status *StatusResponse, err error) {
	for {
		status, err := a.client.Status()
		if err != nil {
			return nil, err
		}
		if status.Game_status == "game_in_progress" {
			return status, nil
		}
		time.Sleep(waitDuration * time.Second)
	}
}

func (a *App) WaitForTurn() error {
	for {
		status, err := a.client.Status()
		if err != nil {
			return err
		}
		if status.Should_fire == true {
			break
		}
		time.Sleep(waitDuration * time.Second)
	}
	return nil
}

func (a *App) Shoot(bd *board.Board) (string, error) {
	err := a.WaitForTurn()
	if err != nil {
		return "", err
	}
	coord, err := getCoord()
	if err != nil {
		return "", err
	}

	result, err := a.client.Shoot(coord)
	if err != nil {
		return "", err
	}
	switch result {
	case "miss":
		bd.Set(board.Right, coord, board.Miss)
	case "hit":
		bd.Set(board.Right, coord, board.Hit)
	case "sunk":
		bd.Set(board.Right, coord, board.Hit)
		bd.CreateBorder(board.Right, coord)
	}

	return result, nil
}

func (a *App) Play(board *board.Board, status *StatusResponse) error {
	for {
		result, err := a.Shoot(board)
		if err != nil {
			fmt.Println(err)
			continue
		}
		show(board, status)
		if result == "miss" {
			break
		}
	}
	fmt.Print("Waiting for bot")
	for i := 0; i < 2; i++ {
		fmt.Print(".")
		time.Sleep(time.Second)
	}
	fmt.Print("\n")
	return nil
}

func show(board *board.Board, status *StatusResponse) {
	red := color.New(color.FgBlack, color.BgRed).SprintFunc()
	green := color.New(color.FgBlack, color.BgGreen).SprintFunc()

	board.Display()
	fmt.Println("Your name: ", green(NICK))
	fmt.Println("Your description: ", green(DESC))
	fmt.Println("\nYour opponent's name: ", red(status.Opponent))
	fmt.Println("Your opponent's description: ", red(status.Opp_desc))
}

func (a *App) OpponentShots(bd *board.Board) error {
	var status *StatusResponse
	var err error
	for {
		status, err = a.client.Status()
		if err != nil {
			return err
		}
		if len(a.oppShots) != len(status.Opp_shots) {
			break
		}
		time.Sleep(waitDuration * time.Second)
	}
	currOppShots := status.Opp_shots
	newOppShots := currOppShots[len(a.oppShots):]
	a.oppShots = currOppShots

	for _, v := range newOppShots {
		bd.HitOrMiss(board.Left, v)
	}
	show(bd, status)
	return nil
}
