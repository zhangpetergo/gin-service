package mid

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"github.com/zhangpetergo/gin-service/foundation/web"
	"net/http"
)

var codeStatus [17]int

// init maps out the error codes to http status codes.
func init() {
	codeStatus[errs.OK.Value()] = http.StatusOK
	codeStatus[errs.Canceled.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.Unknown.Value()] = http.StatusInternalServerError
	codeStatus[errs.InvalidArgument.Value()] = http.StatusBadRequest
	codeStatus[errs.DeadlineExceeded.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.NotFound.Value()] = http.StatusNotFound
	codeStatus[errs.AlreadyExists.Value()] = http.StatusConflict
	codeStatus[errs.PermissionDenied.Value()] = http.StatusForbidden
	codeStatus[errs.ResourceExhausted.Value()] = http.StatusTooManyRequests
	codeStatus[errs.FailedPrecondition.Value()] = http.StatusBadRequest
	codeStatus[errs.Aborted.Value()] = http.StatusConflict
	codeStatus[errs.OutOfRange.Value()] = http.StatusBadRequest
	codeStatus[errs.Unimplemented.Value()] = http.StatusNotImplemented
	codeStatus[errs.Internal.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unavailable.Value()] = http.StatusServiceUnavailable
	codeStatus[errs.DataLoss.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unauthenticated.Value()] = http.StatusUnauthorized
}

func Errors(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ctx := c.Request.Context()
		// 程序处理完毕
		// 判断是否存在错误
		if len(c.Errors) > 0 {
			// 处理第一个错误
			// 在 gin 中，错误是一个数组，这里只处理第一个错误，一般来说我们在程序中遇到错误时，只会返回一个错误
			// 如果出现了例外情况，那么我们需要修改这里的代码
			err := c.Errors[0].Err
			// 记录错误
			log.Error(ctx, "message", "ERROR", err.Error())
			// 清空 c.Errors
			// c.Errors = []*gin.Error{}

			// 返回的是我们自定义的错误
			if errs.IsError(err) {
				var e errs.Error
				errors.As(err, &e)
				// 修改返回内容
				c.JSON(codeStatus[e.Code.Value()], e)
				return
			}
			// 返回的不是我们自定义的错误，返回未知错误
			e := errs.Newf(errs.Unknown, errs.Unknown.String())
			c.JSON(codeStatus[e.Code.Value()], e)

			// 如果是关闭服务的错误，直接返回
			if web.IsShutdown(err) {
				c.Error(err)
			}
		}

	}
}
