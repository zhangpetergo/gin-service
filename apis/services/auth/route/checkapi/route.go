package checkapi

import (
	"github.com/zhangpetergo/gin-service/business/api/auth"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, a *auth.Auth) {
	// 不需要中间件的组
	checkGroup := app.Group("")
	checkGroup.GET("/liveness", liveness)
	checkGroup.GET("/readiness", readiness)
}
