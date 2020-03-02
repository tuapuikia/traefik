package replacepath

import (
	"context"
	"net/http"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	// ReplacedPathHeader is the default header to set the old path to.
	ReplacedPathHeader = "X-Replaced-Path"
	typeName           = "ReplacePath"
)

// ReplacePath is a middleware used to replace the path of a URL request.
type replacePath struct {
	next http.Handler
	path string
	name string
}

// New creates a new replace path middleware.
func New(ctx context.Context, next http.Handler, config dynamic.ReplacePath, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating middleware")

	return &replacePath{
		next: next,
		path: config.Path,
		name: name,
	}, nil
}

func (r *replacePath) GetTracingInformation() (string, ext.SpanKindEnum) {
	return r.name, tracing.SpanKindNoneEnum
}

func (r *replacePath) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.Header.Add(ReplacedPathHeader, req.URL.Path)
	req.URL.Path = r.path
	req.RequestURI = req.URL.RequestURI()

	r.next.ServeHTTP(rw, req)
}
