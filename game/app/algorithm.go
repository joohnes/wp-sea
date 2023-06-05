package app

import (
	"context"
	"fmt"
	"math/rand"
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
	enabled bool
	mode    Mode
	Loop    bool
	tried   []string
	//rest of the options
}

func NewAlgorithm() Algorithm {
	return Algorithm{
		false,
		TargetState,
		false,
		[]string{},
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

func (a *App) SearchShip() (x, y int) {
	if a.algorithm.mode == TargetState {
		for {
			x = rand.Intn(10)
			y = rand.Intn(10)
			if a.enemyStates[x][y] != gui.Hit && a.enemyStates[x][y] != gui.Miss {
				break
			}
		}
		return
	} else {
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
		for _, x := range a.CheckShipPoints(coordX, coordY) {
			cord := helpers.AlphabeticCoords(x.x, x.y)
			if !In(a.algorithm.tried, cord) {
				a.algorithm.tried = append(a.algorithm.tried, a.LastPlayerHit)
				a.LastPlayerHit = cord
				return a.SearchShip()
			}
		}
	}
	for shot := range a.playerShots {
		if shot == "a1" {
			for {
				x = rand.Intn(10)
				y = rand.Intn(10)
				if a.enemyStates[x][y] != gui.Hit && a.enemyStates[x][y] != gui.Miss {
					break
				}
				a.algorithm.mode = TargetState
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
				var coord string
				if a.gameState != StatePlayerTurn {
					return
				}
				coord = helpers.AlphabeticCoords(x, y)
				err := a.Shoot(coord)
				if err != nil {
					errorchan <- err
				}
				resetTime <- 1
				textchan <- fmt.Sprintf("%s: Algorithm shot at %s", a.algorithm.mode, coord)
			} else if a.gameState == StateEnded {
				return
			}
		}
	}
}

func (a *App) HasAlreadyBeenShot(coord string) bool {
	for x := range a.playerShots {
		if x == coord {
			return true
		}
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
		if a.enemyStates[dx][dy] == "Ship" {
			connections = append(connections, point{dx, dy})
		}

	}
	for _, c := range connections {
		a.getShips(c.x, c.y, points)
	}
}
