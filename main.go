package main

import (
	"log"
	"time"
)

func main() {
	setupFlags()
	setupLogs()
	defer logfile.Close()

	if err := initialize(); err != nil {
		log.Panicf("failed to load images to match against, error:\n%s", err.Error())
	}

	log.Println("mangopeeler initialized successfully")

	go progress()

	walker(path)

	log.Printf("\nscanned %d entries (including %d images) in %s", totalReads.Load(), imagesRead.Load(), time.Since(start).String())
}
