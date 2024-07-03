package checkapi

import (
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App) {
	app.GET("/liveness", liveness)
	app.GET("/readiness", readiness)
	app.GET("/testerror", testError)
	app.GET("/testpanic", testPanic)
}
