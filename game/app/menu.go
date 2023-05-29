package app

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/inancgumus/screen"
	table "github.com/jedib0t/go-pretty/v6/table"
	"github.com/joohnes/wp-sea/game/helpers"
	"github.com/joohnes/wp-sea/game/logger"
)

var ErrBack = errors.New("back")

func (a *App) ShowStats() error {
	var data map[string][]int
	err := helpers.ServerErrorWrapper(ShowErrors, func() error {
		var err error
		data, err = a.client.Stats()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	p := make(helpers.PairList, len(data))

	i := 0
	for k, v := range data {
		p[i] = helpers.Pair{Key: k, Values: v}
		i++
	}

	sort.Sort((sort.Reverse(p)))

	t := table.NewWriter()
	t.SetTitle("Stats")

	t.AppendHeader(table.Row{"#", "Nick", "Games", "Points", "Rank", "Wins"})
	counter := 1
	for _, x := range p {
		t.AppendRow(table.Row{counter, x.Key, x.Values[0], x.Values[1], x.Values[2], x.Values[3]})
		counter += 1
	}
	fmt.Println(t.Render())
	fmt.Println("Press enter to go back to the menu")
	_, _ = fmt.Scanln()
	return nil
}
func (a *App) ShowPlayerStats(nick string) error {
	data, err := a.client.StatsPlayer(nick)
	if err != nil {
		return err
	}
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("%s's stats", nick))

	t.AppendHeader(table.Row{"Nick", "Games", "Points", "Rank", "Wins"})

	t.AppendRow(table.Row{nick, data[0], data[1], data[2], data[3]})
	fmt.Println(t.Render())
	fmt.Println("Press enter to go back to the menu")
	_, _ = fmt.Scanln()
	return nil
}

func (a *App) WaitingRefresh() {
	for {
		if a.gameState == StateWaiting {
			err := a.client.Refresh()
			if err != nil {
				time.Sleep(10 * time.Second)
			}
		} else {
			return
		}
	}
}

func PrintOptions(nick string) {
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("Nick: %s", nick))
	t.AppendHeader(table.Row{"#", "Choose an option"})
	t.AppendRow(table.Row{1, "Play with WPBot"})
	t.AppendRow(table.Row{2, "Play with another player"})
	t.AppendRow(table.Row{3, "Top 10 players"})
	t.AppendRow(table.Row{4, "Your stats"})
	t.AppendRow(table.Row{5, "Check someone's stats"})
	t.AppendRow(table.Row{6, "Set up your ships"})
	t.AppendFooter(table.Row{"", "Type 'q' to exit"})
	fmt.Println(t.Render())
	fmt.Print("Option: ")
}

func (a *App) ChoosePlayer() error {
	var playerlist []map[string]string
	log := logger.GetLoggerInstance()
	err := helpers.ServerErrorWrapper(ShowErrors, func() error {
		var err error
		playerlist, err = a.client.PlayerList()
		if err != nil {
			log.Println(err.Error())
			if ShowErrors {
				fmt.Println(err.Error())
			}
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(playerlist) != 0 {
		t := table.NewWriter()
		t.SetTitle("Waiting players")

		t.AppendHeader(table.Row{"#", "Nick"})
		for counter, x := range playerlist {
			t.AppendRow(table.Row{counter + 1, x["nick"]})
		}
		t.AppendFooter(table.Row{"", fmt.Sprintf("Choose a player (1-%d)\nIf you wish to wait, type 'wait'\nTo go back, type 'back'", len(playerlist))})
		fmt.Println(t.Render())

	Again:
		answer, err := helpers.GetAnswer()
		if err != nil {
			return err
		}
		if strings.ToLower(answer) == "wait" {
			fmt.Println("Waiting...")
			err = a.client.InitGame(nil, a.desc, a.nick, "", false)
			if err != nil {
				return err
			}
			a.gameState = StateWaiting

			go a.WaitingRefresh()

			return nil
		} else if strings.ToLower(answer) == "back" {
			return ErrBack
		} else {
			i, err := strconv.Atoi(answer)
			if err != nil {
				log.Printf("Couldn't convert %s to a number\n", answer)
				fmt.Println("Please enter a valid number (1-", len(playerlist), ")")
				goto Again
			}
			if i < 1 || i > len(playerlist)+1 {
				fmt.Println("Please enter a valid number (1-", len(playerlist), ")")
				goto Again
			}
			err = a.client.InitGame(nil, a.desc, a.nick, playerlist[i-1]["nick"], false)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		fmt.Println("No players waiting at the moment")
		fmt.Println("Do you want to wait for another player? y/n")
	NoPlayersAgain:
		answer, err := helpers.GetAnswer()
		if err != nil {
			return err
		}
		switch strings.ToLower(answer) {
		case "y", "yes":
			fmt.Println("Waiting...")
			// err = a.client.InitGame(nil, a.desc, a.nick, "", false)
			// if err != nil {
			// 	return err
			// }
			err := helpers.ServerErrorWrapper(ShowErrors, func() error {
				err := a.client.InitGame(nil, a.desc, a.nick, "", false)
				if err != nil {
					if ShowErrors {
						fmt.Println(err.Error())
					}
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
			a.gameState = StateWaiting

			go a.WaitingRefresh()

			return nil
		case "n", "no":
			return errors.New("back")
		default:
			fmt.Println("Please type 'y' or 'n'")
			goto NoPlayersAgain
		}
	}
}

func (a *App) ChooseOption() error {
	log := logger.GetLoggerInstance()
Start:
	screen.Clear()
	screen.MoveTopLeft()
	PrintOptions(a.nick)
	answer, err := helpers.GetAnswer()
	if err != nil {
		log.Println(err)
		goto Start
	}

	switch strings.ToLower(answer) {
	case "q", "quit":
		os.Exit(0)
	case "1": // play with bot
		err := helpers.ServerErrorWrapper(ShowErrors, func() error {
			err := a.client.InitGame(nil, a.desc, a.nick, "", true)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		a.gameState = StateWaiting

	case "2": // play with another player
		for {
			err := a.ChoosePlayer()
			if errors.Is(err, ErrBack) {
				goto Start
			}
			if err != nil {
				log.Println(err)
				fmt.Println("Error occurred. Please try again")
				continue
			}
			break
		}
	case "3": // top10
		err := a.ShowStats()
		if err != nil {
			return err
		}
		goto Start
	case "4": // stats
		err := a.ShowPlayerStats(a.nick)
		if err != nil {
			return err
		}
		goto Start

	case "5":
		fmt.Print("Enter name: ")
		nick, err := helpers.GetAnswer()
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			goto Start
		}
		err = a.ShowPlayerStats(nick)
		if err != nil {
			return err
		}
		goto Start
	case "6": //set up ships
	default:
		fmt.Println("Please enter a valid number!")
		goto Start
	}
	return nil
}
