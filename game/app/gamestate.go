package app

type Gamestate string

const (
	StateStart      Gamestate = "Start"
	StatePlayerTurn           = "PlayerTurn"
	StateOppTurn              = "OppTurn"
	StateEnded                = "Ended"
)
