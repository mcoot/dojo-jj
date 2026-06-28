package mocks

import "github.com/mcoot/dojo-jj/internal/models"

type JJClient struct {
	jjInstalled bool

	repoRoot         models.RepoRoot
	getRepoRootError error

	workspaces []*models.JJWorkspace

	listWorkspacesError error
}

func NewJJClient() *JJClient {
	return &JJClient{
		repoRoot:    "/default-repo-root",
		jjInstalled: true,
	}
}

func (m *JJClient) IsJJAvailable() bool {
	return m.jjInstalled
}

func (m *JJClient) GetRepoRoot() (models.RepoRoot, error) {
	if m.getRepoRootError != nil {
		return "", m.getRepoRootError
	}
	return m.repoRoot, nil
}

func (m *JJClient) ListWorkspaces() ([]*models.JJWorkspace, error) {
	if m.listWorkspacesError != nil {
		return nil, m.listWorkspacesError
	}
	return m.workspaces, nil
}

func (m *JJClient) AddWorkspace(name, destPath, revSet string) (*models.JJWorkspace, error) {
	return nil, nil
}

func (m *JJClient) SetJJInstalled(installed bool) {
	m.jjInstalled = installed
}

func (m *JJClient) SetRepoRoot(repoRoot models.RepoRoot) {
	m.repoRoot = repoRoot
}

func (m *JJClient) SetGetRepoRootError(err error) {
	m.getRepoRootError = err
}

func (m *JJClient) SetListWorkspacesError(err error) {
	m.listWorkspacesError = err
}
