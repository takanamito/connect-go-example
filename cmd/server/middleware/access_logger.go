package middleware

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

func LoggingInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			startTime := time.Now()

			res, err := next(ctx, req)

			duration := time.Since(startTime)

			group := slog.Group("req",
				slog.String("protocol", req.Peer().Protocol),
				slog.String("method", req.HTTPMethod()),
				slog.String("procedure", req.Spec().Procedure),
				slog.String("query", req.Peer().Query.Encode()),
				slog.String("duration", duration.String()),
			)
			logger.Info("request", group)

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
