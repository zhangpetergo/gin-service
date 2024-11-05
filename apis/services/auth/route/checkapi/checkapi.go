// Package checkapi maintains the web based api for system access.
package checkapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/zhangpetergo/gin-service/business/data/sqldb"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"net/http"
	"os"
	"runtime"
	"time"
)

type api struct {
	build string
	log   *logger.Logger
	db    *sqlx.DB
}

func newAPI(build string, log *logger.Logger, db *sqlx.DB) *api {
	return &api{
		build: build,
		db:    db,
		log:   log,
	}
}

func (api *api) readiness(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK

	if err := sqldb.StatusCheck(ctx, api.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		api.log.Info(ctx, "readiness failure", "status", status)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	c.JSON(statusCode, data)
}

func (api *api) liveness(c *gin.Context) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}
	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS int    `json:"GOMAXPROCS,omitempty"`
	}{
		Status:     "up",
		Build:      api.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	c.JSON(http.StatusOK, data)
}
