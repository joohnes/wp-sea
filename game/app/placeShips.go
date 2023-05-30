package app

import (
	"context"
	"errors"
	"github.com/joohnes/wp-sea/game/helpers"
)

func (a *App) CheckIfChangedMap() bool {
	for _, x := range a.playerStates {
		for _, y := range x {
			if y != "Ship" {
				return true
			}
		}
	}
	return false
}

var (
	ErrInvalidShape = errors.New("invalid shape")
)

func (a *App) PlaceShips(ctx context.Context, cancel context.CancelFunc, shipchannel chan string, errorchan chan error) {
	for {
		select {
		case coord := <-shipchannel:
			coords, err := helpers.NumericCords(coord)
			if err != nil {
				errorchan <- err
				break
			}
			err = a.ValidateShipPlacement(coords, cancel)
			if err != nil {
				errorchan <- err
				break
			}

		case <-ctx.Done():
			return
		}
	}
}

func (a *App) ValidateShipPlacement(coords map[string]uint8, cancel context.CancelFunc) error {
	if a.playerStates[coords["x"]][coords["y"]] == "" {
		err := a.CheckShipLength(coords)
		if err != nil {
			return err
		}
		err = a.CheckCorners(coords)
		if err != nil {
			return err
		}
		err = a.CheckCornerNumber(coords)
		if err != nil {
			return err
		}
		err = a.CheckFigures()
		if err != nil {
			return err
		}
		a.playerStates[coords["x"]][coords["y"]] = "Ship"
	} else if a.playerStates[coords["x"]][coords["y"]] == "Ship" {
		a.playerStates[coords["x"]][coords["y"]] = ""
	}
	return nil
}

func (a *App) CheckCorners(coords map[string]uint8) error {
	points := []point{
		{-1, 1},
		{1, 1},
		{1, -1},
		{-1, -1},
	}

	for _, v := range points {
		dx := int(coords["x"]) + v.x
		dy := int(coords["y"]) + v.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
			continue
		}
		if a.playerStates[dx][dy] == "Ship" {
			vec := []point{
				{v.x, 0},
				{0, v.y},
			}
			for _, s := range vec {
				dx1 := int(coords["x"]) + s.x
				dy1 := int(coords["y"]) + s.y
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

func (a *App) CheckCornerNumber(coords map[string]uint8) error {
	points := []point{
		{-1, 1},
		{1, 1},
		{1, -1},
		{-1, -1},
	}
	corners := 0

	for _, v := range points {
		dx := int(coords["x"]) + v.x
		dy := int(coords["y"]) + v.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
			continue
		}
		if a.playerStates[dx][dy] == "Ship" {
			corners++
			if corners > 1 {
				return errors.New("you can't put ships diagonally")
			}
		}
	}
	return nil
}

func (a *App) CheckShipLength(coords map[string]uint8) error {
	var points []point
	a.countShips(int(coords["x"]), int(coords["y"]), &points)
	if len(points) > 4 {
		return errors.New("too long")
	}
	return nil
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

func (a *App) CheckFigures() error {
	// Wrong figures
	vec := [][]point{
		{
			{1, 0},
			{1, 1},
			{2, 1},
		},
		{
			{1, 0},
			{0, 1},
			{-1, 1},
		},
		{
			{0, 1},
			{-1, 1},
			{-1, 2},
		},
		{
			{0, 1},
			{1, 1},
			{1, 2},
		},
	}

	for i, c := range a.playerStates {
		for j, d := range c {
			if d == "Ship" {
				for _, shape := range vec {
					counter := 0
					x, y := 0, 0
					for _, coord := range shape {
						dx := coord.x + i
						dy := coord.y + j
						if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
							continue
						}
						if a.playerStates[dx][dx] != "Ship" {
							break
						}
						counter++
						x = dx
						y = dy
					}
					if counter > 2 {
						a.playerStates[x][y] = ""
						return ErrInvalidShape
					}
				}
			}
		}
	}
	return nil
}
