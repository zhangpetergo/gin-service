package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"time"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 添加 traceID 和 Now 到上下文中
		v := web.Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		c.Request = c.Request.WithContext(web.SetValues(c.Request.Context(), &v))

		c.Next()
	}
}
