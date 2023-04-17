package app

type StatusResponse struct {
	Game_status      string   `json: "game_status,omitempty"`
	Last_game_status string   `json: "last_game_status,omitempty"`
	Nick             string   `json: "nick,omitempty"`
	Opp_desc         string   `json: "opp_desc,omitempty"`
	Opp_shots        []string `json: "opp_shots,omitempty"`
	Opponent         string   `json: "opponent,omitempty"`
	Should_fire      bool     `json: "should_fire,omitempty"`
	Timer            int      `json: "timer,omitempty"`
}
