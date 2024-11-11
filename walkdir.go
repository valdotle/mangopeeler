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

const maxThreads = 100

const perDirectory = 10

var limit = newPool(maxThreads / perDirectory)

type pool struct {
	add, reset, wait chan any
	size             int
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

func newPool(size int) pool {
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
				break
			}
		}

		p.wait <- nil
	}(p)

	return p
}

func walker(path string) {
	limit.get()
	go walkDir(path)
	limit.finish()
}

func walkDir(path string) {
	defer limit.remove()

	ds, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("error occured walking specified director%s, error:\n%s", func() string {
			if *walk {
				return "ies"
			}
			return "y"
		}(), err.Error())
	}

	entries := len(ds)
	threads := make(chan any, perDirectory-1)
	completed := make(chan any, entries)
	for i := 0; i < perDirectory-1; i++ {
		threads <- nil
	}

	for _, d := range ds {
		<-threads
		go func() {
			if err = walkfunc(filepath.Join(path, d.Name()), d); err != nil {
				log.Fatal(err)
			}

			completed <- nil
			threads <- nil
		}()
	}

	for i := 1; i < entries; i++ {
		<-completed
	}
}

var imageExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

func walkfunc(path string, d fs.DirEntry) error {
	totalReads.Add(1)

	if d.IsDir() {
		limit.get()
		go walkDir(path)

		return nil
	}

	// not a (supported) image
	if !slices.ContainsFunc(imageExtensions, func(e string) bool { return strings.HasSuffix(d.Name(), e) }) {
		return nil
	}

	img, err := images4.Open(path)
	if err != nil {
		logToFile.Printf("failed to read image at %s", path)

		return nil
	}

	imagesRead.Add(1)

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
