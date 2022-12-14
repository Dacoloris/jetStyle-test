package services

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

// Html2Pdf  executes the main process that converts a html file to a pdf file
func Html2Pdf(workdirName, filename string) (string, error) {
	log.Default().Printf("[INFO] service Zip2Pdf was called. WorkDir= %s, File= %s", workdirName, filename)

	if err := changeCurrentDir(workdirName); err != nil {
		return "", err
	}
	defer changeCurrentDir("..")
	log.Default().Println("[INFO] current dir was changed to workdir")

	unzippedFilenames, err := Unzip(filename, ".")
	if err != nil {
		return "", err
	}
	log.Default().Println("[INFO] Uploaded file was unzipped successfully")

	htmlFileName, err := getHtmlFileNameFromUnzippedFileNames(unzippedFilenames)
	if err != nil {
		return "", err
	}
	log.Default().Printf("[INFO] html file name was identified: %s", htmlFileName)

	imagesFileNames := getImageFileNamesFromUnzipFileNames(unzippedFilenames)
	log.Default().Printf("[INFO] all image file names was identified, %v ", imagesFileNames)

	// only if html contains images embed them.
	if len(imagesFileNames) > 0 {
		htmlFileName, err = imagesEmbedder(htmlFileName, imagesFileNames)
		if err != nil {
			return "", err
		}
	}

	pdfFileFullName, err := wkhtmltopdfConvert(workdirName, htmlFileName)
	if err != nil {
		return "", err
	}
	log.Default().Println("[INFO] wkhtmltopdf  has worked successfully")
	log.Default().Println("[INFO] service Zip2Pdf finished successfully")
	return pdfFileFullName, nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src, dest string) ([]string, error) {
	var filenames []string
	format := archiver.Zip{}

	err := format.Unarchive(src, dest)
	if err != nil {
		return filenames, err
	}

	files, err := os.ReadDir(dest)
	if err != nil {
		return filenames, err
	}

	for _, f := range files {
		fpath := filepath.Join(dest, f.Name())
		filenames = append(filenames, fpath)
	}

	return filenames, nil
}

func changeCurrentDir(dirName string) error {
	if err := os.Chdir(dirName); err != nil {
		return err
	}
	return nil
}

func getHtmlFileNameFromUnzippedFileNames(fNames []string) (string, error) {
	for _, name := range fNames {
		if strings.HasSuffix(name, ".html") {
			return name, nil
		}
	}
	return "", errors.New("no html file don't present in uploaded data")
}

func getImageFileNamesFromUnzipFileNames(fNames []string) []string {
	var resp []string
	for _, name := range fNames {
		switch {
		case strings.HasSuffix(name, ".png"):
			resp = append(resp, name)
		case strings.HasSuffix(name, ".jpg"):
			resp = append(resp, name)
		case strings.HasSuffix(name, ".jpeg"):
			resp = append(resp, name)
		case strings.HasSuffix(name, ".gif"):
			resp = append(resp, name)
		}
	}
	return resp
}
