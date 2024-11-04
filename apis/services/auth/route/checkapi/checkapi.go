// Package checkapi maintains the web based api for system access.
package checkapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func liveness(c *gin.Context) {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	c.JSON(http.StatusOK, status)
}

func readiness(c *gin.Context) {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	c.JSON(http.StatusOK, status)
}
