package app

import (
	"fmt"
	"strconv"
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
		fmt.Println(status.Game_status)
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

func (a *App) Shoot() error {
	err := a.WaitForTurn()
	if err != nil {
		return err
	}
	coord, err := a.getCoord()
	if err != nil {
		return err
	}

	result, err := a.client.Shoot(coord)
	if err != nil {
		return err
	}

	switch result {
	case "miss":
		// bd.Set(board.Right, coord, board.Miss)
		a.enemy_states[coord[0]-97][coord[1]-49] = "Miss"
		a.shotsCount += 1
	case "hit":
		// bd.Set(board.Right, coord, board.Hit)
		a.enemy_states[coord[0]-97][coord[1]-49] = "Hit"
		a.shotsCount += 1
		a.shotsHit += 1
	case "sunk":
		// bd.Set(board.Right, coord, board.Hit)
		// bd.CreateBorder(board.Right, coord)
		a.enemy_states[coord[0]-97][coord[1]-49] = "Hit"
		a.shotsCount += 1
		a.shotsHit += 1
	}

	return nil
}

func (a *App) Play(status *StatusResponse) error {
	for {
		err := a.Shoot()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// a.show(board, status)

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
	switch answer {
	case "1":
		err := a.client.InitGame(nil, a.desc, a.nick, "", true)
		if err != nil {
			return err
		}
	case "2":
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
		fmt.Println("Waiting players: ")
		for i, x := range playerlist {
			fmt.Println(i, x["nick"])
		}
		fmt.Println("Choose a player number: ")
		answer, err = a.getAnswer()
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(answer)
		if err != nil {
			return err
		}
		a.client.Refresh()
		time.Sleep(time.Second * 1)
		fmt.Printf("'%s'", playerlist[i]["nick"])
		err = a.client.InitGame(nil, a.desc, a.nick, playerlist[i]["nick"], false)
		if err != nil {
			return err
		}

	default:
		fmt.Println("Please enter a number from the list!")
	}
	return nil
}
