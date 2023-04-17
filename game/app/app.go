package app

type client interface {
	InitGame(coords []string, desc, nick, target_opponent string, wpbot bool) error
	PrintToken()
	Board() ([]string, error)
	// Status() error
}

type app struct {
	c client
}

func New(c client) *app {
	return &app{
		c,
	}
}

func (a *app) Run() {
	a.c.InitGame(nil, "", "", "", true)
	a.c.PrintToken()
	a.c.Board()
}
