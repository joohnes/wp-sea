package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/inancgumus/screen"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/joohnes/wp-sea/game/helpers"
	"github.com/joohnes/wp-sea/game/logger"
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

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

	t := table.NewWriter()
	t.SetTitle("Stats")
	t.SortBy([]table.SortBy{{Name: "Rank", Mode: table.AscNumeric}})

	t.AppendHeader(table.Row{"#", "Nick", "Games", "Points", "Rank", "Wins"})
	counter := 1
	for i, x := range data {
		t.AppendRow(table.Row{counter, i, x[0], x[1], x[2], x[3]})
		counter++
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
func (a *App) WaitingTimer() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	dur := 0
	fmt.Print("\033[s")
	for {
		if a.gameState != StateWaiting {
			return
		}
		<-t.C
		fmt.Print("\033[u\033[K")
		fmt.Printf("Waiting [%v seconds]", dur)
		dur++
	}
}

func (a *App) WaitingRefresh() {
	go a.WaitingTimer()
	for {
		if a.gameState == StateWaiting {
			err := a.client.Refresh()
			if err != nil {
				logger.GetLoggerInstance().Error.Println(err)
			}
			time.Sleep(10 * time.Second)
		} else {
			return
		}
	}
}

func PrintOptions(nick string, changed, algorithm bool) {
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("Nick: %s", nick))
	t.AppendHeader(table.Row{"#", "Choose an option"})
	t.AppendRow(table.Row{1, "Play with WPBot"})
	t.AppendRow(table.Row{2, "Play with another player"})
	t.AppendRow(table.Row{3, "Top 10 players"})
	t.AppendRow(table.Row{4, "Your stats"})
	t.AppendRow(table.Row{5, "Check someone's stats"})
	if changed {
		green := color.New(color.FgGreen).SprintFunc()
		t.AppendRow(table.Row{6, fmt.Sprintf("Set up your ships %s", green("(SET)"))})
	} else {
		t.AppendRow(table.Row{6, "Set up your ships"})
	}
	t.AppendRow(table.Row{7, "Reset ship placement"})
	if !algorithm {
		green := color.New(color.FgGreen).SprintFunc()
		t.AppendRow(table.Row{8, green("Turn on Algorithm")})
	} else {
		red := color.New(color.FgRed).SprintFunc()
		t.AppendRow(table.Row{8, red("Turn off Algorithm")})
	}
	t.AppendRow(table.Row{9, "Show algorithm options"})
	t.AppendRow(table.Row{10, "Show heatmap"})
	t.AppendFooter(table.Row{"", "Type 'q' to exit"})
	fmt.Println(t.Render())
	fmt.Print("Option: ")
}

func (a *App) PrintAlgorithmOptions() error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	for {
		screen.Clear()
		screen.MoveTopLeft()
		t := table.NewWriter()
		t.SetTitle("Algorithm options")
		t.AppendHeader(table.Row{"#", fmt.Sprintf("Choose an option [%s/%s]", green("enabled"), red("disabled"))})
		if a.algorithm.options.Loop {
			t.AppendRow(table.Row{1, fmt.Sprintf("%s: Auto play, ctrl + c to quit", green("Loop"))})
		} else {
			t.AppendRow(table.Row{1, fmt.Sprintf("%s: Auto play, ctrl + c to quit", red("Loop"))})
		}
		if a.algorithm.options.Stats {
			t.AppendRow(table.Row{2, fmt.Sprintf("%s: Uses statistics to shot", green("Stat"))})
			t.AppendRow(table.Row{"", "at the most common spot for ship to be in"})
		} else {
			t.AppendRow(table.Row{2, fmt.Sprintf("%s: Uses statistics to shot", red("Stat"))})
			t.AppendRow(table.Row{"", "at the most common spot for ship to be in"})
		}
		if a.algorithm.options.Density {
			t.AppendRow(table.Row{3, fmt.Sprintf("%s: Uses density map to shot", green("Density"))})
			t.AppendRow(table.Row{"", "where there may be the most remaining ships"})
		} else {
			t.AppendRow(table.Row{3, fmt.Sprintf("%s: Uses density map to shot", red("Density"))})
			t.AppendRow(table.Row{"", "where there may be the most remaining ships"})
		}
		t.AppendFooter(table.Row{"", "Type 'b' to go back"})
		fmt.Println(t.Render())
		fmt.Print("Option: ")
		answer, err := helpers.GetAnswer(false)
		if err != nil {
			return err
		}
		switch strings.ToLower(answer) {
		case "1":
			a.algorithm.options.Loop = !a.algorithm.options.Loop
			continue
		case "2":
			a.algorithm.options.Stats = !a.algorithm.options.Stats
			a.algorithm.options.Density = false
		case "3":
			a.algorithm.options.Stats = false
			a.algorithm.options.Density = !a.algorithm.options.Density
		case "b", "back":
			return ErrBack
		default:
			fmt.Println("Please enter a valid number!")
			time.Sleep(time.Second)
			continue
		}
	}
}

