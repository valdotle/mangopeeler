package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
	"time"
)

func setupLogs() {
	log.SetFlags(0)
	log.SetPrefix("\r")

	if *logPath != "" && *createLogs {
		logfile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Panicf("failed to open logfile, error:\n%s", err.Error())
		}

		defer logfile.Close()

		log.SetOutput(logfile)
	} else if !*createLogs {
		log.SetOutput(io.Discard)
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
