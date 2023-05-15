package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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

func (a *App) Shoot(coord string) error {
	err := a.WaitForTurn()
	if err != nil {
		return err
	}
Again:
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
		goto Again
	case "sunk":
		// bd.Set(board.Right, coord, board.Hit)
		// bd.CreateBorder(board.Right, coord)
		a.enemy_states[coord[0]-97][coord[1]-49] = "Hit"
		a.shotsCount += 1
		a.shotsHit += 1
		goto Again
	}

	return nil
}

////////////////////////////////////////////
// TUTAJ SIE KOŃCZĄ DOBRE FUNKCJE
////////////////////////////////////////////

func (a *App) Play(status *StatusResponse, coordchan chan string) error {
	var coord string
	select {
	case coordchan <- coord:
		err := a.Shoot(coord)
		if err != nil {
			return err
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

func (a *App) OpponentShots() error {
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
	a.oppShots = status.Opp_shots

	for _, v := range a.oppShots {
		err := a.HitOrMiss(strings.ToLower(v))
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) HitOrMiss(coord string) error {
	x := coord[0] - 97
	if x > 9 || x < 0 {
		return errors.New("Wrong coord!")
	}
	y := coord[1] - 49
	if y > 9 || y < 0 {
		return errors.New("Wrong coord!")
	}
	state := a.my_states[coord[0]-97][coord[1]-49]
	switch state {
	case "Ship":
		a.enemy_states[coord[0]-97][coord[1]-49] = "Hit"
	default:
		a.enemy_states[coord[0]-97][coord[1]-49] = "Miss"
	}
	return nil
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
		if len(playerlist) != 0 {

			fmt.Println("Waiting players: ")
			for i, x := range playerlist {
				fmt.Println(i, x["nick"])
			}
			fmt.Println("Do you want to wait for another player? y/n")
			answer, err = a.getAnswer()
			if err != nil {
				return err
			}
			if answer == "y" {
				err = a.client.InitGame(nil, a.desc, a.nick, "", false)
				if err != nil {
					return err
				}
				return nil
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
		} else {
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
			default:
				fmt.Println("Please enter a number from the list!")
			}
		}
	}
	return nil
}
