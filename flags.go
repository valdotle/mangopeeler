package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	path       string
	walk       = flag.Bool("walk", true, "whether to walk subdirectories (if applicable)")
	delete     = flag.Bool("delete", true, "whether to delete located duplicates")
	logPath    = flag.String("log-path", "", "where to store logfiles (if applicable)")
	createLogs = flag.Bool("log", true, "whether to log actions performed by the script")
	sitelist   = siteEnum()
	sites      = sitelist
)

type stringArrayFlag []string

// String is an implementation of the flag.Value interface
func (i *stringArrayFlag) String() string {
	return strings.Join(*i, ", ")
}

// Set is an implementation of the flag.Value interface
func (i *stringArrayFlag) Set(value string) error {
	if !slices.Contains(sitelist, value) {
		return fmt.Errorf("%s is not a valid site value, must be one of %v", value, sitelist)
	}
	*i = append(*i, value)
	return nil
}

func siteEnum() stringArrayFlag {
	ds, err := images.ReadDir("images")
	if err != nil {
		log.Fatalf("failed to read directory of images to match against, error:\n%s", err.Error())
	}

	var sites stringArrayFlag

	for _, d := range ds {
		if !d.IsDir() {
			continue
		}
		sites = append(sites, d.Name())
	}

	return sites
}

func setupFlags() {
	flag.Var(&sites, "site", "which site(s)'s images to check for")

	dir, err := os.Getwd()
	if err != nil {
		log.Panicf("failed to find workdir, error:\n%s", err.Error())
	}

	flag.StringVar(&path, "path", filepath.Clean(dir), "the directory to execute this script in")
}
