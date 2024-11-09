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

	// 额，应该把中间件放在前面，处理程序放在最后
	app.GET("/auth/token/:kid", mid.Trace(), mid.Logger(cfg.Log), basic, mid.Metrics(), mid.Errors(cfg.Log), mid.Panics(), api.token)
	app.GET("/auth/authenticate", mid.Trace(), mid.Logger(cfg.Log), bearer, mid.Metrics(), mid.Errors(cfg.Log), mid.Panics(), api.authenticate)
	app.POST("/auth/authorize", mid.Trace(), mid.Logger(cfg.Log), mid.Metrics(), mid.Errors(cfg.Log), mid.Panics(), api.authorize)
}
