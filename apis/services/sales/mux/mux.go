// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"github.com/zhangpetergo/gin-service/apis/services/sales/route/sys/checkapi"
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"os"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(log *logger.Logger, shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown, mid.Trace(), mid.Logger(log), mid.Metrics(), mid.Errors(log), mid.Panics())

	// Add the routes for the check group.
	checkapi.Routes(mux)
	return mux
}
