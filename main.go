package main

import (
	"game/game/app"
	"game/game/client"
	"time"
)

const (
	serverAddr        = "https://go-pjatk-server.fly.dev"
	httpClientTimeout = 30 * time.Second
)

func main() {
	c := client.New(serverAddr, httpClientTimeout)
	app := app.New(c)

	app.Run()
}
