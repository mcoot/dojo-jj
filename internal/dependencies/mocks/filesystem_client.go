package mocks

type FileSystemClient struct{}

func NewFileSystemClient() *FileSystemClient {
	return &FileSystemClient{}
}

func (m *FileSystemClient) MkdirAll(path string) error {
	return nil
}
