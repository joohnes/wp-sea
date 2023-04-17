package app

type StatusResponse struct {
	game_status      string
	last_game_status string
	nick             string
	opp_desc         string
	opponent         string
	should_fire      bool
	timer            int
}
