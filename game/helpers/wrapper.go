package helpers

import (
	"errors"
	"fmt"
	"github.com/joohnes/wp-sea/game/logger"
	"time"
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
			fmt.Printf("#%d - Server error occurred. Please wait\n", x+1)
		}
		log.Printf("#%d - Server error occurred. Please wait\n", x+1)
		time.Sleep(1 * time.Second)
	}
	return errors.New("max connection tries reached")
}
