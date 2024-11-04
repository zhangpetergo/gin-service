package authapi

import (
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, a *auth.Auth) {

	authen := mid.AuthenticateLocal(a)

	api := newAPI(a)
	app.GET("/auth/token/:kid", api.token, mid.Trace(), mid.Logger(log), authen, mid.Metrics(), mid.Errors(log), mid.Panics())
	app.GET("/auth/authenticate", api.authenticate, mid.Trace(), mid.Logger(log), authen, mid.Metrics(), mid.Errors(log), mid.Panics())
	app.POST("/auth/authorize", api.authorize, mid.Trace(), mid.Logger(log), mid.Metrics(), mid.Errors(log), mid.Panics())
}
