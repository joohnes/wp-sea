package app

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"time"
)

func (a *App) ShowBoard(ctx context.Context, coordchan chan<- string, textchan <-chan string, errorchan <-chan error, timeLeftchan <-chan int) {
	ui := gui.NewGUI(true)
	txt := gui.NewText(1, 1, "Press Ctrl+C to exit", nil)
	ui.Draw(txt)
	myBoard := gui.NewBoard(1, 4, nil)
	ui.Draw(myBoard)
	enemyBoard := gui.NewBoard(50, 4, nil)
	ui.Draw(enemyBoard)

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

	turnText := gui.NewText(50, 1, "", nil)
	ui.Draw(turnText)
	chanText := gui.NewText(1, 31, "test", nil)
	ui.Draw(chanText)
	errorText := gui.NewText(1, 32, "error", nil)
	ui.Draw(errorText)

	myBoard.SetStates(a.myStates)
	enemyBoard.SetStates(a.enemyStates)
	go func() {
		for {
			char := enemyBoard.Listen(context.TODO())
			txt.SetText(fmt.Sprintf("Coordinate: %s", char))
			coordchan <- char
			ui.Log("Coordinate: %s", char)
		}
	}()

	go func() {
		for {
			select {
			case text := <-textchan:
				chanText.SetText(text)

			case err := <-errorchan:
				errorText.SetText(err.Error())
			case timeLeft := <-timeLeftchan:
				timer.SetText(fmt.Sprintf("Time left: %v", timeLeft))
			}
		}
	}()

	go func() {
		for {
			switch a.gameState {
			case StatePlayerTurn:
				turnText.SetText("Your Turn!")
			case StateOppTurn:
				turnText.SetText("Enemy's turn!")
			}
		}
	}()

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			myBoard.SetStates(a.myStates)
			enemyBoard.SetStates(a.enemyStates)
		}
	}()

	ui.Start(ctx, nil)

}

func (a *App) SetUpShips(ctx context.Context, shipchannel chan string, errorchan chan error) {
	ui := gui.NewGUI(true)
	txt := gui.NewText(1, 1, "Press Ctrl+C to exit", nil)
	ui.Draw(txt)
	myBoard := gui.NewBoard(1, 4, nil)
	ui.Draw(myBoard)
	errorText := gui.NewText(1, 32, "error", nil)
	ui.Draw(errorText)
	myBoard.SetStates(a.myStates)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				char := myBoard.Listen(context.TODO())
				txt.SetText(fmt.Sprintf("Coordinate: %s", char))
				shipchannel <- char
				ui.Log("Coordinate: %s", char)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errorchan:
				errorText.SetText(err.Error())
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(200 * time.Millisecond)
				myBoard.SetStates(a.myStates)
			}
		}
	}()

	ui.Start(ctx, nil)
}
