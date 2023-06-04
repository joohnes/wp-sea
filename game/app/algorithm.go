package app

import (
	"context"
	"encoding/csv"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"github.com/joohnes/wp-sea/game/helpers"
	"github.com/joohnes/wp-sea/game/logger"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Mode string

const (
	HuntState   Mode = "Hunt"
	TargetState      = "Target"
)

func (a *App) LoadStatistics() {
	f, err := os.Open("statistics.csv")
	defer f.Close()
	if err != nil {
		_, err = os.Create("statistics.csv")
		if err != nil {
			return
		}
		f, err = os.Open("statistics.csv")
		defer f.Close()
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		logger.GetLoggerInstance().Println("couldn't read statistics")
	}
	var basemap map[string]int
	for _, x := range records {
		number, err := strconv.Atoi(x[1])
		if err != nil {
			logger.GetLoggerInstance().Printf("couldn't load %v\n", x)
		}
		basemap[x[0]] = number
	}
	a.statistics = basemap
}

func (a *App) SaveStatistics() {
	f, err := os.Open("statistics.csv")
	defer f.Close()
	if err != nil {
		_, err = os.Create("statistics.csv")
		if err != nil {
			return
		}
		f, err = os.Open("statistics.csv")
		defer f.Close()
	}

	w := csv.NewWriter(f)
	for coord, occurrences := range a.statistics {
		err := w.Write([]string{coord, strconv.Itoa(occurrences)})
		if err != nil {
			logger.GetLoggerInstance().Printf("couldn't save %s, %d\n", coord, occurrences)
		}
	}
	w.Flush()
}

func (a *App) ClosestShip(x, y int) int {
	vec := []point{
		{-1, 0},
		{-1, -1},
		{0, 1},
		{1, 1},
		{1, 0},
		{-1, 1},
		{0, -1},
		{1, -1},
	}

	for i := 1; i <= 10; i++ {
		for _, v := range vec {
			dx := x + i*v.x
			dy := y + i*v.y
			if a.enemyStates[dx][dy] == gui.Hit {
				return i
			}
		}
	}
	return 0
}

func (a *App) SearchShip() (x, y int) {
	if a.mode == TargetState {
		for {
			x = rand.Intn(10)
			y = rand.Intn(10)
			if a.enemyStates[x][y] != gui.Hit && a.enemyStates[x][y] != gui.Miss {
				break
			}
		}
		return
	} else {
		coord, err := helpers.NumericCords(a.LastPlayerHit)
		if err != nil {
			return
		}
		vec := []point{
			{-1, 0},
			{0, 1},
			{1, 0},
			{0, -1},
		}
		for _, v := range vec {
			dx := int(coord["x"]) + v.x
			dy := int(coord["y"]) + v.y
			if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
				continue
			}
			if a.enemyStates[dx][dy] != gui.Hit && a.enemyStates[dx][dy] != gui.Miss {
				return dx, dy
			}
		}
		a.algorithmTried = append(a.algorithmTried, a.LastPlayerHit)
		for _, v := range vec {
			dx := int(coord["x"]) + v.x
			dy := int(coord["y"]) + v.y
			if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
				continue
			}
			if a.enemyStates[dx][dy] == gui.Hit && !In(a.algorithmTried, a.LastPlayerHit) {
				a.LastPlayerHit = helpers.AlphabeticCoords(dx, dy)
				return a.SearchShip()
			}
		}
	}
	return
}
func In(arr []string, coord string) bool {
	for _, x := range arr {
		if x == coord {
			return true
		}
	}
	return false
}

func (a *App) AlgorithmPlay(ctx context.Context, textchan chan<- string, errorchan chan error, resetTime chan int) {
	t := time.NewTicker(time.Millisecond * 200)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if a.gameState == StatePlayerTurn {
				x, y := a.SearchShip()
				coord := helpers.AlphabeticCoords(x, y)
				err := a.Shoot(coord)
				if err != nil {
					errorchan <- err
				}
				resetTime <- 1
				textchan <- fmt.Sprintf("%s: Algorithm shot at %s", a.mode, coord)
			} else if a.gameState == StateEnded {
				return
			}
		}
	}
}
