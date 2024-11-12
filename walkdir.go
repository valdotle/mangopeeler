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
	addJob          chan string
	removeJob, wait chan any
	jobs            []string
	size            uint
}

func (p pool) add(path string) {
	p.addJob <- path
}

func (p pool) remove() {
	p.removeJob <- nil
}

func (p pool) finish() {
	<-p.wait
}

func newPool() pool {
	p := pool{make(chan string, options.DirEntryThreads*options.DirThreads), make(chan any, options.DirThreads), make(chan any), nil, options.DirThreads}

	go p.run()

	return p
}

func (p pool) run() {
	for available := p.size; ; {
		select {
		case <-p.removeJob:
			available++
		case path := <-p.addJob:
			p.jobs = append(p.jobs, path)
		}

		if available > 0 && len(p.jobs) > 0 {
			go processDir(p.jobs[0])
			p.jobs = p.jobs[1:]
			available--
			continue
		}

		// make sure there are no pending jobs before closing, since the order of select isn't deterministic
		if available == p.size {
			select {
			case path := <-p.addJob:
				p.jobs = append(p.jobs, path)
				continue
			default:
			}

			break
		}
	}

	p.wait <- nil
}

func walker(path string) {
	if dirThreaded {
		limit = newPool()
		limit.add(path)
		limit.finish()

	} else {
		processDir(path)
	}
}

const threadingThreshold = 4

func processDir(path string) {
	if dirThreaded {
		defer limit.remove()
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

	if fileThreaded && len(ds) > threadingThreshold {
		entries := len(ds)
		threads := make(chan any, options.DirEntryThreads)
		completed := make(chan any, entries)
		for i := 0; i < int(options.DirEntryThreads); i++ {
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
			limit.add(path)

		} else if options.Walk {
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
		if options.Delete {
			if err = os.Remove(path); err != nil {
				return err
			}

			logToFile.Println("file deleted successfully")
		}
	}

	return nil
}
