package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/vitali-fedulov/images4"
)

var custom string

func initializeCustomDir() {
	abs, err := filepath.Abs(options.Custom)
	if err != nil {
		log.Fatal(err)
	}

	custom = abs
}

func processDir(path string) {
	if dirThreaded {
		defer limit.Remove()
	}

	current, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	if current == custom {
		return
	}

	ds, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("error occured walking specified director%s, error:\n%s", func() string {
			if options.Walk {
				return "ies"
			}
			return "y"
		}(), err.Error())
	}

	var dirProcessor func(d fs.DirEntry)

	// empty directory
	if entries := len(ds); entries == 0 {
		totalReads.Add(1)

		if options.EmptyDir {
			logToFile.Printf("[empty directory] at %s", path)
			deleteDirEntry(path)
		}

		return

	} else
	// without caring for duplicates
	if !options.Duplicates || entries < 2 {
		dirProcessor = dirEntry(path)

	} else
	// with duplicates check
	{
		var results []dirEntryResponse
		dirProcessor = dirEntryWithResult(path, &results)

		defer matchDuplicates(&results)
	}

	for _, d := range ds {
		totalReads.Add(1)
		dirProcessor(d)
	}
}

func dirEntry(path string) func(d fs.DirEntry) {
	return func(d fs.DirEntry) {
		if _, err := processDirEntry(filePath(path, d), d); err != nil {
			log.Fatal(err)
		}
	}
}

func dirEntryWithResult(path string, results *[]dirEntryResponse) func(d fs.DirEntry) {
	return func(d fs.DirEntry) {
		if result, err := processDirEntry(filePath(path, d), d); err != nil {
			log.Fatal(err)
		} else if result != nil {
			*results = append(*results, dirEntryResponse{filepath: filePath(path, d), icon: *result})
		}
	}
}

func filePath(path string, d fs.DirEntry) string {
	return filepath.Join(path, d.Name())
}

type dirEntryResponse struct {
	filepath string
	icon     images4.IconT
}
