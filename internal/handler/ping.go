package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:operation GET /ping html_to_pdf ping
// Returns pong if the server is running
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "pong"})
}
