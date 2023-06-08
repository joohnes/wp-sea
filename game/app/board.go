package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	gui "github.com/grupawp/warships-gui/v2"
)

var (
	playerCfg = gui.TextConfig{FgColor: gui.Black, BgColor: gui.Green}
	oppCfg    = gui.TextConfig{FgColor: gui.Black, BgColor: gui.Red}
	errCfg    = gui.TextConfig{FgColor: gui.Black, BgColor: gui.Red}
)

var (
	Left  = 1
	Right = 50
)

func (a *App) ShowBoard(ctx context.Context, coordchan chan<- string, textchan, predchan <-chan string, errorchan <-chan error, timeLeftchan <-chan int) {
	// SETUP
	ui := gui.NewGUI(false)
	txt := gui.NewText(Left, 1, "Press Ctrl+C to exit", nil)
	myBoard := gui.NewBoard(Left, 6, nil)
	enemyBoard := gui.NewBoard(Right, 6, nil)

	//TEXTS
	timer := gui.NewText(Left, 2, "Timer: ", nil)
	turnInfo := gui.NewText(Left, 4, "Turn: 1", nil)
	playerNick := gui.NewText(Left, 30, fmt.Sprintf("Player: %s", a.nick), &playerCfg)
	playerDesc := gui.NewText(Left, 31, fmt.Sprintf("Player's desc: %s", a.desc), &playerCfg)
	oppNick := gui.NewText(Left, 33, fmt.Sprintf("Opp: %s", a.oppNick), &oppCfg)
	oppDesc := gui.NewText(Left, 34, fmt.Sprintf("Opp's desc: %s", a.oppDesc), &oppCfg)

	turnText := gui.NewText(Right, 4, "", nil)
	chanText := gui.NewText(Left, 36, "", nil)
	predText := gui.NewText(Right, 36, "", nil)
	errorText := gui.NewText(Left, 37, "", &errCfg)

	shotsCounttxt := gui.NewText(Right, 1, "Shots: 0", nil)
	shotsHittxt := gui.NewText(Right, 2, "Hit: 0", nil)
	accuracytxt := gui.NewText(Right, 3, "Accuracy: %", nil)

	playerBoardIndicator := gui.NewText(Left, 28, "Your Board", &playerCfg)
	enemyBoardIndicator := gui.NewText(Right, 28, "Opponent's Board", &oppCfg)

	legendEmpty := gui.NewText(28, 1, "~ -> Empty space", nil)
	legendShip := gui.NewText(28, 2, "S -> Ship", nil)
	legendHit := gui.NewText(28, 3, "H -> Hit", nil)
	legendMiss := gui.NewText(28, 4, "M -> Miss", nil)
	legendEmpty.SetBgColor(gui.Blue)
	legendShip.SetBgColor(gui.Green)
	legendHit.SetBgColor(gui.Red)
	legendMiss.SetBgColor(gui.Grey)

	shipsleft4 := gui.NewText(
		Right+20,
		1,
		fmt.Sprintf("4 mast: %d", a.enemyShips[4]),
		nil)
	shipsleft3 := gui.NewText(
		Right+20,
		2,
		fmt.Sprintf("3 mast: %d", a.enemyShips[3]),
		nil)
	shipsleft2 := gui.NewText(
		Right+20,
		3,
		fmt.Sprintf("2 mast: %d", a.enemyShips[2]),
		nil)
	shipsleft1 := gui.NewText(
		Right+20,
		4,
		fmt.Sprintf("1 mast: %d", a.enemyShips[1]),
		nil)

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
		turnInfo,
		turnText,
		chanText,
		predText,
		errorText,
		playerNick,
		playerDesc,
		oppNick,
		oppDesc,
		shotsHittxt,
		shotsCounttxt,
		accuracytxt,
		shipsleft4,
		shipsleft3,
		shipsleft2,
		shipsleft1,
		playerBoardIndicator,
		enemyBoardIndicator,
		legendEmpty,
		legendHit,
		legendMiss,
		legendShip,
	)

	if a.Requirements() {
		a.myStates = a.playerStates
	}
	myBoard.SetStates(a.myStates)
	enemyBoard.SetStates(a.enemyStates)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				char := enemyBoard.Listen(context.TODO())
				coordchan <- char
				ui.Log("Coordinate: %s", char)
			}
		}
	}()

	go func() {
		t := time.NewTicker(time.Second)
		counter := 0
		for {
			select {
			case <-ctx.Done():
				return
			case text := <-textchan:
				if text == "You have won the game!" {
					chanText.SetBgColor(gui.Green)
					chanText.SetText(text)
				} else if text == "You have lost the game!" {
					chanText.SetBgColor(gui.Red)
					chanText.SetText(text)
				} else if text == "color reset" {
					chanText.SetBgColor(gui.White)
					chanText.SetText("")
					predText.SetText("")
				} else {
					chanText.SetText(text)
				}

			case pred := <-predchan:
				predText.SetText(fmt.Sprintf("Algorithm suggests: %s", pred))

			case err := <-errorchan:
				errorText.SetText(err.Error())
				counter = 3
			case <-t.C:
				if counter > 1 {
					counter--
				} else {
					errorText.SetText("")
				}

			case timeLeft := <-timeLeftchan:
				timer.SetText(fmt.Sprintf("Time left: %v", timeLeft))
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				switch a.gameState {
				case StatePlayerTurn:
					turnText.SetText("Your Turn!")
					turnText.SetBgColor(gui.Green)
				case StateOppTurn:
					turnText.SetText("Enemy's turn!")
					turnText.SetBgColor(gui.Red)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(50 * time.Millisecond)
				turnInfo.SetText("Turn: " + strconv.Itoa(a.turn))
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

				if a.enemyShips[4] == 0 {
					shipsleft4.SetBgColor(gui.Green)
				}
				if a.enemyShips[3] == 0 {
					shipsleft3.SetBgColor(gui.Green)
				}
				if a.enemyShips[2] == 0 {
					shipsleft2.SetBgColor(gui.Green)
				}
				if a.enemyShips[1] == 0 {
					shipsleft1.SetBgColor(gui.Green)
				}
				shipsleft4.SetText(fmt.Sprintf("4 mast: %d", a.enemyShips[4]))
				shipsleft3.SetText(fmt.Sprintf("3 mast: %d", a.enemyShips[3]))
				shipsleft2.SetText(fmt.Sprintf("2 mast: %d", a.enemyShips[2]))
				shipsleft1.SetText(fmt.Sprintf("1 mast: %d", a.enemyShips[1]))
			}
		}
	}()

	ui.Start(ctx, nil)
}

