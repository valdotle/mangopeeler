package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	setupFlags()
	flag.Parse()
	setupLogs()

	if err := initialize(); err != nil {
		log.Panicf("failed to load images to match against, error:\n%s", err.Error())
	}

	log.Println("mangopeeler initialized successfully")

	go progress()

	if err := walker(); err != nil {
		log.Fatalf("error occured walking specified director%s, error:\n%s", func() string {
			if *walk {
				return "ies"
			}
			return "y"
		}(), err.Error())
	}

	log.Printf("scanned %d entries (including %d images) in %s", totalReads.Load(), imagesRead.Load(), time.Since(start).String())
}
