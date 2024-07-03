package checkapi

import (
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger) {
	// 添加组的路
	checkGroup := app.Group("")
	checkGroup.GET("/liveness", liveness)
	checkGroup.GET("/readiness", readiness)
	testGroup := app.Group("")
	testGroup.Use(mid.Trace(), mid.Logger(log), mid.Metrics(), mid.Errors(log), mid.Panics())
	testGroup.GET("/testerror", testError)
	testGroup.GET("/testpanic", testPanic)
}
