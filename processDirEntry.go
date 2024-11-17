package main

import (
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/vitali-fedulov/images4"
)

var imageExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

func processDirEntry(path string, d fs.DirEntry, result *images4.IconT) error {
	defer totalReads.Add(1)

	if d.IsDir() {
		if dirThreaded {
			limit.Add(path)

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

	icon := images4.Icon(img)

	if result != nil {
		*result = icon
	}

	if matchAggregator(icon) {
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
