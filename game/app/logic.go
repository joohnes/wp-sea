package app

import (
	board "github.com/grupawp/warships-lightgui"
	"time"
)

const (
	waitDuration = 3
)

func (a *App) WaitForStart() (status *StatusResponse, err error) {
	for {
		status, err := a.client.Status()
		if err != nil {
			return nil, err
		}
		if status.Game_status == "game_in_progress" {
			return status, nil
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
		if status.Should_fire == true {
			break
		}
		time.Sleep(waitDuration * time.Second)
	}
	return nil
}

func (a *App) Shoot(bd *board.Board) (string, error) {
	err := a.WaitForTurn()
	if err != nil {
		return "", err
	}
	coord, err := getCoord()
	if err != nil {
		return "", err
	}

	result, err := a.client.Shoot(coord)
	if err != nil {
		return "", err
	}
	switch result {
	case "miss":
		bd.Set(board.Right, coord, board.Miss)
	case "hit":
		bd.Set(board.Right, coord, board.Hit)
	case "sunk":
		bd.Set(board.Right, coord, board.Hit)
		bd.CreateBorder(board.Right, coord)
	}

	return result, nil
}
