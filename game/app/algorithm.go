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
	HuntState    Mode = "HUNT"
	TargetState  Mode = "TARGET"
	DensityState Mode = "DENSITY"
	MixedState   Mode = "MIXED"
)

type Algorithm struct {
	enabled       bool
	assistance    bool
	mode          Mode
	tried         []string
	shot          []string
	statList      PairList
	densityMap    [10][10]int
	predictedShot string
	options       Options
}

type Options struct {
	Loop    bool
	Stats   bool
	Density bool
	Mixed   bool
	Map     bool
}

/*
	This algorithm was created with the help of the article:
	http://www.datagenetics.com/blog/december32011/
*/

func NewAlgorithm() Algorithm {
	return Algorithm{
		false,
		false,
		TargetState,
		[]string{},
		[]string{},
		PairList{},
		[10][10]int{},
		"",
		Options{
			false,
			false,
			false,
			false,
			true,
		},
	}
}

func (a *App) LoadMap() {
	basemap := []string{"a1", "a5", "a3", "a7", "a9", "b1", "b3", "b5", "b7", "b9", "c1", "c3", "c5", "d1", "i9", "j1", "j3", "j5", "j7", "j9"}
	for _, b := range basemap {
		x, y, err := helpers.NumericCords(b)
		if err != nil {
			logger.GetLoggerInstance().Error.Println("couldn't convert coord during LoadMap func")
		} else {
			a.playerStates[x][y] = gui.Ship
		}
	}
	for x := range a.placeShips {
		a.placeShips[x] = 0
	}
}

func (a *App) getRandomCoord() (x, y int) {
	for {
		x = rand.Intn(10)
		y = rand.Intn(10)
		if a.enemyStates[x][y] != gui.Hit && a.enemyStates[x][y] != gui.Miss && !In(a.algorithm.shot, helpers.AlphabeticCoords(x, y)) {
			return
		}
	}
}

func (a *App) SearchShip() (x, y int) {
	if a.algorithm.mode == TargetState {
		return a.getTargetCoords()
	} else if a.algorithm.mode == DensityState || a.algorithm.mode == MixedState {
		return a.getDensityCoords()
	} else if a.algorithm.mode == HuntState {
		return a.getHuntCoords()
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

func (a *App) getTargetCoords() (x, y int) {
	if a.algorithm.options.Stats || len(a.algorithm.shot) < 10 {
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
	}
	return a.getRandomCoord()
}

func (a *App) getHuntCoords() (x, y int) {
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
	return a.getRandomCoord()
}

func (a *App) getDensityCoords() (x, y int) {
	densityMap := [10][10]int{}
	for i := range a.enemyStates {
		for j, v := range a.enemyStates[i] {
			if v != gui.Hit && v != gui.Miss {
				for _, shape := range shapes {
					if a.enemyShips[len(shape)+1] == 0 {
						continue
					}
					if a.DoesShapeFit(i, j, shape) {
						for _, c := range shape {
							dx := i + c.x
							dy := j + c.y
							densityMap[dx][dy] += 1
						}
					}
				}
			}
		}
	}
	a.algorithm.densityMap = densityMap
	var maxX, maxY, number int
	for i := range densityMap {
		for j := range densityMap[i] {
			var c int
			if a.enemyStates[i][j] == gui.Hit || a.enemyStates[i][j] == gui.Miss || In(a.algorithm.shot, helpers.AlphabeticCoords(i, j)) {
				continue
			}
			if a.algorithm.options.Mixed {
				c = densityMap[i][j] * a.statistics[helpers.AlphabeticCoords(i, j)]
			} else {
				c = densityMap[i][j]
			}
			if c > number {
				number = c
				maxX = i
				maxY = j
			}
		}
	}
	if !In(a.algorithm.shot, helpers.AlphabeticCoords(maxX, maxY)) {
		return maxX, maxY
	}
	return a.getRandomCoord()
}

func (a *App) DoesShapeFit(x, y int, shape []point) bool {
	for _, c := range shape {
		dx := x + c.x
		dy := y + c.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 || a.enemyStates[dx][dy] == gui.Hit || a.enemyStates[dx][dy] == gui.Miss {
			return false
		}
		//if a.enemyStates[dx][dy] != gui.Empty {
		//	return false
		//}
	}
	return true
}

var shapes = [][]point{
	{
		{0, 1}, // #
		{1, 1}, // ##
		{1, 2}, //  #
	},
	{
		{0, 1},  //  #
		{-1, 1}, // ##
		{-1, 2}, // #
	},
	{
		{1, 0}, // ##
		{1, 1}, //  ##
		{2, 1}, //
	},
	{
		{1, 0},  //  ##
		{0, 1},  // ##
		{-1, 1}, //
	},
	{
		{1, 0}, // ##
		{0, 1}, // ##
		{1, 1}, //
	},
	{
		{1, 0}, // horizontal straight
		{2, 0},
		{3, 0},
	},
	{
		{0, 1}, // vertical straight
		{0, 2},
		{0, 3},
	},
	{
		{0, 1}, // #
		{0, 2}, // #
		{1, 2}, // ##
	},
	{
		{0, 1},  //  #
		{0, 2},  //  #
		{-1, 2}, // ##
	},
	//{
	//	{1, 0},  //  #   WE HAVE TO DO ANOTHER CHECK
	//	{1, -1}, //  #   BUT FROM DIFFERENT SIDE
	//	{1, -2}, // ##   THIS ONE IS FROM LEFT DOWN
	//},
	{
		{1, 0}, // ##
		{1, 1}, //  #
		{1, 2}, //  #
	},
	{
		{1, 0}, // ##
		{0, 1}, // #
		{0, 2}, // #
	},
	{
		{0, 1}, // #
		{1, 1}, // ###
		{2, 1},
	},
	{
		{1, 0}, // ###
		{2, 0}, // #
		{0, 1},
	},
	{
		{0, 1},  //   #
		{-1, 1}, // ###
		{-2, 1},
	},
	//{
	//	{1, 0}, //   #   SAME SITUATION HERE
	//	{2, 0}, // ###   LEFT DOWN
	//	{2, -1},
	//},
	{
		{1, 0}, // ###
		{2, 0}, //   #
		{2, 1},
	},
	{
		{1, 0}, // ###
		{2, 0}, //  #
		{1, 1},
	},
	{
		{0, 1}, // #
		{1, 1}, // ##
		{0, 2}, // #
	},
	{
		{0, 1},  //  #
		{-1, 1}, // ##
		{0, 2},  //  #
	},
	{
		{1, 1}, //  #
		{0, 1}, // ###
		{-1, 1},
	},
	// ALL FOURS
	{
		{1, 1}, // #
		{0, 1}, // ##
	},
	{
		{0, 1}, // ##
		{0, 1}, // #
	},
	{
		{1, 1}, // ##
		{1, 0}, //  #
	},
	{
		{1, 1},  //  #
		{-1, 1}, // ##
	},
	// ALL THREES
	{
		{0, 1},
	},
	{
		{1, 0},
	},
	// ALL TWOS
}
