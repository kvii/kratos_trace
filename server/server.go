package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
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

	// 2. 创建带参数的日志对象
	lg := log.DefaultLogger
	lg = log.With(lg, "trace_id", tracing.TraceID())
	lg = log.With(lg, "span_id", tracing.SpanID())

	// 3. 创建带中间件的服务
	srv := http.NewServer(
		http.Address(":9090"),
		http.Middleware(
			tracing.Server(),
			logging.Server(lg),
		),
	)
	srv.Route("/").GET("/ping", func(ctx http.Context) error {
		// 4. 写路由逻辑
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			var out emptypb.Empty
			return &out, nil
		})

		var in emptypb.Empty
		err := ctx.Bind(&in)
		if err != nil {
			return err
		}
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		return ctx.Result(200, out)
	})

	s := kratos.New(
		kratos.Server(srv),
	)
	if err := s.Run(); err != nil {
		panic(err)
	}
}
