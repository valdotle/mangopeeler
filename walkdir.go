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

func walker() error {
	if *walk {
		return filepath.WalkDir(path, walkfunc)
	}

	ds, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, d := range ds {
		if err = walkfunc(filepath.Join(path, d.Name()), d, nil); err != nil {
			return err
		}
	}

	return nil
}

var imageExtensions = []string{".png", ".jpg", ".jpeg", ".gif"}

func walkfunc(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	totalReads.Add(1)

	if d.IsDir() {
		return nil
	}

	// not a (supported) image
	if !slices.ContainsFunc(imageExtensions, func(e string) bool { return strings.HasSuffix(d.Name(), e) }) {
		return nil
	}

	img, err := images4.Open(path)
	if err != nil {
		log.Printf("failed to read image at %s ", path)

		return nil
	}

	imagesRead.Add(1)

	if match(images4.Icon(img)) {
		log.Printf("found matching file at %s", path)
		if *delete {
			if err = os.Remove(path); err != nil {
				return err
			}

			log.Println("file deleted successfully")
		}
	}

	return nil
}
