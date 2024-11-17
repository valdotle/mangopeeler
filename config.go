package main

import _ "embed"

type config struct {
	Delete          bool            `json:"delete"`
	Dir             string          `json:"dir"`
	DirThreads      uint            `json:"directory-threads"`
	DirEntryThreads uint            `json:"directory-entry-threads"`
	Log             bool            `json:"log"`
	LogAt           string          `json:"log-at"`
	Sites           stringArrayFlag `json:"site"`
	Walk            bool            `json:"walk"`
	Duplicates      bool            `json:"duplicates"`
}

//go:embed config.json
var data []byte

const configFileName = "./config.json"
