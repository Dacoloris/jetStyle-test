// Html To Pdf API
//
// This is a simple converter from html archive to pdf.
//
//		Schemes: http
//	 Host: localhost:8080
//		BasePath: /
//		Version: 1.0.0
//		Contact: https://github.com/Dacoloris
//
//		Consumes:
//		- application/zip
//
//		Produces:
//		- application/json
//
// swagger:meta
package main

import (
	"context"
	"jetStyle-test/internal/app"
	"jetStyle-test/internal/config"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println(a.HTTPServer.ListenAndServe())
		cancel()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
	}
	a.HTTPServer.Shutdown(context.Background())
}
