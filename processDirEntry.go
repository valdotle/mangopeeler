package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/vitali-fedulov/images4"
)

func processDirEntry(path string, d fs.DirEntry) (iconP *images4.IconT, err error) {
	if d.IsDir() {
		if dirThreaded {
			limit.Add(path)

		} else if options.Walk {
			processDir(path)
		}

		return
	}

	// not a (supported) image extension
	if !isSupportedImage(d) {
		return
	}

	defer imagesRead.Add(1)

	img, err := images4.Open(path)
	if err != nil {
		logToFile.Printf("[corrupted image] failed to read image at %s", path)

		if options.Broken {
			deleteDirEntry(path)
		}

		return nil, nil
	}

	icon := images4.Icon(img)

	if matchCustom(icon) {
		logToFile.Printf("[custom image] found at %s", path)
		deleteDirEntry(path)
		return
	}

	return &icon, nil
}

func deleteDirEntry(path string) {
	dirEntriesFound.Add(1)
	if options.Delete {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Panicf("failed to delete file, error:\n%s", err.Error())
		} else {
			logToFile.Printf("file deleted successfully")
		}
	}
}

var imageExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

func isSupportedImage(d fs.DirEntry) bool {
	return slices.ContainsFunc(imageExtensions, func(e string) bool { return strings.HasSuffix(strings.ToLower(d.Name()), e) })
}
