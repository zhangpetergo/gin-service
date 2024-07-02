// Package checkapi maintains the web based api for system access.
package checkapi

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"math/rand"
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

func testError(c *gin.Context) {
	if n := rand.Intn(100); n%2 == 0 {
		c.Error(errs.Newf(errs.FailedPrecondition, "this message is trused"))
		return
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	c.JSON(200, status)

}
