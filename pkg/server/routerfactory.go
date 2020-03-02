package server

import (
	"context"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/responsemodifiers"
	"github.com/containous/traefik/v2/pkg/server/middleware"
	"github.com/containous/traefik/v2/pkg/server/router"
	routertcp "github.com/containous/traefik/v2/pkg/server/router/tcp"
	routerudp "github.com/containous/traefik/v2/pkg/server/router/udp"
	"github.com/containous/traefik/v2/pkg/server/service"
	"github.com/containous/traefik/v2/pkg/server/service/tcp"
	"github.com/containous/traefik/v2/pkg/server/service/udp"
	tcpCore "github.com/containous/traefik/v2/pkg/tcp"
	"github.com/containous/traefik/v2/pkg/tls"
	udpCore "github.com/containous/traefik/v2/pkg/udp"
)

// RouterFactory the factory of TCP/UDP routers.
type RouterFactory struct {
	entryPointsTCP []string
	entryPointsUDP []string

	managerFactory *service.ManagerFactory

	chainBuilder *middleware.ChainBuilder
	tlsManager   *tls.Manager
}

// NewRouterFactory creates a new RouterFactory
func NewRouterFactory(staticConfiguration static.Configuration, managerFactory *service.ManagerFactory, tlsManager *tls.Manager, chainBuilder *middleware.ChainBuilder) *RouterFactory {
	var entryPointsTCP, entryPointsUDP []string
	for name, cfg := range staticConfiguration.EntryPoints {
		protocol, err := cfg.GetProtocol()
		if err != nil {
			// Should never happen because Traefik should not start if protocol is invalid.
			log.WithoutContext().Errorf("Invalid protocol: %v", err)
		}

		if protocol == "udp" {
			entryPointsUDP = append(entryPointsUDP, name)
		} else {
			entryPointsTCP = append(entryPointsTCP, name)
		}
	}

	return &RouterFactory{
		entryPointsTCP: entryPointsTCP,
		entryPointsUDP: entryPointsUDP,
		managerFactory: managerFactory,
		tlsManager:     tlsManager,
		chainBuilder:   chainBuilder,
	}
}

// CreateRouters creates new TCPRouters and UDPRouters
func (f *RouterFactory) CreateRouters(conf dynamic.Configuration) (map[string]*tcpCore.Router, map[string]udpCore.Handler) {
	ctx := context.Background()

	rtConf := runtime.NewConfig(conf)

	// HTTP
	serviceManager := f.managerFactory.Build(rtConf)

	middlewaresBuilder := middleware.NewBuilder(rtConf.Middlewares, serviceManager)
	responseModifierFactory := responsemodifiers.NewBuilder(rtConf.Middlewares)

	routerManager := router.NewManager(rtConf, serviceManager, middlewaresBuilder, responseModifierFactory, f.chainBuilder)

	handlersNonTLS := routerManager.BuildHandlers(ctx, f.entryPointsTCP, false)
	handlersTLS := routerManager.BuildHandlers(ctx, f.entryPointsTCP, true)

	serviceManager.LaunchHealthCheck()

	// TCP
	svcTCPManager := tcp.NewManager(rtConf)

	rtTCPManager := routertcp.NewManager(rtConf, svcTCPManager, handlersNonTLS, handlersTLS, f.tlsManager)
	routersTCP := rtTCPManager.BuildHandlers(ctx, f.entryPointsTCP)

	// UDP
	svcUDPManager := udp.NewManager(rtConf)
	rtUDPManager := routerudp.NewManager(rtConf, svcUDPManager)
	routersUDP := rtUDPManager.BuildHandlers(ctx, f.entryPointsUDP)

	rtConf.PopulateUsedBy()

	return routersTCP, routersUDP
}
