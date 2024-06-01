package models

type ContainerInfo struct {
	ContainerID   string `json:"containerId"`
	ContainerName string `json:"containerName"`
	Port          int    `json:"port"`
}
