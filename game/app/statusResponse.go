package app

type StatusResponse struct {
	Game_status      string
	Last_game_status string
	Nick             string
	Opp_desc         string
	Opponent         string
	Should_fire      bool
	Timer            int
}
