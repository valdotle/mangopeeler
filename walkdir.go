package main

import (
	"github.com/valdotle/mangopeeler/internal"
)

var limit internal.Pool

func walker(path string) {
	if dirThreaded {
		limit = internal.NewPool(options.DirThreads, options.DirEntryThreads, processDir)
		limit.Add(path)
		limit.Finish()

	} else {
		processDir(path)
	}
}

/*
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
}*/

/*
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
*/
