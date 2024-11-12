package main

import (
	"embed"
	"image"
	"io/fs"
	"slices"
	"strings"

	"github.com/vitali-fedulov/images4"
)

func match(icon images4.IconT) bool {
	for _, match := range matchImages {
		if images4.Similar(icon, match) {
			return true
		}
	}

	return false
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
