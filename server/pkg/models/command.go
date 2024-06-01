package models

type Command struct {
	Directory string `json:"dir"`
	Command   string `json:"command"`
	Type      string `json:"type"`
	Data      string `json:"data"`
	IsFile    string `json:"isFile"`
}
