package models

type Output struct {
	OldDirectory string `json:"oldDirectory"`
	Directory    string `json:"directory"`
	Output       string `json:"ouput"`
	Error        string `json:"error"`
	Type         string `json:"typ"`
	IsFile       string `json:"isFile"`
	Command      string `json:"command"`
}
