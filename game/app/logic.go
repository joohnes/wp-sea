package app

import (
	"context"
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
			a.actualStatus = *status
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
	//	err := a.WaitForTurn()
	//	if err != nil {
	//		return err
	//	}
	//Again:
	//	if err != nil {
	//		return err
	//	}
	coordmap, err := numericCords(coord)
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
		a.enemyStates[coordmap["x"]][coordmap["y"]] = "Miss"
		a.shotsCount += 1
		a.gameState = StateOppTurn
	case "hit":
		// bd.Set(board.Right, coord, board.Hit)
		a.enemyStates[coordmap["x"]][coordmap["y"]] = "Hit"
		a.shotsCount += 1
		a.shotsHit += 1
	case "sunk":
		// bd.Set(board.Right, coord, board.Hit)
		// bd.CreateBorder(board.Right, coord)
		a.enemyStates[coordmap["x"]][coordmap["y"]] = "Hit"
		a.shotsCount += 1
		a.shotsHit += 1
	}

	return nil
}

////////////////////////////////////////////
// TUTAJ SIE KOŃCZĄ DOBRE FUNKCJE
////////////////////////////////////////////

func (a *App) Play(ctx context.Context, coordchan <-chan string, textchan chan<- string, errorchan chan<- error, resettime chan int) {
	var coord string
	for {
		select {
		case coord = <-coordchan:
			_, err := numericCords(coord)
			if err != nil {
				errorchan <- err
			}
			err = a.Shoot(coord)
			if err != nil {
				errorchan <- err
			}
			resettime <- 1
			textchan <- fmt.Sprintf("Shot at %s", coord)
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) HitOrMiss(coord string) error {
	coordmap, err := numericCords(coord)
	if err != nil {
		return err
	}

	state := a.myStates[coordmap["x"]][coordmap["y"]]
	switch state {
	case "Ship":
		a.myStates[coordmap["x"]][coordmap["y"]] = "Hit"
	case "Hit":
		a.myStates[coordmap["x"]][coordmap["y"]] = "Hit"
	default:
		a.myStates[coordmap["x"]][coordmap["y"]] = "Miss"
	}
	//for i, x := range a.myStates {
	//	for j, y := range x {
	//		fmt.Println(i, j, y)
	//	}
	//}
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
			err = a.client.Refresh()
			if err != nil {
				return err
			}
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

func (a *App) CheckStatus(ctx context.Context, cancel context.CancelFunc, textchan chan<- string) {
	statusTicker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-statusTicker.C:
			status, err := a.client.Status()
			if err != nil {
				continue
			}
			a.actualStatus = *status
			if status.Game_status == "ended" {
				switch status.Last_game_status {
				case "win":
					textchan <- "You have won the game!"
				case "lose":
					textchan <- "You have lost the game!"
				}
				time.Sleep(5 * time.Second)
				cancel()
			}

		case <-ctx.Done():
			return
		}
	}
}

func (a *App) OpponentShots(ctx context.Context, errorchan chan<- error) {
	oppShotTicker := time.NewTicker(time.Millisecond * 500)
	checkTicker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-oppShotTicker.C:
			diff := difference(a.actualStatus.Opp_shots, a.oppShots)
			a.oppShots = a.actualStatus.Opp_shots
			if len(diff) != 0 {
				for _, v := range diff {
					err := a.HitOrMiss(strings.ToLower(v))
					if err != nil {
						errorchan <- err
					}
				}
			}
			a.gameState = StatePlayerTurn
			// additional check
		case <-checkTicker.C:
			for _, v := range a.actualStatus.Opp_shots {
				err := a.HitOrMiss(strings.ToLower(v))
				if err != nil {
					errorchan <- err
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) GetBoard() error {
	coords, err := a.client.Board()
	if err != nil {
		return err
	}
	for _, x := range coords {
		coordmap, err := numericCords(x)
		if err != nil {
			return err
		}
		a.myStates[coordmap["x"]][coordmap["y"]] = "Ship"

	}
	return nil
}

func (a *App) Timer(ctx context.Context, timeLeftchan, resetTimerchan chan int) {
	second := time.NewTicker(time.Second)
	syncTimer := time.NewTicker(5 * time.Second)
	timeLeft := 60
	for {
		select {
		case <-second.C:
			timeLeftchan <- timeLeft
			timeLeft -= 1
		case <-syncTimer.C:
			timeLeft = a.actualStatus.Timer
		case <-resetTimerchan:
			second.Reset(time.Second)
			timeLeft = 60
			timeLeftchan <- timeLeft
		case <-ctx.Done():
			return
		}
	}
}
