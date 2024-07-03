package mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func Panics() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				trace := debug.Stack()
				err := fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				c.Error(err)
				c.Abort()
			}
		}()
		c.Next()
	}
}
