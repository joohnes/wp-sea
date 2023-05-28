package app

type StatusResponse struct {
	GameStatus     string   `json:"game_status,omitempty"`
	LastGameStatus string   `json:"last_game_status,omitempty"`
	Nick           string   `json:"nick,omitempty"`
	OppDesc        string   `json:"opp_desc,omitempty"`
	OppShots       []string `json:"opp_shots,omitempty"`
	Opponent       string   `json:"opponent,omitempty"`
	ShouldFire     bool     `json:"should_fire,omitempty"`
	Timer          int      `json:"timer,omitempty"`
}
