package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type logging struct {
	filename string
	*log.Logger
}

var logger *logging
var once sync.Once

func GetLoggerInstance() *logging {
	once.Do(func() {
		dt := time.Now()
		t := dt.Format("2006-01-02::15-04-05")
		logger = createLogger(fmt.Sprintf("%s.log", t))
	})
	return logger
}

func createLogger(fname string) *logging {
	absPath, err := filepath.Abs("logs")
	if err != nil {
		fmt.Println("Error reading given path:", err)
	}
	file, _ := os.OpenFile(absPath+"/"+fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &logging{
		filename: fname,
		Logger:   log.New(file, "Error: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
