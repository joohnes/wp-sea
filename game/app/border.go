package app

import (
	"context"
)

type point struct {
	x, y int
}

func (a *App) MarkBorders(ctx context.Context, coordmap map[string]uint8, errorchan chan error) {
	//var c context.Context
	//if ctx.Value(fmt.Sprintf("%v-%v", coordmap["x"], coordmap["y"])) == nil {
	//	c = context.WithValue(ctx, fmt.Sprintf("%v-%v", coordmap["x"], coordmap["y"]), []uint8{coordmap["x"], coordmap["y"]})
	//	errorchan <- errors.New("dziala")
	//} else {
	//	return
	//}

	points := []point{}
	a.searchShips(int(coordmap["x"]), int(coordmap["y"]), &points)

	for _, i := range points {
		a.drawBorder(i)
	}

}

func (a *App) searchShips(x, y int, points *[]point) {
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
	connections := []point{}

	for _, v := range vec {
		dx := x + v.x
		dy := y + v.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
			continue
		}
		if a.enemyStates[x+v.x][y+v.y] == "Ship" || a.enemyStates[x+v.x][y+v.y] == "Hit" {
			connections = append(connections, point{dx, dy})
		}
	}

	// Run the method recursively on each linked element
	for _, c := range connections {
		a.searchShips(c.x, c.y, points)
	}
}

func (a *App) drawBorder(p point) {
	//
	//    XXX
	//    XOX
	//    XXX
	//
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

	for _, v := range vec {
		dx := p.x + v.x
		dy := p.y + v.y
		if dx < 0 || dx >= 10 || dy < 0 || dy >= 10 {
			continue

		}

		prev := a.enemyStates[dx][dy]
		if !(prev == "Ship" || prev == "Hit" || prev == "Miss") { // don't overwrite already marked
			a.enemyStates[dx][dy] = "Miss"
		}
	}
}
