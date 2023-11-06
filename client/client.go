package main

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	// 1. 设置全局 trace provider
	otel.SetTracerProvider(trace.NewTracerProvider())

	ctx := context.Background()

	// 2. 创建带参数的日志对象
	lg := log.DefaultLogger
	lg = log.With(lg, "trace_id", tracing.TraceID())
	lg = log.With(lg, "span_id", tracing.SpanID())

	// 3. 创建带中间件的客户端
	c, err := http.NewClient(ctx,
		http.WithEndpoint(":9090"),
		http.WithMiddleware(
			tracing.Client(),
			logging.Client(lg),
		),
	)
	if err != nil {
		panic(err)
	}

	var in emptypb.Empty
	var out emptypb.Empty

	// 4. 发起请求
	err = c.Invoke(ctx, "GET", "/ping", &in, &out)
	if err != nil {
		panic(err)
	}
}
