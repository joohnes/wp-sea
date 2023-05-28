package app

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"time"
)

func (a *App) ShowBoard(ctx context.Context, coordchan chan<- string, textchan <-chan string, errorchan <-chan error, timeLeftchan <-chan int) {

	// SETUP
	ui := gui.NewGUI(true)
	Left := 1
	Right := 50

	txt := gui.NewText(Left, 1, "Press Ctrl+C to exit", nil)
	myBoard := gui.NewBoard(Left, 6, nil)
	enemyBoard := gui.NewBoard(Right, 6, nil)
	shipsleft := gui.NewText(
		Right,
		29,
		fmt.Sprintf("4 mast: %d | 3 mast: %d | 2 mast: %d | 1 mast: %d", a.enemyShips[4], a.enemyShips[3], a.enemyShips[2], a.enemyShips[1]),
		nil)

	//TEXTS
	timer := gui.NewText(Left, 2, "Timer: ", nil)
	playerNick := gui.NewText(Left, 29, fmt.Sprintf("Player: %s", a.nick), nil)
	playerDesc := gui.NewText(Left, 30, fmt.Sprintf("Player's desc: %s", a.desc), nil)
	oppNick := gui.NewText(Left, 31, fmt.Sprintf("Opp: %s", a.oppNick), nil)
	oppDesc := gui.NewText(Left, 32, fmt.Sprintf("Opp's desc: %s", a.oppDesc), nil)

	turnText := gui.NewText(Right, 4, "", nil)
	chanText := gui.NewText(Left, 34, "", nil)
	errorText := gui.NewText(Left, 35, "", nil)
	errorText.SetBgColor(gui.Red)

	shotsCounttxt := gui.NewText(Right, 1, "Shots: 0", nil)
	shotsHittxt := gui.NewText(Right, 2, "Hit: 0", nil)
	accuracytxt := gui.NewText(Right, 3, "Accuracy: %", nil)

	DrawList := func(a ...gui.Drawable) {
		for _, d := range a {
			ui.Draw(d)
		}
	}
	DrawList(
		txt,
		myBoard,
		enemyBoard,
		timer,
		turnText,
		chanText,
		errorText,
		playerNick,
		playerDesc,
		oppNick,
		oppDesc,
		shotsHittxt,
		shotsCounttxt,
		accuracytxt,
		shipsleft,
	)

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
				if text == "You have won the game!" {
					chanText.SetBgColor(gui.Green)
				} else if text == "You have lost the game!" {
					chanText.SetBgColor(gui.Red)
				}

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
				turnText.SetBgColor(gui.Green)
			case StateOppTurn:
				turnText.SetText("Enemy's turn!")
				turnText.SetBgColor(gui.Red)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			myBoard.SetStates(a.myStates)
			enemyBoard.SetStates(a.enemyStates)
			shotsCounttxt.SetText(fmt.Sprintf("Shots: %d", a.shotsCount))
			shotsHittxt.SetText(fmt.Sprintf("Hits: %d", a.shotsHit))
			var accuracy float64
			if a.shotsCount != 0 {
				accuracy = float64(a.shotsHit) / float64(a.shotsCount) * 100
			}
			accuracytxt.SetText(fmt.Sprintf("Accuracy: %.2f%%", accuracy))
			if accuracy > 60 {
				accuracytxt.SetBgColor(gui.Green)
			} else {
				accuracytxt.SetBgColor(gui.White)
			}

			shipsleft.SetText(fmt.Sprintf("4 mast: %d | 3 mast: %d | 2 mast: %d | 1 mast: %d", a.enemyShips[4], a.enemyShips[3], a.enemyShips[2], a.enemyShips[1]))
		}
	}()

	ui.Start(ctx, nil)

}

func (a *App) SetUpShips(ctx context.Context, shipchannel chan string, errorchan chan error) {
	ui := gui.NewGUI(true)
	Left := 1
	//Right := 50
	txt := gui.NewText(Left, 1, "Press Ctrl+C to exit", nil)
	ui.Draw(txt)
	myBoard := gui.NewBoard(Left, 4, nil)
	ui.Draw(myBoard)
	errorText := gui.NewText(Left, 32, "error", nil)
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
