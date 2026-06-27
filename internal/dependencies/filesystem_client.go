package dependencies

import "os"

type FileSystemClient interface {
	MkdirAll(path string) error
}

type FileSystemClientImpl struct{}

func NewFileSystemClient() *FileSystemClientImpl {
	return &FileSystemClientImpl{}
}

func (f *FileSystemClientImpl) MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}
