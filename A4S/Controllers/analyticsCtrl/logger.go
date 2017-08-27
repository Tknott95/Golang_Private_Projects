package analytics

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ErrInvalidID    = errors.New("Invalid ID")
	ErrInvalidEmail = errors.New("Invalid email")
)

func CreateLogger(filename string, msgToLog string) *log.Logger {
	file, err := os.OpenFile("./Analytics/LogFiles/"+"LOG_"+filename+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) /*  Better number for prod app ( this reads and writes) */
	if err != nil {
		panic(err)
	}
	if len(msgToLog) < 1 {
		msgToLog = "No Msg"
	}
	logger := log.New(file, " - "+msgToLog+" - ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

func Time(logger *log.Logger, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		elapsed := time.Since(start)
		logger.Println(elapsed)
	})
}
