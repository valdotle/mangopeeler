package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
)

var (
	options     config
	dirThreaded bool
)

var flagAliases = map[string]string{
	"broken-images":   "bri",
	"custom-images":   "c",
	"delete":          "del",
	"directory":       "dir",
	"duplicates":      "dup",
	"empty-directory": "edr",
	"log-at":          "lat",
	"log":             "l",
	"threads":         "t",
	"walk":            "w",
}

func setupFlags() {
	// read config file
	fileData, err := os.ReadFile(configFileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Panicf("failed to open config, error:\n%s", err.Error())
	} else if err != nil {
		data = fileData
	}

	// parse config
	if err := json.Unmarshal(data, &options); err != nil {
		log.Panicf("failed to read config, error\n:%s", err.Error())
	}

	// set flags
	flag.StringVar(&options.Custom, "custom-images", options.Custom, "path to directory of custom filter images to use")
	flag.BoolVar(&options.Broken, "broken-images", options.Broken, "whether to check for non-parseable images")
	flag.BoolVar(&options.Delete, "delete", options.Delete, "whether to delete directory entries matching your search criteria")
	flag.StringVar(&options.Dir, "directory", options.Dir, "the directory to execute this script in")
	flag.BoolVar(&options.Duplicates, "duplicates", options.Duplicates, "whether to check for duplicate images within directories")
	flag.BoolVar(&options.EmptyDir, "empty-directory", options.EmptyDir, "whether to check for empty directories")
	flag.BoolVar(&options.Log, "log", options.Log, "whether to create logfiles for actions performed by the script")
	flag.StringVar(&options.LogAt, "log-at", options.LogAt, "where to store logfiles (if applicable)")
	flag.UintVar(&options.Threads, "threads", options.Threads, "how many directories to process simultaneously (if applicable)")
	flag.BoolVar(&options.Walk, "walk", options.Walk, "whether to walk subdirectories (if applicable)")

	// set flag aliases
	for from, to := range flagAliases {
		flagSet := flag.Lookup(from)
		flag.Var(flagSet.Value, to, "shorthand for "+flagSet.Name)
	}

	flag.Parse()

	dirThreaded = options.Threads > 1 && options.Walk
}
