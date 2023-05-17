package app

import (
	"fmt"
	"strconv"
	"time"
)

func (a *App) ChooseOption() error {
	fmt.Println("1. Play with WPBot")
	fmt.Println("2. Play with another player")
	fmt.Println("3. Top 10")
	fmt.Println("4. Check your stats")
	fmt.Println("5. Set up your ships")
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
				fmt.Println("Waiting...")
				err = a.client.InitGame(nil, a.desc, a.nick, "", false)
				if err != nil {
					return err
				}
				go func() {
					if a.actualStatus.Game_status == "waiting" {
						_ = a.client.Refresh()
						time.Sleep(10 * time.Second)
					} else {
						return
					}
				}()
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
	case "3": // top10
		table, err := a.client.Stats()
		if err != nil {
			return err
		}
		fmt.Println("Top 10 players")
		fmt.Println("Nick | Games | Points | Rank | Wins")
		for i, x := range table {
			text := i
			for _, y := range x {
				text += string(y) + "   "
			}
			fmt.Println(text)
		}
	case "4": // stats
		table, err := a.client.StatsPlayer(a.nick)
		if err != nil {
			return err
		}
		fmt.Println("Stats for player: ", a.nick)
		fmt.Println("Games | Points | Rank | Wins")
		text := ""
		for _, x := range table {
			text += string(x) + " "
		}
	case "5": //set up ships
	}
	return nil
}
