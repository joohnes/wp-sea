package main

import (
	"fmt"
	"github.com/joohnes/wp-sea/game/app"
	"github.com/joohnes/wp-sea/game/client"
	"time"
)

const (
	serverAddr        = "https://go-pjatk-server.fly.dev"
	httpClientTimeout = 30 * time.Second
)

func main() {
	c := client.New(serverAddr, httpClientTimeout)
	application := app.New(c)

	err := application.Run()
	if err != nil {
		fmt.Print(err)
	}
}
