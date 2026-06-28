package models

type RepoRoot string

type WorkspaceID string

type LeaseID string

type JJWorkspace struct {
	Name     WorkspaceID
	ChangeID string
	Root     string
}

type WorkspaceLease struct {
	JJWorkspace
	RepoRoot RepoRoot
	Lease    LeaseID
}