func (a *App) ChoosePlayer() error {
	var playerlist []map[string]string
	log := logger.GetLoggerInstance()
	err := helpers.ServerErrorWrapper(ShowErrors, func() error {
		var err error
		playerlist, err = a.client.PlayerList()
		if err != nil {
			logger.GetLoggerInstance().Error.Println(err)
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
		answer, err := helpers.GetAnswer(false)
		if err != nil {
			return err
		}
		if strings.ToLower(answer) == "wait" {
			var err error
			if a.CheckIfChangedMap() && a.Requirements() {
				err = helpers.ServerErrorWrapper(ShowErrors, func() error {
					err := a.client.InitGame(a.TranslateMap(), a.desc, a.nick, "", false)
					if err != nil {
						return err
					}
					return nil
				})
				fmt.Println("Connecting to server...")
			} else {
				err = helpers.ServerErrorWrapper(ShowErrors, func() error {
					err := a.client.InitGame(nil, a.desc, a.nick, "", false)
					if err != nil {
						return err
					}
					return nil
				})
				fmt.Println("Connecting to server...")
			}
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
				log.Error.Printf("Couldn't convert %s to a number\n", answer)
				fmt.Println("Please enter a valid number (1-", len(playerlist), ")")
				goto Again
			}
			if i < 1 || i > len(playerlist)+1 {
				fmt.Println("Please enter a valid number (1-", len(playerlist), ")")
				goto Again
			}
			if a.CheckIfChangedMap() && a.Requirements() {
				err = helpers.ServerErrorWrapper(ShowErrors, func() error {
					err := a.client.InitGame(a.TranslateMap(), a.desc, a.nick, playerlist[i-1]["nick"], false)
					if err != nil {
						return err
					}
					return nil
				})
				fmt.Println("Connecting to server...")
			} else {
				err = helpers.ServerErrorWrapper(ShowErrors, func() error {
					err := a.client.InitGame(nil, a.desc, a.nick, playerlist[i-1]["nick"], false)
					if err != nil {
						return err
					}
					return nil
				})
				fmt.Println("Connecting to server...")
			}
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		t := table.NewWriter()
		t.SetTitle("No players waiting at the moment")
		t.AppendRow(table.Row{"Do you want to wait for another player? [y/n]"})
		fmt.Println(t.Render())
	NoPlayersAgain:
		answer, err := helpers.GetAnswer(false)
		if err != nil {
			return err
		}
		switch strings.ToLower(answer) {
		case "y", "yes":
			if a.CheckIfChangedMap() && a.Requirements() {
				err = helpers.ServerErrorWrapper(ShowErrors, func() error {
					err := a.client.InitGame(a.TranslateMap(), a.desc, a.nick, "", false)
					if err != nil {
						return err
					}
					return nil
				})
				fmt.Println("Connecting to server...")
			} else {
				err = helpers.ServerErrorWrapper(ShowErrors, func() error {
					err := a.client.InitGame(nil, a.desc, a.nick, "", false)
					if err != nil {
						return err
					}
					return nil
				})
				fmt.Println("Connecting to server...")
			}
			if err != nil {
				return err
			}
			a.gameState = StateWaiting

			go a.WaitingRefresh()

			return nil
		case "n", "no":
			return ErrBack
		default:
			fmt.Println("Please type '(y)es' or '(n)o'")
			goto NoPlayersAgain
		}
	}
}

func (a *App) ChooseOption(ctx context.Context, shipchannel chan string, errChan chan error) error {
	log := logger.GetLoggerInstance()
	if a.algorithm.options.Loop {
		err := helpers.ServerErrorWrapper(ShowErrors, func() error {
			var err error
			if a.CheckIfChangedMap() && a.Requirements() {
				err = a.client.InitGame(a.TranslateMap(), a.desc, a.nick, "", true)
			} else {
				err = a.client.InitGame(nil, a.desc, a.nick, "", true)
			}
			if err != nil {
				return err
			}
			fmt.Println("Connecting to server...")
			return nil
		})
		if err != nil {
			return err
		}
		a.gameState = StateWaiting
		return nil
	}

	var Break = false
	for {
		if Break {
			break
		}
		screen.Clear()
		screen.MoveTopLeft()
		PrintOptions(a.nick, a.Requirements(), a.algorithm.enabled)
		answer, err := helpers.GetAnswer(false)
		if err != nil {
			log.Error.Println(err)
			continue
		}

		switch strings.ToLower(answer) {
		case "q", "quit":
			os.Exit(0)
		case "1": // play with bot
			err := helpers.ServerErrorWrapper(ShowErrors, func() error {
				var err error
				if a.CheckIfChangedMap() && a.Requirements() {
					err = a.client.InitGame(a.TranslateMap(), a.desc, a.nick, "", true)
				} else {
					err = a.client.InitGame(nil, a.desc, a.nick, "", true)
				}
				if err != nil {
					return err
				}
				fmt.Println("Connecting to server...")
				return nil
			})
			if err != nil {
				return err
			}
			a.gameState = StateWaiting
			Break = true

		case "2": // play with another player
			err := a.ChoosePlayer()
			if errors.Is(err, ErrBack) {
				continue
			}
			if err != nil {
				logger.GetLoggerInstance().Error.Println(err)
			}
			Break = true

		case "3": // top10
			err := a.ShowStats()
			if err != nil {
				return err
			}

		case "4": // stats
			err := a.ShowPlayerStats(a.nick)
			if err != nil {
				return err
			}

		case "5": // specific player stats
			fmt.Print("Enter name: ")
			nick, err := helpers.GetAnswer(true)
			if err != nil {
				log.Error.Println(err)
				fmt.Println(err)
				continue
			}
			err = a.ShowPlayerStats(nick)
			if err != nil {
				return err
			}

		case "6": //set up ships
			go a.PlaceShips(ctx, shipchannel, errChan)
			a.SetUpShips(ctx, shipchannel, errChan)

		case "7": //reset ship placement
			for i := range a.playerStates {
				for j := range a.playerStates[i] {
					a.playerStates[i][j] = ""
				}
			}
			a.placeShips = map[int]int{4: 1, 3: 2, 2: 3, 1: 4}
			fmt.Println("Ship placement has been reset!")
			time.Sleep(time.Second * 2)

		case "8": // turn on/off algorithm
			a.algorithm.enabled = !a.algorithm.enabled

		case "9": // algorithm options
			err = a.PrintAlgorithmOptions()
			if errors.Is(err, ErrBack) {
				continue
			}
			if err != nil {
				logger.GetLoggerInstance().Error.Println(err)
			}
		case "10": // Show heatmap with shot statistics
			a.ShowStatistics()
		default:
			fmt.Println("Please enter a valid number!")
			time.Sleep(time.Second)
		}
	}
	return nil
}
