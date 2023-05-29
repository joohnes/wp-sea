package app

import (
	"context"
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
			}

		case <-ctx.Done():
			return
		}
	}
}

func (a *App) ValidateShipPlacement(coords map[string]uint8, cancel context.CancelFunc) error {

	if a.playerStates[coords["x"]][coords["y"]] == "Ship" {
		a.playerStates[coords["x"]][coords["y"]] = ""
	} else if a.playerStates[coords["x"]][coords["y"]] == "" {

		a.playerStates[coords["x"]][coords["y"]] = "Ship"
	}

	return nil
}
