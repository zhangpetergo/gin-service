package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/metrics"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 统计请求个数
		ctx := metrics.Set(c.Request.Context())
		// 为了在请求处理过程中传递 ctx，必须这样做
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		n := metrics.AddRequests(ctx)

		if n%1000 == 0 {
			metrics.AddGoroutines(ctx)
		}

		if len(c.Errors) > 0 {
			metrics.AddErrors(ctx)
		}
	}
}
