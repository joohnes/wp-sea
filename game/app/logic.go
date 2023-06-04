package app

import (
	"context"
	"fmt"
	"github.com/joohnes/wp-sea/game/logger"
	"strings"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
	"github.com/joohnes/wp-sea/game/helpers"

	"github.com/fatih/color"
)

const (
	waitDuration = 3
)

func (a *App) WaitForStart() error {
	for {
		status, err := a.client.Status()
		if err != nil {
			return err
		}
		if status.GameStatus == "game_in_progress" {
			a.actualStatus = *status
			return nil
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
		if status.ShouldFire {
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
	switch status.LastGameStatus {
	case "win":
		green := color.New(color.FgBlack, color.BgGreen).SprintFunc()
		fmt.Println(green("You have won the game!"))
		a.gameState = StateEnded
		return true
	case "lose":
		red := color.New(color.FgBlack, color.BgRed).SprintFunc()
		fmt.Println(red("You have lost the game!"))
		a.gameState = StateEnded
		return true
	}
	return false
}

func (a *App) Shoot(coord string) error {
	coordmap, err := helpers.NumericCords(coord)
	if err != nil {
		return err
	}
	result, err := a.client.Shoot(coord)
	if err != nil {
		return err
	}

	switch result {
	case "miss":
		a.enemyStates[coordmap["x"]][coordmap["y"]] = "Miss"
		a.shotsCount += 1
		a.gameState = StateOppTurn
	case "hit":
		a.enemyStates[coordmap["x"]][coordmap["y"]] = "Hit"
		a.shotsCount += 1
		a.shotsHit += 1
		a.LastPlayerHit = coord
		a.mode = HuntState
		a.statistics[coord] += 1
	case "sunk":
		a.enemyStates[coordmap["x"]][coordmap["y"]] = "Hit"
		a.shotsCount += 1
		a.shotsHit += 1
		a.MarkBorders(coordmap)
		a.mode = TargetState
		a.statistics[coord] += 1
	}
	a.playerShots[coord] = result
	return nil
}

func (a *App) Play(ctx context.Context, coordchan <-chan string, textchan chan<- string, errorchan chan error, resetTime chan int) {
	var coord string
	for {
		select {
		case coord = <-coordchan:
			if a.gameState == StatePlayerTurn {
				_, err := helpers.NumericCords(coord)
				if err != nil {
					errorchan <- err
				}
				if a.playerShots[coord] != "" {
					textchan <- fmt.Sprintf("You have already shot at %s", coord)
					break
				}
				err = a.Shoot(coord)
				if err != nil {
					errorchan <- err
				}

				resetTime <- 1
				textchan <- fmt.Sprintf("Shot at %s", coord)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) HitOrMiss(coord string) error {
	coordmap, err := helpers.NumericCords(coord)
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
	return nil
}

func (a *App) CheckStatus(ctx context.Context, cancel context.CancelFunc, textchan chan<- string) {
	statusTicker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-statusTicker.C:
			status, err := a.client.Status()
			if err != nil {
				continue
			}
			switch status.ShouldFire {
			case true:
				a.gameState = StatePlayerTurn
			case false:
				a.gameState = StateOppTurn
			}
			a.actualStatus = *status
			if status.GameStatus == "ended" {
				a.gameState = StateEnded
				switch status.LastGameStatus {
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
			diff := helpers.Difference(a.actualStatus.OppShots, a.oppShots)
			a.oppShots = a.actualStatus.OppShots
			if len(diff) != 0 {
				for _, v := range diff {
					err := a.HitOrMiss(strings.ToLower(v))
					if err != nil {
						errorchan <- err
					}
				}
				a.gameState = StatePlayerTurn
			}
			// additional check
		case <-checkTicker.C:
			for _, v := range a.actualStatus.OppShots {
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
		coordmap, err := helpers.NumericCords(x)
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
			if timeLeft == 0 {
				second.Stop()
				syncTimer.Stop()
			}
			timeLeft -= 1
		case <-syncTimer.C:
			timeLeft = a.actualStatus.Timer
			second.Reset(time.Second)
		case <-resetTimerchan:
			second.Reset(time.Second)
			timeLeft = 60
			timeLeftchan <- timeLeft
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) Reset() {
	a.oppShots = []string{}
	a.oppNick = ""
	a.oppDesc = ""
	a.actualStatus = StatusResponse{}
	a.shotsCount = 0
	a.shotsHit = 0
	a.myStates = [10][10]gui.State{}
	a.enemyStates = [10][10]gui.State{}
	a.gameState = StateStart
	a.enemyShips = map[int]int{4: 1, 3: 2, 2: 3, 1: 4}
	a.playerShots = map[string]string{}
	a.mode = TargetState
	a.LastPlayerHit = ""
	a.algorithm = false
	a.algorithmTried = []string{}
	a.client.ResetToken()
}

// TranslateMap Translates [][]gui.State to coord list, so we can send it in a request
func (a *App) TranslateMap() (coords []string) {
	for i, x := range a.playerStates {
		for j, y := range x {
			if y != "" {
				coords = append(coords, helpers.AlphabeticCoords(i, j))
			}
		}
	}
	return
}

func (a *App) LogError(err error) {
	logger.GetLoggerInstance().Println(err)
	if ShowErrors {
		fmt.Println(err)
	}
}
