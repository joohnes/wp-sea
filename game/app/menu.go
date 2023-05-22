package app

import (
	"fmt"
	"strconv"
	"time"

	table "github.com/jedib0t/go-pretty/v6/table"
)

func (a *App) ShowStats() error {
	data, err := a.client.Stats()
	if err != nil {
		return err
	}
	t := table.NewWriter()
	t.SetTitle("Stats")

	t.AppendHeader(table.Row{"#", "Nick", "Games", "Points", "Rank", "Wins"})
	counter := 1
	for i, x := range data {
		t.AppendRow(table.Row{counter, i, x[0], x[1], x[2], x[3]})
		counter += 1
	}
	fmt.Println(t.Render())
	return nil
}
func (a *App) ShowPlayerStats() error {
	data, err := a.client.StatsPlayer(a.nick)
	if err != nil {
		return err
	}
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("%s's stats", a.nick))

	t.AppendHeader(table.Row{"Nick", "Games", "Points", "Rank", "Wins"})

	t.AppendRow(table.Row{a.nick, data[0], data[1], data[2], data[3]})
	fmt.Println(t.Render())
	return nil
}

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
		err := a.ShowStats()
		if err != nil {
			return err
		}
	case "4": // stats
		err := a.ShowPlayerStats()
		if err != nil {
			return err
		}

	case "5": //set up ships
	}
	return nil
}
