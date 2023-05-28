package app

type Gamestate string

const (
	StateStart      Gamestate = "Start"
	StateWaiting              = "Waiting"
	StatePlayerTurn           = "PlayerTurn"
	StateOppTurn              = "OppTurn"
	StateEnded                = "Ended"
)
