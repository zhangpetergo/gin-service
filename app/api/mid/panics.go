package mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/metrics"
	"runtime/debug"
)

func Panics() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		defer func() {
			if rec := recover(); rec != nil {
				trace := debug.Stack()
				err := fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

				metrics.AddPanics(ctx)

				c.Error(err)
				c.Abort()

			}
		}()
		c.Next()
	}
}
