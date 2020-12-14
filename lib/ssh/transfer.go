package ssh

import (
	"github.com/hardboiledalex/go-tools/lib/utils"
	"github.com/melbahja/goph"
	"github.com/pkg/sftp"
	"io/ioutil"
	"os"
)

type BinaryFile struct {
	Data []byte
	Mode os.FileMode
}

type TextFile struct {
	Strings []string
	Mode    os.FileMode
}

// Downloads binary file from the remote machine to memory
// Returns file data in bytes, original file permissions and error if happened
func DownloadBinaryFileToMemory(client *goph.Client, downloadPath string) (*BinaryFile, error) {
	if client != nil {
		sftpClient, err := sftp.NewClient(client.Client)
		if err != nil {
			return nil, err
		}
		defer sftpClient.Close()

		remoteFile, err := sftpClient.Open(downloadPath)
		if err != nil {
			return nil, err
		}
		defer remoteFile.Close()

		fileInfo, err := remoteFile.Stat()
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(remoteFile)
		if err != nil {
			return nil, err
		}

		return &BinaryFile{data, fileInfo.Mode()}, nil
	} else {
		localFileInfo, err := os.Stat(downloadPath)
		if err != nil {
			return nil, err
		}
		fileBytes, err := ioutil.ReadFile(downloadPath)
		if err != nil {
			return nil, err
		}
		return &BinaryFile{fileBytes, localFileInfo.Mode()}, nil
	}
}

// Uploads binary file from memory to the remote machine
// Returns an error if happened
func UploadBinaryFileFromMemory(client *goph.Client, uploadPath string, binaryFile *BinaryFile) error {
	if client != nil {
		sftpClient, err := sftp.NewClient(client.Client)
		if err != nil {
			return err
		}
		defer sftpClient.Close()

		file, err := sftpClient.Create(uploadPath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = file.Chmod(binaryFile.Mode)
		if err != nil {
			return err
		}

		_, err = file.Write(binaryFile.Data)
		return err
	} else {
		return ioutil.WriteFile(uploadPath, binaryFile.Data, binaryFile.Mode)
	}
}

// Downloads text file from the remote machine to memory
// Returns an empty structure if some error happened during download
func DownloadTextFileToMemory(client *goph.Client, downloadPath string) *TextFile {
	binaryFile, err := DownloadBinaryFileToMemory(client, downloadPath)
	if err != nil {
		return &TextFile{make([]string, 0), 0644}
	}
	return &TextFile{utils.BytesToStrings(binaryFile.Data), binaryFile.Mode}
}

// Uploads text file from memory to the remote machine
// Returns an error if happened
func UploadTextFileFromMemory(client *goph.Client, uploadPath string, textFile *TextFile) error {
	return UploadBinaryFileFromMemory(client, uploadPath, &BinaryFile{utils.StringsToBytes(textFile.Strings), textFile.Mode})
}
