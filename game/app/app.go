package app

import "fmt"

const (
	NICK = "Nerien"
	DESC = "test"
)

type client interface {
	InitGame(coords []string, desc, nick, targetOpponent string, wpbot bool) error
	PrintToken()
	Board() ([]string, error)
	Status() (*StatusResponse, error)
	Shoot(coord string) (string, error)
	Resign() error
	GetOppDesc() (string, string, error)
	Refresh() error
	PlayerList() ([][]string, error)
	Stats() (map[string][]int, error)
	StatsPlayer(nick string) ([]int, error)
}

type App struct {
	client   client
	oppShots []string
	oppNick  string
	oppDesc  string
}

func New(c client) *App {
	return &App{
		c,
		[]string{},
		"",
		"",
	}
}

func (a *App) Run() error {

	stats, err := a.client.StatsPlayer("Nerien")
	if err != nil {
		return err
	}
	fmt.Println(stats)
	return nil

	//answer, err := a.getAnswer()
	//if err != nil {
	//	return err
	//}
	//if answer == "1" {
	//	err := a.client.InitGame(nil, DESC, NICK, "", true)
	//	if err != nil {
	//		return err
	//	}
	//} else if answer == "2" {
	//	playerlist, err := a.client.PlayerList()
	//	if err != nil {
	//		return err
	//	}
	//	fmt.Println(playerlist)
	//} else {
	//	fmt.Println("Please enter a number from the list!")
	//	return nil
	//}
	//
	//boardCoords, err := a.client.Board()
	//if err != nil {
	//	return err
	//}
	//status, err := a.WaitForStart()
	//if err != nil {
	//	return err
	//}
	//
	//a.opp_nick, a.opp_desc, err = a.client.GetOppDesc()
	//if err != nil {
	//	return err
	//}
	//
	//board := gui.New(
	//	gui.NewConfig(),
	//)
	//board.Import(boardCoords)
	//a.show(board, status)
	//
	//// MAIN GAME LOOP
	//for {
	//	if a.CheckIfWon() {
	//		return nil
	//	}
	//	err = a.Play(board, status)
	//	if err != nil {
	//		return err
	//	}
	//	err = a.OpponentShots(board)
	//	if err != nil {
	//		return err
	//	}
	//}
}
