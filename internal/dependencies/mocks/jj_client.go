package mocks

import "github.com/mcoot/dojo-jj/internal/models"

type JJClient struct {
	jjInstalled bool
}

func NewJJClient() *JJClient {
	return &JJClient{
		jjInstalled: true,
	}
}

func (m *JJClient) IsJJAvailable() bool {
	return m.jjInstalled
}

func (m *JJClient) AddWorkspace(name, destPath, revSet string) (*models.JJWorkspace, error) {
	return nil, nil
}

func (m *JJClient) SetJJInstalled(installed bool) {
	m.jjInstalled = installed
}
