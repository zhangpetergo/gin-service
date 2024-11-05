package checkapi

import (
	"github.com/jmoiron/sqlx"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(build string, app *web.App, log *logger.Logger, db *sqlx.DB) {
	// 不需要中间件的组

	api := newAPI(build, log, db)

	checkGroup := app.Group("")
	checkGroup.GET("/liveness", api.liveness)
	checkGroup.GET("/readiness", api.readiness)
}
