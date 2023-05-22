package app

import (
	"errors"
	"fmt"
	"time"
)

const maxConnectionTries int = 10

func ServerErrorWrapper(f func() error) error {
	log := GetLoggerInstance()
	for x := 0; x < maxConnectionTries; x++ {
		err := f()
		if err == nil {
			return nil
		}
		if err != nil && showErrors {
			fmt.Printf("#%d - Server error occured. Please wait\n", x+1)
		}
		log.Printf("#%d - Server error occured. Please wait\n", x+1)
		time.Sleep(1 * time.Second)
	}
	return errors.New("max connection tries reached")
}
