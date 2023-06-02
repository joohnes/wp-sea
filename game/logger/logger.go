package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logging struct {
	filename string
	*log.Logger
}

var logger *Logging
var once sync.Once

func GetLoggerInstance() *Logging {
	once.Do(func() {
		dt := time.Now()
		t := dt.Format("2006-01-02::15-04-05")
		logger = createLogger(fmt.Sprintf("%s.log", t))
	})
	return logger
}

func createLogger(fname string) *Logging {
	absPath, err := filepath.Abs("logs")
	if err != nil {
		fmt.Println("Error reading given path:", err)
	}
	if _, err := os.Stat("logs"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("logs", os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	}
	file, _ := os.OpenFile(absPath+"/"+fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &Logging{
		filename: fname,
		Logger:   log.New(file, "Error: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
