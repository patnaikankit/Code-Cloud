package models

type Output struct {
	OldDirectory string `json:"oldDir"`
	Directory    string `json:"dir"`
	Output       string `json:"out"`
	Error        string `json:"error"`
	Type         string `json:"type"`
	IsFile       string `json:"isFile"`
	Command      string `json:"command"`
}
