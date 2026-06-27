package dependencies

import (
	"os/exec"

	"github.com/mcoot/dojo-jj/internal/models"
)

type JJClient interface {
	IsJJAvailable() bool
	AddWorkspace(name, destPath, revSet string) (*models.JJWorkspace, error)
}

type JJClientImpl struct{}

func NewJJClient() *JJClientImpl {
	return &JJClientImpl{}
}

func (c *JJClientImpl) IsJJAvailable() bool {
	_, err := exec.LookPath("jj")
	return err == nil
}

func (c *JJClientImpl) AddWorkspace(name, destPath, revSet string) (*models.JJWorkspace, error) {
	return nil, nil
}
