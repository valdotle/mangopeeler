package main

import (
	"log"
	"time"

	"github.com/valdotle/mangopeeler/internal"
)

func main() {
	setupFlags()
	setupLogs()
	defer closeLogs()

	if err := initialize(); err != nil {
		log.Panicf("failed to load images to match against, error:\n%s", err.Error())
	}

	log.Println("mangopeeler initialized successfully")

	go progress()

	walker(options.Dir)

	log.Printf("\nscanned %d entries (including %d images) in %s", totalReads.Load(), imagesRead.Load(), time.Since(start).String())
}

var limit internal.Pool

func walker(path string) {
	if dirThreaded {
		limit = internal.NewPool(options.Threads, processDir)
		limit.Add(path)
		limit.Finish()

	} else {
		processDir(path)
	}
}
