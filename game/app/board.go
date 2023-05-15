package app

import (
	"context"
	"fmt"

	gui "github.com/grupawp/warships-gui/v2"
)

func (a *App) ShowBoard(coordchan chan string) {
	ui := gui.NewGUI(true)
	txt := gui.NewText(1, 1, "Press Ctrl+C to exit", nil)
	ui.Draw(txt)
	my_board := gui.NewBoard(1, 4, nil)
	ui.Draw(my_board)
	enemy_board := gui.NewBoard(50, 4, nil)
	ui.Draw(enemy_board)

	//TEXTS
	timer := gui.NewText(1, 2, "Timer: ", nil)
	ui.Draw(timer)
	playerNick := gui.NewText(1, 26, fmt.Sprintf("Player: %s", a.nick), nil)
	playerDesc := gui.NewText(1, 27, fmt.Sprintf("Player's desc: %s", a.desc), nil)
	oppNick := gui.NewText(1, 28, fmt.Sprintf("Opp: %s", a.oppNick), nil)
	oppDesc := gui.NewText(1, 29, fmt.Sprintf("Opp's desc: %s", a.oppDesc), nil)
	ui.Draw(playerNick)
	ui.Draw(playerDesc)
	ui.Draw(oppNick)
	ui.Draw(oppDesc)

	turnText := gui.NewText(1, 30, "", nil)
	ui.Draw(turnText)

	for i := range a.my_states {
		a.my_states[i] = [10]gui.State{}
		a.enemy_states[i] = [10]gui.State{}
	}
	my_board.SetStates(a.my_states)
	enemy_board.SetStates(a.enemy_states)
	go func() {
		for {
			char := enemy_board.Listen(context.TODO())
			txt.SetText(fmt.Sprintf("Coordinate: %s", char))
			coordchan <- char
			ui.Log("Coordinate: %s", char) // logs are displayed after the game exits
		}
	}()

	// go func(cause error) {
	// 	for {
	// 		select {

	// 		}
	// 	}
	// }
	ui.Start(nil)
}
