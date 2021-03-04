package service

import (
	"bytes"
	"fmt"
	"os"
)

// FileStore is interface for FileServer
type FileStore interface {
	Save(SpeakerID string, ModelID string, fileName string, fileData bytes.Buffer) error
}

// DiskFileStore is store file info
type DiskFileStore struct{}

// NewDiskFileStore return a new DiskFileStore
func NewDiskFileStore() *DiskFileStore {
	return &DiskFileStore{}
}

// Save add a new file on disk
func (store *DiskFileStore) Save(
	SpeakerID string,
	ModelID string,
	fileName string,
	fileData bytes.Buffer) error {
	fmt.Printf("start Save : %s, %s, %s\n", SpeakerID, ModelID, fileName)
	// create directory
	filePath := fmt.Sprintf("/Users/jjkim/workspace/src/sampleProject/grpcProject/files/%s/%s",
		SpeakerID, ModelID)
	os.MkdirAll(filePath, 0755)

	fullPath := fmt.Sprintf("%s/%s", filePath, fileName)
	fmt.Println("fullPath : ", fullPath)
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("Create : ", err)
		return err
	}
	defer file.Close()

	_, err = fileData.WriteTo(file)
	if err != nil {
		fmt.Println("WriteTo : ", err)
		return err
	}

	return nil
}
