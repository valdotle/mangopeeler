package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vitali-fedulov/images4"
)

var limit pool

type pool struct {
	add, reset, wait chan any
	size             uint
}

func (p pool) get() {
	<-p.add
}

func (p pool) remove() {
	p.reset <- nil
}

func (p pool) finish() {
	<-p.wait
}

func newPool(size uint) pool {
	p := pool{make(chan any), make(chan any, size), make(chan any), size}

	go func(p pool) {
		for available := p.size; ; {
			if available > 0 {
				select {
				case <-p.reset:
					available++
				case p.add <- nil:
					available--
				}
			} else {
				<-p.reset
				available++
			}

			if available == size {
				select {
				case p.add <- nil:
					continue
				default:
				}

				break
			}
		}

		p.wait <- nil
	}(p)

	return p
}

func walker(path string) {
	if dirThreaded {
		limit = newPool(*dirThreads)
		limit.get()
		go processDir(path)
		limit.finish()

	} else {
		processDir(path)
	}
}

func processDir(path string) {
	if dirThreaded {
		defer limit.remove()
	}

	ds, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("error occured walking specified director%s, error:\n%s", func() string {
			if *walk {
				return "ies"
			}
			return "y"
		}(), err.Error())
	}

	if fileThreaded {
		entries := len(ds)
		threads := make(chan any, *dirEntryThreads)
		completed := make(chan any, entries)
		for i := 0; i < int(*dirEntryThreads); i++ {
			threads <- nil
		}

		for _, d := range ds {
			<-threads
			go func() {
				if err = processDirEntry(filepath.Join(path, d.Name()), d); err != nil {
					log.Fatal(err)
				}

				completed <- nil
				threads <- nil
			}()
		}

		for i := 1; i < entries; i++ {
			<-completed
		}

	} else {
		for _, d := range ds {
			if err = processDirEntry(filepath.Join(path, d.Name()), d); err != nil {
				log.Fatal(err)
			}
		}
	}
}

var imageExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

func processDirEntry(path string, d fs.DirEntry) error {
	defer totalReads.Add(1)

	if d.IsDir() {
		if dirThreaded {
			limit.get()
			go processDir(path)

		} else if *walk {
			processDir(path)
		}

		return nil
	}

	// not a (supported) image
	if !slices.ContainsFunc(imageExtensions, func(e string) bool { return strings.HasSuffix(d.Name(), e) }) {
		return nil
	}

	defer imagesRead.Add(1)

	img, err := images4.Open(path)
	if err != nil {
		logToFile.Printf("failed to read image at %s", path)

		return nil
	}

	if match(images4.Icon(img)) {
		logToFile.Printf("found matching file at %s", path)
		if *deleteMatches {
			if err = os.Remove(path); err != nil {
				return err
			}

			logToFile.Println("file deleted successfully")
		}
	}

	return nil
}
