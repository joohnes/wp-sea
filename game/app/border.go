package app

import gui "github.com/grupawp/warships-gui/v2"

type point struct {
	x, y int
}

func (a *App) MarkBorders(x, y int, board *[10][10]gui.State, enemy bool) {
	var points []point
	a.searchShips(x, y, &points, board)
	for _, i := range points {
		a.drawBorder(i, board)
	}
	if enemy {
		a.enemyShips[len(points)] -= 1
	}
}

func (a *App) searchShips(x, y int, points *[]point, board *[10][10]gui.State) {
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
		//if a.enemyStates[x+v.x][y+v.y] == "Ship" || a.enemyStates[x+v.x][y+v.y] == "Hit" {
		//	connections = append(connections, point{dx, dy})
		//}
		if board[x+v.x][y+v.y] == "Ship" || board[x+v.x][y+v.y] == "Hit" {
			connections = append(connections, point{dx, dy})
		}
	}

	for _, c := range connections {
		a.searchShips(c.x, c.y, points, board)
	}
}

func (a *App) drawBorder(p point, board *[10][10]gui.State) {
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

		//prev := a.enemyStates[dx][dy]
		//if !(prev == "Ship" || prev == "Hit" || prev == "Miss") {
		//	a.enemyStates[dx][dy] = "Miss"
		//}
		prev := board[dx][dy]
		if !(prev == "Ship" || prev == "Hit" || prev == "Miss") {
			board[dx][dy] = "Miss"
		}
	}
}
