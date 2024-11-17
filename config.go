package main

import _ "embed"

type config struct {
	Delete     bool            `json:"delete"`
	Dir        string          `json:"directory"`
	Duplicates bool            `json:"duplicates"`
	Log        bool            `json:"log"`
	LogAt      string          `json:"log-at"`
	Sites      stringArrayFlag `json:"site"`
	Threads    uint            `json:"threads"`
	Walk       bool            `json:"walk"`
}

//go:embed config.json
var data []byte

const configFileName = "./config.json"
