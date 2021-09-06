package proxy

import (
	"context"
	"go-reverse-proxy/app/common/metrics"
	"go-reverse-proxy/app/values"
	"strconv"
	"time"
)

const (
	RequestCount      = "request_count"
	LatencySeconds    = "latency_seconds"
	ForwardMethodName = "Forward"
)

type InstrumentationMiddleware struct {
	Next Handler
	MC   *metrics.MetricsContext
}

func (mw InstrumentationMiddleware) recordRequestCount(ctx context.Context, method string, moduleErr error) {
	lvs := []string{"method", method, "success", strconv.FormatBool(moduleErr == nil)}
	if err := metrics.Record(ctx, RequestCount, 1, lvs...); err != nil {
		mw.MC.Logger.Log("metrics", RequestCount, "method", method, "err", err)
	}
}

func (mw InstrumentationMiddleware) recordLatencySeconds(
	ctx context.Context,
	begin time.Time,
	method string,
	moduleErr error,
) {
	lvs := []string{"method", method, "success", strconv.FormatBool(moduleErr == nil)}
	if err := metrics.RecordTiming(ctx, LatencySeconds, time.Since(begin), lvs...); err != nil {
		mw.MC.Logger.Log("metrics", LatencySeconds, "method", method, "err", err)
	}
}

func (mw InstrumentationMiddleware) Forward(
	initCtx context.Context,
	request *values.Request,
) (
	[]byte,
	int,
	error,
) {
	var err error
	ctx := metrics.IntoContext(initCtx, mw.MC)
	defer func(ctx context.Context, begin time.Time) {
		mw.recordRequestCount(ctx, ForwardMethodName, err)
		mw.recordLatencySeconds(ctx, begin, ForwardMethodName, err)
	}(ctx, time.Now().UTC())

	return mw.Next.Forward(ctx, request)
}
