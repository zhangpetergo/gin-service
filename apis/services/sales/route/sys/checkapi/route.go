package checkapi

import (
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, a *auth.Auth) {

	authen := mid.Authorization(a)
	athAdminOnly := mid.Authorize(a, auth.RuleAdminOnly)

	// 添加组的路
	checkGroup := app.Group("")
	checkGroup.GET("/liveness", liveness)
	checkGroup.GET("/readiness", readiness)
	testGroup := app.Group("")
	testGroup.Use(mid.Trace(), mid.Logger(log), mid.Metrics(), mid.Errors(log), mid.Panics())
	testGroup.GET("/testerror", testError)
	testGroup.GET("/testpanic", testPanic)

	// 测试权限组
	authGroup := app.Group("")
	authGroup.Use(mid.Trace(), mid.Logger(log), authen, athAdminOnly, mid.Metrics(), mid.Errors(log), mid.Panics())
	authGroup.GET("/testauth", liveness)
}
