package app

import gui "github.com/grupawp/warships-gui/v2"

type State int

type Algorithm struct {
	board      [10][10]gui.State
	statistics map[string]int
}

func newAlgorithm() *Algorithm {
	return &Algorithm{
		[10][10]gui.State{},
		map[string]int{},
	}
}

func (a *Algorithm) LoadStatistics() {}
func (a *Algorithm) ChangeState(x, y int, state gui.State) {
	if !(x < 0 || x >= 10 || y < 0 || y >= 10) {
		a.board[x][y] = state
	}
}

func (a *Algorithm) ClosestShip(x, y int) int {
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
		for _, point := range vec {
			dx := x + i*point.x
			dy := y + i*point.y
			if a.board[dx][dy] == gui.Hit || a.board[dx][dy] == "Sunk" {
				return i
			}
		}
	}
	return 0
}
