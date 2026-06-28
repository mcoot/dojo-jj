package models

type WorkspaceState struct {
	CurrentLease *LeaseID `json:"current_lease,omitempty"`
}

type State struct {
	Pools map[RepoRoot]map[WorkspaceID]WorkspaceState `json:"pools"`
}
