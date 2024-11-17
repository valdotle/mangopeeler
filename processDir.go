package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/valdotle/mangopeeler/internal"
	"github.com/vitali-fedulov/images4"
)

const threadingThreshold = 4

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
		entries      = len(ds)
		thread       = fileThreaded && threadingThreshold < entries
	)

	// synchronous without caring for duplicates
	if !options.Duplicates && !thread {
		dirProcessor = func(d fs.DirEntry) { dirEntry(path, d) }

	} else
	// synchronous with duplicates check
	if !thread {
		var results []internal.DirEntryResponse
		dirProcessor = func(d fs.DirEntry) {
			results = append(results, internal.DirEntryResponse{Filename: fileName(path, d), Icon: dirEntryWithResult(path, d)})
		}

		deferfunc = func() { matchDuplicates(results) }

	} else
	// threaded with duplicate check
	if dirPool := internal.NewDirPool(int(options.DirEntryThreads), entries); options.Duplicates {
		dirProcessor = func(d fs.DirEntry) { dirEntryThreadedWithResult(d, path, dirPool) }
		deferfunc = func() { matchDuplicates(dirPool.WaitForResults()) }

	} else
	// threaded without duplicate check
	{
		dirProcessor = func(d fs.DirEntry) { dirEntryThreaded(d, path, dirPool) }
		deferfunc = dirPool.Wait
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

func dirEntryThreaded(d fs.DirEntry, path string, dirPool internal.DirPool) {
	dirPool.Add()
	go func() {
		defer dirPool.Remove()
		dirEntry(path, d)
	}()
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

func dirEntryThreadedWithResult(d fs.DirEntry, path string, dirPool internal.DirPool) {
	dirPool.Add()
	go dirPool.RemoveWithResult(fileName(path, d), dirEntryWithResult(path, d))
}
