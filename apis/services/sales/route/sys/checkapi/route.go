package checkapi

import (
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, authClient *authclient.Client) {

	authen := mid.AuthenticateService(log, authClient)
	athAdminOnly := mid.AuthorizeService(log, authClient, auth.RuleAdminOnly)

	// 不需要中间件的组
	checkGroup := app.Group("")
	checkGroup.GET("/liveness", liveness)
	checkGroup.GET("/readiness", readiness)

	// 测试中间件的组
	testGroup := app.Group("")
	testGroup.Use(mid.Trace(), mid.Logger(log), mid.Metrics(), mid.Errors(log), mid.Panics())
	testGroup.GET("/testerror", testError)
	testGroup.GET("/testpanic", testPanic)

	// 测试权限组
	authGroup := app.Group("")
	authGroup.Use(mid.Trace(), mid.Logger(log), authen, athAdminOnly, mid.Metrics(), mid.Errors(log), mid.Panics())
	authGroup.GET("/testauth", liveness)
}
