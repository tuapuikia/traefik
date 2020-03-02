package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/containous/traefik/v2/pkg/config/runtime"
)

type serviceManager interface {
	BuildHTTP(rootCtx context.Context, serviceName string, responseModifier func(*http.Response) error) (http.Handler, error)
	LaunchHealthCheck()
}

// InternalHandlers is the internal HTTP handlers builder.
type InternalHandlers struct {
	api        http.Handler
	dashboard  http.Handler
	rest       http.Handler
	prometheus http.Handler
	ping       http.Handler
	serviceManager
}

// NewInternalHandlers creates a new InternalHandlers.
func NewInternalHandlers(api func(configuration *runtime.Configuration) http.Handler, configuration *runtime.Configuration, rest http.Handler, metricsHandler http.Handler, pingHandler http.Handler, dashboard http.Handler, next serviceManager) *InternalHandlers {
	var apiHandler http.Handler
	if api != nil {
		apiHandler = api(configuration)
	}

	return &InternalHandlers{
		api:            apiHandler,
		dashboard:      dashboard,
		rest:           rest,
		prometheus:     metricsHandler,
		ping:           pingHandler,
		serviceManager: next,
	}
}

// BuildHTTP builds an HTTP handler.
func (m *InternalHandlers) BuildHTTP(rootCtx context.Context, serviceName string, responseModifier func(*http.Response) error) (http.Handler, error) {
	if strings.HasSuffix(serviceName, "@internal") {
		return m.get(serviceName)
	}

	return m.serviceManager.BuildHTTP(rootCtx, serviceName, responseModifier)
}

func (m *InternalHandlers) get(serviceName string) (http.Handler, error) {
	switch serviceName {
	case "api@internal":
		if m.api == nil {
			return nil, errors.New("api is not enabled")
		}
		return m.api, nil

	case "dashboard@internal":
		if m.dashboard == nil {
			return nil, errors.New("dashboard is not enabled")
		}
		return m.dashboard, nil

	case "rest@internal":
		if m.rest == nil {
			return nil, errors.New("rest is not enabled")
		}
		return m.rest, nil

	case "ping@internal":
		if m.ping == nil {
			return nil, errors.New("ping is not enabled")
		}
		return m.ping, nil

	case "prometheus@internal":
		if m.prometheus == nil {
			return nil, errors.New("prometheus is not enabled")
		}
		return m.prometheus, nil

	default:
		return nil, fmt.Errorf("unknown internal service %s", serviceName)
	}
}
