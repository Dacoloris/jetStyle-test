package services

import (
	"fmt"
	"log"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func wkhtmltopdfConvert(workDirName, htmlFileName string) (string, error) {
	f, err := os.Open(htmlFileName)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Fatal(err)
	}
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", err
	}
	r := wkhtmltopdf.NewPageReader(f)
	pdfg.Dpi.Set(300)
	pdfg.NoCollate.Set(false)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.MarginBottom.Set(10)

	r.EnableLocalFileAccess.Set(true)
	pdfg.AddPage(r)

	err = pdfg.Create()
	if err != nil {
		return "", err
	}

	resp := fmt.Sprintf("%s/output.pdf", workDirName)
	err = pdfg.WriteFile("output.pdf")
	if err != nil {
		return "", err
	}

	return resp, nil
}
