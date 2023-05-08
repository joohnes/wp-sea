package app

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	board "github.com/grupawp/warships-lightgui"
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
		if status.Should_fire {
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
	coord, err := a.getCoord()
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
		a.show(board, status)
		if result == "miss" {
			a.shotsCount += 1
			break
		} else {
			a.shotsCount += 1
			a.shotsHit += 1
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
	a.show(bd, status)
	return nil
}

func (a *App) CheckIfWon() bool {
	status, err := a.client.Status()
	if err != nil {
		fmt.Println("Could not get status")
	}
	switch status.Last_game_status {
	case "win":
		green := color.New(color.FgBlack, color.BgGreen).SprintFunc()
		fmt.Println(green("You have won the game!"))
		return true
	case "lose":
		red := color.New(color.FgBlack, color.BgRed).SprintFunc()
		fmt.Println(red("You have lost the game!"))
		return true
	}
	return false
}

func (a *App) ChooseOption() error {
	fmt.Println("1. Play with WPBot")
	fmt.Println("2. Play with another player")
	fmt.Println("Choose an option (number): ")
	answer, err := a.getAnswer()
	if err != nil {
		return err
	}
	if answer == "1" {
		err := a.client.InitGame(nil, a.desc, a.nick, "", true)
		if err != nil {
			return err
		}
	} else if answer == "2" {
		playerlist, err := a.client.PlayerList()
		if err != nil {
			return err
		}
		if len(playerlist) == 0 {
			fmt.Println("No players waiting at the moment")
			fmt.Println("Do you want to wait for another player? y/n")
			answer, err = a.getAnswer()
			if err != nil {
				return err
			}
			switch answer {
			case "y":
				err = a.client.InitGame(nil, a.desc, a.nick, "", false)
				if err != nil {
					return err
				}
				return nil
			case "n":
				return nil
			}
		}
	} else {
		fmt.Println("Please enter a number from the list!")
	}
	return nil
}
