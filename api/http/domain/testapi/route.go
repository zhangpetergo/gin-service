package testapi

import (
	"github.com/zhangpetergo/gin-service/app/api/auth"
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	authen := mid.Authenticate(cfg.Log, cfg.AuthClient)
	athAdminOnly := mid.Authorize(cfg.Log, cfg.AuthClient, auth.RuleAdminOnly)
	api := newAPI()

	testGroup := app.Group("")
	testGroup.Use(mid.Trace(), mid.Logger(cfg.Log), mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())
	testGroup.GET("/testerror", api.testError)
	testGroup.GET("/testpanic", api.testPanic)

	testAuthGroup := app.Group("")
	testAuthGroup.Use(mid.Trace(), mid.Logger(cfg.Log), authen, athAdminOnly, mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())
	testAuthGroup.GET("/testauth", api.testAuth)
}