func (a *App) SetUpShips(ctx context.Context, shipchannel chan string, errorchan chan error) {
	ui := gui.NewGUI(false)
	txt := gui.NewText(Left, 1, "Press Ctrl+C to exit", nil)
	myBoard := gui.NewBoard(Left, 4, nil)
	infoText := gui.NewText(Left, 27, "Press any field to put a ship there", nil)
	shipsInfo := gui.NewText(Right, 4, "Ships you must place to start a game:", nil)
	shipsText4 := gui.NewText(Right, 6, fmt.Sprintf("4 mast: %d left", a.placeShips[4]), &playerCfg)
	shipsText3 := gui.NewText(Right, 8, fmt.Sprintf("3 mast: %d left", a.placeShips[3]), &playerCfg)
	shipsText2 := gui.NewText(Right, 10, fmt.Sprintf("2 mast: %d left", a.placeShips[2]), &playerCfg)
	shipsText1 := gui.NewText(Right, 12, fmt.Sprintf("1 mast: %d left", a.placeShips[1]), &playerCfg)
	shipsReqText := gui.NewText(Right, 16, "", nil)
	errorText := gui.NewText(Left, 31, "", &errCfg)
	myBoard.SetStates(a.playerStates)

	DrawList := func(a ...gui.Drawable) {
		for _, d := range a {
			ui.Draw(d)
		}
	}
	DrawList(
		txt,
		myBoard,
		errorText,
		infoText,
		shipsInfo,
		shipsText4,
		shipsText3,
		shipsText2,
		shipsText1,
		shipsReqText,
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				char := myBoard.Listen(context.TODO())
				shipchannel <- char
			}
		}
	}()
	go func() {
		t := time.NewTicker(time.Second)
		counter := 0
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errorchan:
				errorText.SetText(err.Error())
				counter = 3
			case <-t.C:
				if counter > 1 {
					counter--
				} else {
					errorText.SetText("")
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(50 * time.Millisecond)
				myBoard.SetStates(a.playerStates)
				shipsText4.SetText(fmt.Sprintf("4 mast: %d left", a.placeShips[4]))
				shipsText3.SetText(fmt.Sprintf("3 mast: %d left", a.placeShips[3]))
				shipsText2.SetText(fmt.Sprintf("2 mast: %d left", a.placeShips[2]))
				shipsText1.SetText(fmt.Sprintf("1 mast: %d left", a.placeShips[1]))

				if !a.Requirements() {
					shipsReqText.SetBgColor(gui.Red)
					shipsReqText.SetText("You need to place all the ships for it to be sent to server")
				} else {
					shipsReqText.SetBgColor(gui.Green)
					shipsReqText.SetText("Your board is ready to start the game!")
				}
			}
		}
	}()

	ui.Start(ctx, nil)
}
