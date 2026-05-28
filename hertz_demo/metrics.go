package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterMetrics(r *server.Hertz) {
	r.GET("/metrics", func(ctx context.Context, c *app.RequestContext) {
		adaptor.HertzHandler(promhttp.Handler())(ctx, c)
	})
}
