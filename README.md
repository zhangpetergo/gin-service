# gin-service







## 1.0 - logger



bill 根据 slog 自定义了一个第三方 logger 包，这个好处是我们到时候可以替换成自己想要的日志

logger 单例，这里我看用到日志的不多，中间件处理 web 请求输出日志，还有 db 日志输出 SQL。

用到地方多的话可以考虑单例

bill 采用的是在 main() 开始初始化`logger`，然后`logger` 作为参数传递给下面想使用`logger`的函数









### main-run 



`run()` 这种想法是如果程序失败了，那么我们可以在`main()`中处理这个错误

```go
package main

import (
	"context"
	"errors"
	"github.com/zhangpetergo/GoLab/slog-example/foundation/logger"
	"math/rand"
	"os"
	"runtime"
)

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return ""
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)

	// -------------------------------------------------------------------------

	ctx := context.Background()
    
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "message", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	defer log.Info(ctx, "shutdown", "GOMAXPROCS", runtime.GOMAXPROCS(0))
    
	if n := rand.Intn(100) % 2; n == 0 {
		return errors.New("ohh bad thing")
	}

	return nil
}

```

### os.Exit(1)

总之，`os.Exit(1)` 用于在程序中出现错误或需要以非正常方式退出时，终止程序并返回一个非零状态码给操作系统。





