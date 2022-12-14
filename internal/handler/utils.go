package handler

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// 2 GB
const limitSize = 2 << 34

func createWorkDir(baseName, prefix string) (string, error) {
	u := uuid.NewV4().String()
	var dirName = ""
	if prefix == "" {
		dirName = fmt.Sprintf("%s_%s", baseName, u)
	} else {
		dirName = fmt.Sprintf("%s_%s_%s", baseName, prefix, u)
	}

	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}

	log.Default().Println("[INFO] workdir was created")

	err = os.Chmod(dirName, 0777)
	if err != nil {
		return "", err
	}

	log.Default().Println("[INFO] chmod was applied on workdir")

	return dirName, nil

}

func removeWorkDir(name string) error {
	err := os.RemoveAll(name)
	if err != nil {
		return err
	}
	return nil
}

// validateRequestContent validates if the number of uploaded files = 1 and  its suffixes are allowed
func (h *Handler) validatePayload(files []*multipart.FileHeader) (int64, error) {
	if len(files) == 0 {
		return 0, errors.New("no files uploaded in request")
	}

	if len(files) > 1 {
		return 0, errors.New("more than 1 file was uploaded")
	}

	allowedSuffixes := []string{".zip"}
	for _, file := range files {
		fileSize := file.Size
		if fileSize >= limitSize {
			h.logger.Error("file size must be less than 2 GB")
			return 0, errors.New("file is too large")
		}

		for _, suffix := range allowedSuffixes {
			if strings.HasSuffix(strings.ToLower(file.Filename), suffix) {
				return fileSize, nil
			}
		}
	}

	h.logger.Error("file must be a zip archive")
	return 0, errors.New("invalid suffix of uploaded file name. Allows only [.zip]")
}
