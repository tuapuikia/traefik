package middlewares

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
)

// GetLoggerCtx creates a logger context with the middleware fields.
func GetLoggerCtx(ctx context.Context, middleware string, middlewareType string) context.Context {
	return log.With(ctx, log.Str(log.MiddlewareName, middleware), log.Str(log.MiddlewareType, middlewareType))
}
