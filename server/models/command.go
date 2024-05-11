package models

type Command struct {
	Directory string `json:"directory"`
	Command   string `json:"command"`
	Type      string `json:"type"`
	Data      string `json:"data"`
	IsFile    string `json:"isFile"`
}
