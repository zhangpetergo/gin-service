package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"time"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
		r := c.Request

		// 从上下文中获取 traceID
		v := web.GetValues(c)

		// 程序运行之前打印
		log.Info(ctx, "request started", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

		c.Next()

		statusCode := c.Writer.Status()

		log.Info(ctx, "request completed", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr,
			"statuscode", statusCode, "since", time.Since(v.Now).String())

	}
}
