// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"github.com/zhangpetergo/gin-service/apis/services/auth/route/authapi"
	"github.com/zhangpetergo/gin-service/apis/services/auth/route/checkapi"
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"os"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(log *logger.Logger, auth *auth.Auth, shutdown chan os.Signal) *web.App {
	//app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	app := web.NewApp(shutdown)

	checkapi.Routes(app, auth)
	authapi.Routes(app, log, auth)

	return app
}
