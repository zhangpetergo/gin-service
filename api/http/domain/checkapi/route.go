package checkapi

import (
	"github.com/jmoiron/sqlx"
	"github.com/zhangpetergo/gin-service/app/api/mid"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	api := newAPI(cfg.Build, cfg.Log, cfg.DB)

	checkGroup := app.Group("", mid.Trace(), mid.Logger(cfg.Log), mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())

	checkGroup.GET("/liveness", api.liveness)
	checkGroup.GET("/readiness", api.readiness)
}
