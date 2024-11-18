package main

import _ "embed"

type config struct {
	Custom     string          `json:"custom-images"`
	Delete     bool            `json:"delete"`
	Dir        string          `json:"directory"`
	Duplicates bool            `json:"duplicates"`
	EmptyDir   bool            `json:"empty-dir"`
	Log        bool            `json:"log"`
	LogAt      string          `json:"log-at"`
	Sites      stringArrayFlag `json:"site"`
	Threads    uint            `json:"threads"`
	Walk       bool            `json:"walk"`
}

//go:embed config.json
var data []byte

const configFileName = "./config.json"
