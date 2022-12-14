package handler

import (
	"fmt"
	"jetStyle-test/services"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ConvertHtmlToPDF(c *gin.Context) {
	start := time.Now()

	form, err := c.MultipartForm()
	if err != nil {
		msgErr := fmt.Sprintf("[ERROR] receiving request. detail: %s ", err.Error())
		log.Default().Print(msgErr)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": msgErr})
		return
	}
	files := form.File["files"]
	log.Default().Println("[INFO] multipart form request was received successfully")

	fileSize, err := h.validatePayload(files)
	if err != nil {
		msgErr := "[ERROR] validating request payload. detail: " + err.Error()
		log.Default().Print(msgErr)
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": msgErr})
		return
	}
	log.Default().Println("[INFO]", "request payload was validated")

	// creates workdir
	workDirName, err := createWorkDir("html2pdf", "")
	if err != nil {
		msgErr := "[ERROR] creating workdir. detail: " + err.Error()
		log.Default().Print(msgErr)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": msgErr})
		return
	}
	log.Default().Printf("[INFO] workdir %s was created", workDirName)

	// save uploaded file in the workdir
	var uploadedFileName string
	for _, file := range files {
		pathSep := fmt.Sprintf("%c", os.PathSeparator)
		uploadedFileName = filepath.Base(file.Filename)
		fileNameInWorkDir := workDirName + pathSep + uploadedFileName
		if err := c.SaveUploadedFile(file, fileNameInWorkDir); err != nil {
			// Removes workdir
			err = removeWorkDir(workDirName)
			if err != nil {
				log.Default().Printf("[ERROR] deleting used workdir %s. detail: %s", workDirName, err.Error())
			}
			log.Default().Printf("[INFO] executed clean up of workdir: %s", workDirName)

			msgErr := fmt.Sprintf("[ERROR] saving uploaded file %s in workdir. detail: %s", fileNameInWorkDir, err.Error())
			log.Default().Print(msgErr)
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": msgErr})
			return
		}
	}
	log.Default().Printf("[INFO] file %v was uploaded and saved properly in the server in the workdir", uploadedFileName)

	// process the html to pdf convertion
	var pdfFilePath string
	pdfFilePath, err = services.Html2Pdf(workDirName, uploadedFileName)
	if err != nil {
		// Removes workdir
		err = removeWorkDir(workDirName)
		if err != nil {
			log.Default().Printf("[ERROR] deleting used workdir %s. detail: %s", workDirName, err.Error())
		}
		log.Default().Printf("[INFO] executed clean up of workdir: %s", workDirName)

		msgErr := fmt.Sprintf("[ERROR] executing zip2pdf service, detail: %s ", err.Error())
		log.Default().Print(msgErr)
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": msgErr})
		return
	}
	log.Default().Println("[INFO]", "zip to pdf service was executed")

	// response: pdf file
	c.File(pdfFilePath)
	log.Default().Printf("[INFO] html to pdf executed successfully. pdf file generated: %s", pdfFilePath)

	// Removes workdir
	err = removeWorkDir(workDirName)
	if err != nil {
		log.Default().Printf("[ERROR] deleting used workdir %s. detail: %s", workDirName, err.Error())
	}
	log.Default().Printf("[INFO] executed clean up of workdir: %s", workDirName)

	h.logger.Info(
		"complete",
		zap.String("filename", uploadedFileName),
		zap.String("date", time.Now().Format("2006-01-02")),
		zap.String("execution_time", time.Since(start).String()),
		zap.String("cost_size", strconv.FormatInt(fileSize, 10)),
	)
}
