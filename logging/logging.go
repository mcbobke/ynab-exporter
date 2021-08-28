package logging

import (
	"io"
	"log"
	"os"
)

var ModLogger *log.Logger

func init() {
	logFile, err := os.Create("/tmp/ynab-exporter.log")
	if err != nil {
		panic(err)
	}
	writer := io.MultiWriter(logFile, os.Stdout)

	ModLogger = log.New(writer, "ynab-exporter: ", log.LstdFlags)
}
