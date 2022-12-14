package services

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
)

func imagesEmbedder(htmlFileName string, imageFileNames []string) (string, error) {

	content, err := os.ReadFile(htmlFileName)
	if err != nil {
		return "", err
	}
	html := string(content)

	imgSrcs, err := getAllImgSources(&html)
	if err != nil {
		return "", err
	}

	replaceMap, err := getReplaceMap(imgSrcs, imageFileNames)
	if err != nil {
		return "", nil
	}

	for imgsrc, imgBase64 := range replaceMap {
		html = embedImageInHtml(imgsrc, imgBase64, html)
	}

	newHtmlFileName := createNewHtmlFileName(htmlFileName)

	err = writeHtmlToFile(newHtmlFileName, html)
	if err != nil {
		return "", err
	}

	return newHtmlFileName, nil
}

func createNewHtmlFileName(oldHtmlFileName string) string {
	basename := strings.TrimSuffix(oldHtmlFileName, ".html")
	return basename + "_base64_images" + ".html"
}

func writeHtmlToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}

func getAllImgSources(html *string) ([]*string, error) {

	// doc, err := htmlquery.LoadDoc(*html)
	doc, err := htmlquery.Parse(strings.NewReader((*html)))
	if err != nil {
		return nil, err
	}

	imgs, err := htmlquery.QueryAll(doc, "//img")
	if err != nil {
		return nil, err
	}

	resp := []*string{}
	for _, img := range imgs {
		src := htmlquery.SelectAttr(img, "src")
		resp = append(resp, &src)
	}

	return resp, nil
}

func embedImageInHtml(imgsrc *string, imgBase64 *string, html string) string {

	resp := strings.ReplaceAll(html, *imgsrc, *imgBase64)

	return resp
}

func getReplaceMap(imgSrcs []*string, imagesFileNames []string) (map[*string]*string, error) {

	imgSrcUploadedFileNameMap, err := buildImgSrcUploadedFileNameMap(imgSrcs, imagesFileNames)
	if err != nil {
		return nil, err
	}

	resp := make(map[*string]*string)
	for imgSrc, fileName := range imgSrcUploadedFileNameMap {
		resp[imgSrc] = calcImageBase64String(fileName)
	}

	return resp, nil
}

func buildImgSrcUploadedFileNameMap(imgSrcs []*string, imageFileNames []string) (map[*string]string, error) {
	resp := map[*string]string{}
	for _, uploadedFileName := range imageFileNames {
		src := discoverAssociatedHtmlImgSrc(uploadedFileName, imgSrcs)
		if src != nil {
			resp[src] = uploadedFileName
		}
	}
	return resp, nil
}

func discoverAssociatedHtmlImgSrc(fName string, imgSrcs []*string) *string {

	// try exact matching
	respExactMatching := exactMatching(fName, imgSrcs)
	if respExactMatching != nil {
		return respExactMatching
	}

	// try pattern 'ID=<number>' where number = image ID

	return nil
}

func exactMatching(fName string, imgSrcs []*string) *string {
	for _, imgsrc := range imgSrcs {
		if fName == *imgsrc {
			return imgsrc
		}
	}
	return nil
}

func calcImageBase64String(fName string) *string {

	bytes, err := os.ReadFile(fName)
	if err != nil {
		log.Default().Printf("error reading file: %s", fName)
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	case "image/gif":
		base64Encoding += "data:image/gif;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	return &base64Encoding
}

// toBase64 return the base64 string representation of the []byte
func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
