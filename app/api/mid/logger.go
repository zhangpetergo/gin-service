package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/foundation/logger"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		r := c.Request
		// 程序运行之前打印
		log.Info(ctx, "request started", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
		c.Next()

		log.Info(ctx, "request completed", "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

		// 程序运行之后打印
	}
}
