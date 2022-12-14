package router

import (
	"jetStyle-test/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(h *handler.Handler) http.Handler {
	r := gin.Default()

	r.POST("/upload", h.ConvertHtmlToPDF)
	r.GET("/ping", h.Ping)

	return r
}
