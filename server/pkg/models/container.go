package models

type ContainerInfo struct {
	ContainerID   string `json:"containerID"`
	ContainerName string `json:"containername"`
	Port          int    `json:"port"`
}
