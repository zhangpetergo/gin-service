package authapi

import (
	"github.com/zhangpetergo/gin-service/app/api/auth"
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth *auth.Auth
	Log  *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	bearer := mid.Bearer(cfg.Auth)
	basic := mid.Basic()
	api := newAPI(cfg.Auth)

	app.GET("/auth/token/:kid", api.token, mid.Trace(), mid.Logger(cfg.Log), basic, mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())
	app.GET("/auth/authenticate", api.authenticate, mid.Trace(), mid.Logger(cfg.Log), bearer, mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())
	app.POST("/auth/authorize", api.authorize, mid.Trace(), mid.Logger(cfg.Log), mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())
}
