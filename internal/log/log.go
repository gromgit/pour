package log

import (
	"log"
	"os"
)

// Private logger instance
var logger *log.Logger

func File(f *os.File) {
	logger = log.New(f, "pour", log.LstdFlags)
}

func Log(v ...interface{}) {
	if logger != nil {
		logger.Println(v...)
	}
}

func Logf(s string, v ...interface{}) {
	if logger != nil {
		logger.Printf(s, v...)
	}
}
