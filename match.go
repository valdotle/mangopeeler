package main

import (
	"embed"
	"image"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/vitali-fedulov/images4"
)

func matchAggregator(icon images4.IconT) bool {
	for _, match := range matchImages {
		if images4.Similar(icon, match) {
			return true
		}
	}

	return false
}

func matchDuplicates(d []dirEntryResponse) {
	for i, search := range d {
		zero := i + 1
		if zero > len(d) {
			break
		}
		for j, match := range d[zero:] {
			if images4.Similar(search.icon, match.icon) {
				var path string
				var dir = filepath.Dir(search.filename)
				if len(d)-j > 2*zero {
					path = match.filename
					logToFile.Printf("[duplicate image] at %s | %s is a duplicate of %s", dir, path, search.filename)
					d = slices.Delete(d, j+zero, j+zero+1)
				} else {
					path = search.filename
					logToFile.Printf("[duplicate image] at %s | %s is a duplicate of %s", dir, path, match.filename)
					d = slices.Delete(d, i, i+1)
				}

				deleteDirEntry(path)

				matchDuplicates(d)
			}
		}
	}
}

//go:embed images
var images embed.FS

func initialize() error {
	return fs.WalkDir(images, "images", initializer)
}

func initializer(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}

	submatch := func(m string) bool { return strings.Contains(path, m) }

	if !slices.ContainsFunc(options.Sites, submatch) {
		return nil
	}

	data, err := images.Open(path)
	if err != nil {
		return err
	}

	defer data.Close()

	img, _, err := image.Decode(data)
	if err != nil {
		return err
	}

	matchImages = append(matchImages, images4.Icon(img))

	return nil
}

var matchImages []images4.IconT
