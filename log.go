package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var (
	logToFile = log.New(io.Discard, "", 0)
	logfile   *os.File
)

func setupLogs() {
	if *createLogs {
		if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
			log.Panicf("failed to open logfile directory, error:\n%s", err.Error())
		}

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Panicf("failed to open logfile, error:\n%s", err.Error())
		}

		logfile = file

		logToFile = log.New(logfile, "", 0)
	}

	log.SetPrefix("\r")
	log.SetFlags(log.Ltime)
}

func closeLogs() {
	if logfile != nil {
		logfile.Close()
	}
}

var (
	imagesRead, totalReads atomic.Uint32
	start                  time.Time
)

func progress() {
	start = time.Now()
	for {
		time.Sleep(time.Second)

		seconds := time.Since(start).Seconds()
		total := totalReads.Load()
		images := imagesRead.Load()
		fmt.Printf("\rdirectory entries scanned: %d (%.2f/s); thereof images: %d (%.2f/s)", total, float64(total)/seconds, images, float64(images)/seconds)
	}
}
