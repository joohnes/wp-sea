package app

import (
	"context"
	"fmt"
	"github.com/joohnes/wp-sea/game/logger"
	"math/rand"
	"strings"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
	"github.com/joohnes/wp-sea/game/helpers"
)

type Mode string

const (
	HuntState   Mode = "Hunt"
	TargetState Mode = "Target"
)

type Algorithm struct {
	enabled  bool
	mode     Mode
	tried    []string
	shot     []string
	statList PairList
	options  Options
}

type Options struct {
	Loop  bool
	Stats bool
}

func NewAlgorithm() Algorithm {
	return Algorithm{
		false,
		TargetState,
		[]string{},
		[]string{},
		PairList{},
		Options{
			false,
			false,
		},
	}
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

func (a *App) getRandomCoord() (x, y int) {
	for {
		x = rand.Intn(10)
		y = rand.Intn(10)
		if a.enemyStates[x][y] != gui.Hit && a.enemyStates[x][y] != gui.Miss && !In(a.algorithm.shot, helpers.AlphabeticCoords(x, y)) {
			break
		}
		return
	}
	return x, y
}

func (a *App) SearchShip() (x, y int) {
	if a.algorithm.mode == TargetState {
		if a.algorithm.options.Stats {
			for _, v := range a.algorithm.statList {
				x, y, err := helpers.NumericCords(v.Key)
				if err != nil {
					return a.getRandomCoord()
				}
				a.algorithm.statList = a.algorithm.statList[1:]
				if a.HasAlreadyBeenShot(helpers.AlphabeticCoords(x, y)) {
					continue
				}
				return x, y
			}
		} else {
			return a.getRandomCoord()
		}
	} else if a.algorithm.mode == HuntState {
		coordX, coordY, err := helpers.NumericCords(a.LastPlayerHit)
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
			dx := coordX + v.x
			dy := coordY + v.y
			if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
				continue
			}
			if a.enemyStates[dx][dy] != gui.Hit && a.enemyStates[dx][dy] != gui.Miss {
				return dx, dy
			}
		}
		a.algorithm.tried = append(a.algorithm.tried, a.LastPlayerHit)
		for _, v := range a.CheckShipPoints(coordX, coordY) {
			cord := helpers.AlphabeticCoords(v.x, v.y)
			if !In(a.algorithm.tried, cord) {
				a.LastPlayerHit = cord
				return a.SearchShip()
			}
		}
	}
	return a.getRandomCoord()
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
	a.algorithm.statList = a.getSortedStatistics()
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if a.gameState == StatePlayerTurn {
				x, y := a.SearchShip()
				var coord string
				coord = helpers.AlphabeticCoords(x, y)
				err := a.Shoot(coord)
				if err != nil {
					errorchan <- err
				} else {
					a.algorithm.shot = append(a.algorithm.shot, strings.ToLower(coord))
					resetTime <- 1
					textchan <- fmt.Sprintf("%s: Algorithm shot at %s", a.algorithm.mode, strings.ToUpper(coord))
				}
			} else if a.gameState == StateEnded {
				return
			}
		}
	}
}

func (a *App) HasAlreadyBeenShot(coord string) bool {
	coordX, coordY, err := helpers.NumericCords(coord)
	if err != nil {
		logger.GetLoggerInstance().Error.Println("couldn't convert coords")
		return false
	}
	if a.enemyStates[coordX][coordY] == gui.Hit || a.enemyStates[coordX][coordY] == gui.Miss {
		return true
	}
	return false
}

func (a *App) CheckShipPoints(x, y int) []point {
	var points []point
	a.getShips(x, y, &points)
	return points
}

func (a *App) getShips(x, y int, points *[]point) {
	vec := []point{
		{-1, 0},
		{0, 1},
		{1, 0},
		{0, -1},
	}

	for _, i := range *points {
		if i.x == x && i.y == y {
			return
		}
	}
	*points = append(*points, point{x, y})
	var connections []point

	for _, v := range vec {
		dx := x + v.x
		dy := y + v.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
			continue
		}
		if a.enemyStates[dx][dy] == gui.Hit {
			connections = append(connections, point{dx, dy})
		}

	}
	for _, c := range connections {
		a.getShips(c.x, c.y, points)
	}
}
