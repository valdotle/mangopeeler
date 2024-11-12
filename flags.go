package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var (
	options                   config
	sitelist                  = siteEnum()
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
	"delete":                  "del",
	"directory-entry-threads": "det",
	"directory-threads":       "dt",
	"log-at":                  "lat",
	"log":                     "l",
	"site":                    "s",
	"walk":                    "w",
}

type config struct {
	Delete          bool            `json:"delete"`
	Dir             string          `json:"dir"`
	DirThreads      uint            `json:"directory-threads"`
	DirEntryThreads uint            `json:"directory-entry-threads"`
	Log             bool            `json:"log"`
	LogAt           string          `json:"log-at"`
	Sites           stringArrayFlag `json:"site"`
	Walk            bool            `json:"walk"`
}

//go:embed config.json
var data []byte

const configFileName = "./config.json"

func setupFlags() {
	// read config file
	fileData, err := os.ReadFile(configFileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Panicf("failed to open config, error:\n%s", err.Error())
	} else if err == nil {
		data = fileData
	}

	// parse config
	if err := json.Unmarshal(data, &options); err != nil {
		log.Panicf("failed to read config, error\n:%s", err.Error())
	}

	// set flags
	flag.BoolVar(&options.Delete, "delete", options.Delete, "whether to delete located duplicates")
	flag.StringVar(&options.Dir, "dir", options.Dir, "the directory to execute this script in")
	flag.UintVar(&options.DirThreads, "directory-threads", options.DirThreads, "how many directories to process simultaneously (if applicable)")
	flag.UintVar(&options.DirEntryThreads, "directory-entry-threads", options.DirEntryThreads, "how many directory entries to process simultaneously")
	flag.BoolVar(&options.Log, "log", options.Log, "whether to create logfiles for actions performed by the script")
	flag.StringVar(&options.LogAt, "log-at", options.LogAt, "where to store logfiles (if applicable)")
	flag.Var(&options.Sites, "site", "which site(s)'s images to check for")
	flag.BoolVar(&options.Walk, "walk", options.Walk, "whether to walk subdirectories (if applicable)")

	// set flag aliases
	for from, to := range flagAliases {
		flagSet := flag.Lookup(from)
		flag.Var(flagSet.Value, to, "shorthand for "+flagSet.Name)
	}

	flag.Parse()

	dirThreaded = options.DirThreads > 1 && options.Walk
	fileThreaded = options.DirEntryThreads > 1
}
