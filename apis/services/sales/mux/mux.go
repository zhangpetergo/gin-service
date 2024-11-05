// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"github.com/jmoiron/sqlx"
	"github.com/zhangpetergo/gin-service/apis/services/sales/route/sys/checkapi"
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"os"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(build string, log *logger.Logger, db *sqlx.DB, authClient *authclient.Client, shutdown chan os.Signal) *web.App {
	// mux := web.NewApp(shutdown, mid.Trace(), mid.Logger(log), mid.Metrics(), mid.Errors(log), mid.Panics())
	app := web.NewApp(shutdown)
	// Add the routes for the check group.
	checkapi.Routes(build, app, log, db, authClient)
	return app
}
