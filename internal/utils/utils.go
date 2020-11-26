package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// RenameToSendToIPFS rename all videos to send to ipfs
func RenameToSendToIPFS(path, resourceID string) {
	var idx int
	idx = 2
	entries, err := ioutil.ReadDir(path)
	SendError("utils.RenameToSendToIPFS.ioutil.ReadDir", err)
	for _, entry := range entries {
		extension := filepath.Ext(entry.Name())
		if extension == ".mp4" {
			sourcePath := filepath.Join(path, entry.Name())
			newPath := filepath.Join(path, fmt.Sprintf("%s_v%d.mp4", resourceID, idx))
			err := os.Rename(sourcePath, newPath)
			SendError("os.Rename", err)
			idx++
		}
	}
}

// VerifyBeforeSendToIPFS verify if has the necessary to send to ipfs
func VerifyBeforeSendToIPFS(path string) bool {
	var hasExtension int
	entries, err := ioutil.ReadDir(path)
	SendError("utils.VerifyBeforeSendToIPFS.ioutil.ReadDir", err)
	for _, entry := range entries {
		extension := filepath.Ext(entry.Name())
		if extension == ".mp4" {
			hasExtension++
		}
	}

	if hasExtension >= 5 {
		return true
	}
	return false
}
