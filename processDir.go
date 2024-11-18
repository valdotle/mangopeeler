package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/vitali-fedulov/images4"
)

func processDir(path string) {
	if dirThreaded {
		defer limit.Remove()
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

	var (
		dirProcessor func(d fs.DirEntry)
		deferfunc    = func() {}
	)

	// empty directory
	if entries := len(ds); entries == 0 {
		if options.EmptyDir {
			logToFile.Printf("[empty directory] at %s", path)
			deleteDirEntry(path)
		}

		return

	} else
	// without caring for duplicates
	if !options.Duplicates || entries < 2 {
		dirProcessor = func(d fs.DirEntry) { dirEntry(path, d) }

	} else
	// with duplicates check
	{
		var results []dirEntryResponse
		dirProcessor = func(d fs.DirEntry) {
			results = append(results, dirEntryResponse{filename: fileName(path, d), icon: dirEntryWithResult(path, d)})
		}

		deferfunc = func() { matchDuplicates(results) }
	}

	defer deferfunc()

	for _, d := range ds {
		dirProcessor(d)
	}
}

func dirEntry(path string, d fs.DirEntry) {
	if err := processDirEntry(fileName(path, d), d, nil); err != nil {
		log.Fatal(err)
	}
}

func dirEntryWithResult(path string, d fs.DirEntry) images4.IconT {
	var result images4.IconT
	if err := processDirEntry(fileName(path, d), d, &result); err != nil {
		log.Fatal(err)
	}

	return result
}

func fileName(path string, d fs.DirEntry) string {
	return filepath.Join(path, d.Name())
}

type dirEntryResponse struct {
	filename string
	icon     images4.IconT
}
