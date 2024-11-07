package testapi

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"math/rand"
	"net/http"
)

type api struct{}

func newAPI() *api {
	return &api{}
}

func (api *api) testError(c *gin.Context) {
	if n := rand.Intn(100); n%2 == 0 {
		c.Error(errs.Newf(errs.FailedPrecondition, "this message is trused"))
		return
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	c.JSON(http.StatusOK, status)

}

func (api *api) testPanic(c *gin.Context) {
	if n := rand.Intn(100); n%2 == 0 {
		panic("WE ARE PANICKING!!!")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	c.JSON(http.StatusOK, status)
}

func (api *api) testAuth(c *gin.Context) {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	c.JSON(http.StatusOK, status)
}
