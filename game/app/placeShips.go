package app

import (
	"context"
	"errors"

	"github.com/joohnes/wp-sea/game/helpers"
)

var (
	ErrInvalidShape = errors.New("invalid shape")
)

func (a *App) CheckIfChangedMap() bool {
	for _, x := range a.playerStates {
		for _, y := range x {
			if y == "Ship" {
				return true
			}
		}
	}
	return false
}

func (a *App) Requirements() bool {
	for _, x := range a.placeShips {
		if x != 0 {
			return false
		}
	}
	return true
}

func (a *App) PlaceShips(ctx context.Context, shipchannel chan string, errorchan chan error) {
	for {
		select {
		case coord := <-shipchannel:
			coords, err := helpers.NumericCords(coord)
			if err != nil {
				errorchan <- err
				break
			}
			err = a.ValidateShipPlacement(int(coords["x"]), int(coords["y"]))
			if err != nil {
				errorchan <- err
			}

		case <-ctx.Done():
			return
		}
	}
}

func (a *App) ValidateShipPlacement(x, y int) error {
	if a.playerStates[x][y] == "" {
		points := a.CheckShipLength(x, y)
		if a.placeShips[len(points)] == 0 {
			return errors.New("you can't place anymore ships of that type")
		}

		//if len(points) == 4 {
		//	err := a.CheckForWrongFigures(points)
		//	if err != nil {
		//		return err
		//	}
		//} else if len(points) > 4 {
		//	return errors.New("too long")
		//}
		/*
			It seems that i imagined a few of the rules and added extra checks
			I'll leave it here for the sake of my lost time
		*/

		if len(points) > 4 {
			return errors.New("too long")
		}
		err := a.CheckCorners(x, y)
		if err != nil {
			return err
		}
		a.playerStates[x][y] = "Ship"
		a.CheckAllShipsLength()
	} else if a.playerStates[x][y] == "Ship" {
		a.playerStates[x][y] = ""
		points := a.CheckShipLength(x, y)
		if len(points) > 2 {
			a.CheckForLoners(x, y)
		}
		a.CheckLeftShips(points)
		a.CheckAllShipsLength()
	}
	return nil
}

func (a *App) CheckCorners(x, y int) error {
	points := []point{
		{-1, 1},
		{1, 1},
		{1, -1},
		{-1, -1},
	}

	for _, v := range points {
		dx := x + v.x
		dy := y + v.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
			continue
		}
		if a.playerStates[dx][dy] == "Ship" {
			vec := []point{
				{v.x, 0},
				{0, v.y},
			}
			for _, s := range vec {
				dx1 := x + s.x
				dy1 := y + s.y
				if dx1 < 0 || dx1 >= 10 || dy1 < 0 || dy1 >= 10 {
					continue
				}
				if a.playerStates[dx1][dy1] == "Ship" {
					return nil
				}
			}
			return errors.New("you can't put ships diagonally")
		}
	}
	return nil
}

func (a *App) CheckShipLength(x, y int) []point {
	var points []point
	a.countShips(x, y, &points)
	return points
}

func (a *App) countShips(x, y int, points *[]point) {
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
		if a.playerStates[dx][dy] == "Ship" {
			connections = append(connections, point{dx, dy})
		}

	}
	for _, c := range connections {
		a.countShips(c.x, c.y, points)
	}
}

func IsIn(arr []point, x, y int) bool {
	for _, pair := range arr {
		if pair.x == x && pair.y == y {
			return true
		}
	}
	return false
}

func (a *App) CheckAllShipsLength() {
	var checked []point
	basemap := map[int]int{4: 1, 3: 2, 2: 3, 1: 4}
	for i := range a.playerStates {
		for j, state := range a.playerStates[i] {
			if state == "Ship" {
				if IsIn(checked, i, j) {
					continue
				}
				points := a.CheckShipLength(i, j)
				basemap[len(points)]--
				checked = append(checked, points...)

			}
		}
	}
	a.placeShips = basemap
}

// CheckForWrongFigures will check if ship (of length 4) is in the correct shape.
// If it is not, return ErrInvalidShape
// Well it happens that function is redundant
func (a *App) CheckForWrongFigures(points []point) error {
	correctShapes := [][]point{
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
		{
			{1, 0},  //  #   WE HAVE TO DO ANOTHER CHECK
			{1, -1}, //  #   BUT FROM DIFFERENT SIDE
			{1, -2}, // ##   THIS ONE IS FROM LEFT DOWN
		},
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
		{
			{1, 0}, //   #   SAME SITUATION HERE
			{2, 0}, // ###   LEFT DOWN
			{2, -1},
		},
		{
			{1, 0}, // ###
			{2, 0}, //   #
			{2, 1},
		},
	}

	minX := 10
	minY := 10
	for _, p := range points {
		if p.x <= minX && p.y <= minY {
			minX = p.x
			minY = p.y
		}
	}
	for _, shape := range correctShapes {
		check := true
		for _, block := range shape {
			dx := minX + block.x
			dy := minY + block.y
			if !IsIn(points, dx, dy) {
				check = false
			}
		}
		if check {
			return nil
		}
	}
	return ErrInvalidShape
}

// CheckForLoners checks the table for a "loner", which is ship on the corner
// of another ship that was left, when player deleted the part of a ship that was connecting them
func (a *App) CheckForLoners(x, y int) {
	vec := []point{
		{-1, 0}, // left
		{0, -1}, // up
		{1, 0},  // right
		{0, 1},  // down
	}

	for _, v := range vec {
		dx := x + v.x
		dy := y + v.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
			continue
		}
		if a.playerStates[dx][dy] == "Ship" {
			points := a.CheckShipLength(dx, dy)
			if len(points) == 1 {
				a.playerStates[dx][dy] = ""
				return
			}
		}
	}
}

// CheckLeftShips checks if the length of a ship left after the deletion of the part of a ship
// can be placed on the board (avoid having -1 as a number of ships left to place in a.placeShips)
func (a *App) CheckLeftShips(points []point) {
	var counter int
	for _, v := range points {
		if a.playerStates[v.x][v.y] == "Ship" {
			counter++
		}
	}
	if a.placeShips[counter] < 1 {
		for _, v := range points {
			a.playerStates[v.x][v.y] = ""
		}
	}
}
