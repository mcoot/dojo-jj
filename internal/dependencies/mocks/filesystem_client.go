package mocks

type FileSystemClient struct{}

func (m *FileSystemClient) MkdirAll(path string) error {
	return nil
}
