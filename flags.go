package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var (
	path, logPath             string
	walk                      = flag.Bool("walk", true, "whether to walk subdirectories (if applicable)")
	deleteMatches             = flag.Bool("delete", true, "whether to delete located duplicates")
	createLogs                = flag.Bool("log", true, "whether to create logfiles for actions performed by the script")
	dirThreads                = flag.Uint("directory-threads", 20, "how many directories to process simultaneously (if applicable)")
	dirEntryThreads           = flag.Uint("directory-entry-threads", 5, "how many directory entries to process simultaneously")
	sitelist                  = siteEnum()
	sites                     = sitelist
	dirThreaded, fileThreaded bool
)

type stringArrayFlag []string

// String is an implementation of the flag.Value interface
func (i *stringArrayFlag) String() string {
	return strings.Join(*i, ", ")
}

// Set is an implementation of the flag.Value interface
func (i *stringArrayFlag) Set(value string) error {
	if !slices.Contains(sitelist, value) {
		return fmt.Errorf("%s is not a valid site value, must be one of %s", value, sitelist.String())
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

var flagAliases = map[string]string{
	"delete":                  "d",
	"directory-entry-threads": "det",
	"directory-threads":       "dt",
	"log-at":                  "lat",
	"log":                     "l",
	"site":                    "s",
	"walk":                    "w",
}

func setupFlags() {
	flag.Var(&sites, "site", "which site(s)'s images to check for")

	dir, err := os.Getwd()
	if err != nil {
		log.Panicf("failed to find workdir, error:\n%s", err.Error())
	}

	logfileName := strings.ReplaceAll(time.Now().Local().Format(time.DateTime)+".log", ":", "-")
	flag.StringVar(&path, "dir", dir, "the directory to execute this script in")
	flag.StringVar(&logPath, "log-at", filepath.Join(dir, "mango peels", logfileName), "where to store logfiles (if applicable)")

	for from, to := range flagAliases {
		flagSet := flag.Lookup(from)
		flag.Var(flagSet.Value, to, "shorthand for "+flagSet.Name)
	}
	flag.Parse()

	dirThreaded = *dirThreads > 1 && *walk
	fileThreaded = *dirEntryThreads > 1
}
