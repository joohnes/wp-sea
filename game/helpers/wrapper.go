package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/joohnes/wp-sea/game/logger"
)

const maxConnectionTries int = 10

func ServerErrorWrapper(showErrors bool, f func() error) error {
	log := logger.GetLoggerInstance()
	for x := 0; x < maxConnectionTries; x++ {
		err := f()
		if err == nil {
			return nil
		}
		if showErrors {
			fmt.Printf("#%d - %s\n", x+1, err.Error())
		}
		log.Printf("#%d - %s\n", x+1, err.Error())
		time.Sleep(1 * time.Second)
	}
	return errors.New("max connection tries reached")
}
