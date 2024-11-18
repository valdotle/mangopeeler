package main

import (
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/vitali-fedulov/images4"
)

func processDirEntry(path string, d fs.DirEntry, result *images4.IconT) error {
	if d.IsDir() {
		if dirThreaded {
			limit.Add(path)

		} else if options.Walk {
			processDir(path)
		}

		return nil
	}

	// not a (supported) image
	if !isSupportedImage(d) {
		return nil
	}

	defer imagesRead.Add(1)

	img, err := images4.Open(path)
	if err != nil {
		logToFile.Printf("[corrupted image] failed to read image at %s", path)

		return nil
	}

	icon := images4.Icon(img)

	if result != nil {
		*result = icon
	}

	if matchAggregator(icon) {
		logToFile.Printf("[aggregator image] found at %s", path)
		deleteDirEntry(path)
	}

	return err
}

func deleteDirEntry(path string) {
	dirEntriesFound.Add(1)
	if options.Delete {
		if err := os.Remove(path); err != nil {
			log.Panicf("failed to delete file, error:\n%s", err.Error())
		} else {
			logToFile.Printf("file deleted successfully")
		}
	}
}

var imageExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

func isSupportedImage(d fs.DirEntry) bool {
	return slices.ContainsFunc(imageExtensions, func(e string) bool { return strings.HasSuffix(d.Name(), e) })
}
