package tools

import (
	"path/filepath"

	"github.com/patnaikankit/Code-Cloud/server/models"
)

type ContainerMap map[string]models.ContainerInfo

func ReadContainersData() (ContainerMap, error) {
	absContextDirectory, err := filepath.Abs("./data/containers.info.json")
	if err != nil {
		return nil, err
	}
}
