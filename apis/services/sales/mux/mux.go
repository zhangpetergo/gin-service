// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"github.com/zhangpetergo/gin-service/apis/services/sales/route/sys/checkapi"
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"os"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(log *logger.Logger, auth *auth.Auth, shutdown chan os.Signal) *web.App {
	// mux := web.NewApp(shutdown, mid.Trace(), mid.Logger(log), mid.Metrics(), mid.Errors(log), mid.Panics())
	mux := web.NewApp(shutdown)
	// Add the routes for the check group.
	checkapi.Routes(mux, log, auth)
	return mux
}
