package crd

import (
	"context"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/provider"
	"github.com/containous/traefik/v2/pkg/tls"
	"github.com/stretchr/testify/assert"
)

var _ provider.Provider = (*Provider)(nil)

func Int(v int) *int    { return &v }
func Bool(v bool) *bool { return &v }

func TestLoadIngressRouteTCPs(t *testing.T) {
	testCases := []struct {
		desc         string
		ingressClass string
		paths        []string
		expected     *dynamic.Configuration
	}{
		{
			desc: "Empty",
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "Simple Ingress Route, with foo entrypoint",
			paths: []string{"tcp/services.yml", "tcp/simple.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with two different rules",
			paths: []string{"tcp/services.yml", "tcp/with_two_rules.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
						},
						"default-test.route-f44ce589164e656d231c": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-f44ce589164e656d231c",
							Rule:        "HostSNI(`bar.com`)",
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
						"default-test.route-f44ce589164e656d231c": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with two identical rules",
			paths: []string{"tcp/services.yml", "tcp/with_two_identical_rules.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with two different services",
			paths: []string{"tcp/services.yml", "tcp/with_two_services.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							Weighted: &dynamic.TCPWeightedRoundRobin{
								Services: []dynamic.TCPWRRService{
									{
										Name:   "default-test.route-fdd3e9338e47a45efefc-whoamitcp-8000",
										Weight: func(i int) *int { return &i }(2),
									},
									{
										Name:   "default-test.route-fdd3e9338e47a45efefc-whoamitcp2-8080",
										Weight: func(i int) *int { return &i }(3),
									},
								},
							},
						},
						"default-test.route-fdd3e9338e47a45efefc-whoamitcp-8000": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
									},
									{
										Address: "10.10.0.2:8000",
									},
								},
							},
						},
						"default-test.route-fdd3e9338e47a45efefc-whoamitcp2-8080": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.3:8080",
									},
									{
										Address: "10.10.0.4:8080",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with different services namespaces",
			paths: []string{"tcp/services.yml", "tcp/with_different_services_ns.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							Weighted: &dynamic.TCPWeightedRoundRobin{
								Services: []dynamic.TCPWRRService{
									{
										Name:   "default-test.route-fdd3e9338e47a45efefc-whoamitcp-8000",
										Weight: func(i int) *int { return &i }(2),
									},
									{
										Name:   "default-test.route-fdd3e9338e47a45efefc-whoamitcp2-8080",
										Weight: func(i int) *int { return &i }(3),
									},
									{
										Name:   "default-test.route-fdd3e9338e47a45efefc-whoamitcp3-8083",
										Weight: func(i int) *int { return &i }(4),
									},
								},
							},
						},
						"default-test.route-fdd3e9338e47a45efefc-whoamitcp-8000": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
									},
									{
										Address: "10.10.0.2:8000",
									},
								},
							},
						},
						"default-test.route-fdd3e9338e47a45efefc-whoamitcp2-8080": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.3:8080",
									},
									{
										Address: "10.10.0.4:8080",
									},
								},
							},
						},
						"default-test.route-fdd3e9338e47a45efefc-whoamitcp3-8083": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.7:8083",
									},
									{
										Address: "10.10.0.8:8083",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:         "Ingress class does not match",
			paths:        []string{"tcp/services.yml", "tcp/simple.yml"},
			ingressClass: "tchouk",
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "Route with empty rule value is ignored",
			paths: []string{"tcp/services.yml", "tcp/with_no_rule_value.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "check rule quoting validity",
			paths: []string{"tcp/services.yml", "tcp/with_bad_host_rule.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "TLS",
			paths: []string{"tcp/services.yml", "tcp/with_tls.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Certificates: []*tls.CertAndStores{
						{
							Certificate: tls.Certificate{
								CertFile: tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								KeyFile:  tls.FileOrContent("-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----"),
							},
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS:         &dynamic.RouterTCPTLSConfig{},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS with passthrough",
			paths: []string{"tcp/services.yml", "tcp/with_tls_passthrough.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS: &dynamic.RouterTCPTLSConfig{
								Passthrough: true,
							},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "TLS with tls options",
			paths: []string{"tcp/services.yml", "tcp/with_tls_options.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
							CipherSuites: []string{
								"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
								"TLS_RSA_WITH_AES_256_GCM_SHA384",
							},
							ClientAuth: tls.ClientAuth{
								CAFiles: []tls.FileOrContent{
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								},
								ClientAuthType: "VerifyClientCertIfGiven",
							},
							SniStrict:                true,
							PreferServerCipherSuites: true,
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS: &dynamic.RouterTCPTLSConfig{
								Options: "default-foo",
							},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS with tls options and specific namespace",
			paths: []string{"tcp/services.yml", "tcp/with_tls_options_and_specific_namespace.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"myns-foo": {
							MinVersion: "VersionTLS12",
							CipherSuites: []string{
								"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
								"TLS_RSA_WITH_AES_256_GCM_SHA384",
							},
							ClientAuth: tls.ClientAuth{
								CAFiles: []tls.FileOrContent{
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								},
								ClientAuthType: "VerifyClientCertIfGiven",
							},
							SniStrict: true,
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS: &dynamic.RouterTCPTLSConfig{
								Options: "myns-foo",
							},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS with bad tls options",
			paths: []string{"tcp/services.yml", "tcp/with_bad_tls_options.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
							CipherSuites: []string{
								"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
								"TLS_RSA_WITH_AES_256_GCM_SHA384",
							},
							ClientAuth: tls.ClientAuth{
								CAFiles: []tls.FileOrContent{
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								},
								ClientAuthType: "VerifyClientCertIfGiven",
							},
							SniStrict: true,
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS: &dynamic.RouterTCPTLSConfig{
								Options: "default-foo",
							},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS with unknown tls options",
			paths: []string{"tcp/services.yml", "tcp/with_unknown_tls_options.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS: &dynamic.RouterTCPTLSConfig{
								Options: "default-unknown",
							},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS with unknown tls options namespace",
			paths: []string{"tcp/services.yml", "tcp/with_unknown_tls_options_namespace.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS: &dynamic.RouterTCPTLSConfig{
								Options: "unknown-foo",
							},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS with ACME",
			paths: []string{"tcp/services.yml", "tcp/with_tls_acme.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS:         &dynamic.RouterTCPTLSConfig{},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "TCP with terminationDelay",
			paths: []string{"tcp/services.yml", "tcp/with_termination_delay.yml"},
			expected: &dynamic.Configuration{
				TLS: &dynamic.TLSConfiguration{},
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
								TerminationDelay: Int(500),
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS with tls Store",
			paths: []string{"tcp/services.yml", "tcp/with_tls_store.yml"},
			expected: &dynamic.Configuration{
				TLS: &dynamic.TLSConfiguration{
					Stores: map[string]tls.Store{
						"default": {
							DefaultCertificate: &tls.Certificate{
								CertFile: tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								KeyFile:  tls.FileOrContent("-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----"),
							},
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers: map[string]*dynamic.TCPRouter{
						"default-test.route-fdd3e9338e47a45efefc": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-fdd3e9338e47a45efefc",
							Rule:        "HostSNI(`foo.com`)",
							TLS:         &dynamic.RouterTCPTLSConfig{},
						},
					},
					Services: map[string]*dynamic.TCPService{
						"default-test.route-fdd3e9338e47a45efefc": {
							LoadBalancer: &dynamic.TCPServersLoadBalancer{
								Servers: []dynamic.TCPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			if test.expected == nil {
				return
			}

			p := Provider{IngressClass: test.ingressClass}
			conf := p.loadConfigurationFromCRD(context.Background(), newClientMock(test.paths...))
			assert.Equal(t, test.expected, conf)
		})
	}
}

func TestLoadIngressRoutes(t *testing.T) {
	testCases := []struct {
		desc         string
		ingressClass string
		paths        []string
		expected     *dynamic.Configuration
	}{
		{
			desc: "Empty",
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "Simple Ingress Route, with foo entrypoint",
			paths: []string{"services.yml", "simple.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"foo"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "Simple Ingress Route with middleware",
			paths: []string{"services.yml", "with_middleware.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test2-route-23c7f4c450289ee29016": {
							EntryPoints: []string{"web"},
							Service:     "default-test2-route-23c7f4c450289ee29016",
							Rule:        "Host(`foo.com`) && PathPrefix(`/tobestripped`)",
							Priority:    12,
							Middlewares: []string{"default-stripprefix", "foo-addprefix"},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{
						"default-stripprefix": {
							StripPrefix: &dynamic.StripPrefix{
								Prefixes: []string{"/tobestripped"},
							},
						},
						"foo-addprefix": {
							AddPrefix: &dynamic.AddPrefix{
								Prefix: "/tobeadded",
							},
						},
					},
					Services: map[string]*dynamic.Service{
						"default-test2-route-23c7f4c450289ee29016": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "Simple Ingress Route with middleware crossprovider",
			paths: []string{"services.yml", "with_middleware_crossprovider.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test2-route-23c7f4c450289ee29016": {
							EntryPoints: []string{"web"},
							Service:     "default-test2-route-23c7f4c450289ee29016",
							Rule:        "Host(`foo.com`) && PathPrefix(`/tobestripped`)",
							Priority:    12,
							Middlewares: []string{"default-stripprefix", "foo-addprefix", "basicauth@file", "redirect@file"},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{
						"default-stripprefix": {
							StripPrefix: &dynamic.StripPrefix{
								Prefixes: []string{"/tobestripped"},
							},
						},
						"foo-addprefix": {
							AddPrefix: &dynamic.AddPrefix{
								Prefix: "/tobeadded",
							},
						},
					},
					Services: map[string]*dynamic.Service{
						"default-test2-route-23c7f4c450289ee29016": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "One ingress Route with two different rules",
			paths: []string{"services.yml", "with_two_rules.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Service:     "default-test-route-6b204d94623b3df4370c",
							Priority:    14,
						},
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-77c62dfe9517144aeeaa",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-test-route-77c62dfe9517144aeeaa": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with two different services",
			paths: []string{"services.yml", "with_two_services.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-77c62dfe9517144aeeaa",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-77c62dfe9517144aeeaa": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami-80",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "default-whoami2-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-whoami-80": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-whoami2-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "one kube service (== servers lb) in a services wrr",
			paths: []string{"with_services_lb0.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-wrr1",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-wrr1": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "traefik service without ingress route",
			paths: []string{"with_services_only.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-wrr1": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "One ingress Route with two different services, each with two services, balancing servers nested",
			paths: []string{"with_services_lb1.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-77c62dfe9517144aeeaa",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-77c62dfe9517144aeeaa": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-wrr1",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "default-wrr2",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-wrr1": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami4-80",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-whoami4-80": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-wrr2": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami6-80",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "default-whoami7-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-whoami6-80": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.5:80",
									},
									{
										URL: "http://10.10.0.6:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-whoami7-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.7:8080",
									},
									{
										URL: "http://10.10.0.8:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "one wrr and one kube service (== servers lb) in a wrr",
			paths: []string{"with_services_lb2.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-wrr1",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-wrr1": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-wrr2",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-wrr2": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "services lb, servers lb, and mirror service, all in a wrr",
			paths: []string{"with_services_lb3.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-wrr1",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-wrr1": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-wrr2",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "default-mirror1",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-wrr2": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-mirror1": {
							Mirroring: &dynamic.Mirroring{
								Service: "default-whoami5-8080",
								Mirrors: []dynamic.MirrorService{
									{Name: "default-whoami4-8080", Percent: 50},
								},
							},
						},
						"default-whoami4-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:8080",
									},
									{
										URL: "http://10.10.0.2:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "services lb, servers lb, and mirror service, all in a wrr with different namespaces",
			paths: []string{"with_namespaces.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-77c62dfe9517144aeeaa",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-77c62dfe9517144aeeaa": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "baz-whoami6-8080",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "foo-wrr1",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "foo-mirror2",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "foo-mirror3",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "foo-mirror4",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"baz-whoami6-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.5:8080",
									},
									{
										URL: "http://10.10.0.6:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"foo-wrr1": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "foo-whoami4-8080",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "baz-whoami6-8080",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "foo-mirror1",
										Weight: func(i int) *int { return &i }(1),
									},
									{
										Name:   "bar-wrr2",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"foo-whoami4-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:8080",
									},
									{
										URL: "http://10.10.0.2:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"foo-mirror1": {
							Mirroring: &dynamic.Mirroring{
								Service: "foo-whoami5-8080",
								Mirrors: []dynamic.MirrorService{
									{Name: "foo-whoami4-8080"},
									{Name: "baz-whoami6-8080"},
									{Name: "bar-mirrored"},
								},
							},
						},
						"foo-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"bar-mirrored": {
							Mirroring: &dynamic.Mirroring{
								Service: "baz-whoami6-8080",
								Mirrors: []dynamic.MirrorService{
									{Name: "foo-whoami4-8080", Percent: 50},
								},
							},
						},
						"foo-mirror2": {
							Mirroring: &dynamic.Mirroring{
								Service: "foo-whoami5-8080",
								Mirrors: []dynamic.MirrorService{
									{Name: "foo-whoami4-8080"},
									{Name: "baz-whoami6-8080"},
									{Name: "bar-mirrored"},
									{Name: "foo-wrr1"},
								},
							},
						},
						"foo-mirror3": {
							Mirroring: &dynamic.Mirroring{
								Service: "foo-wrr1",
								Mirrors: []dynamic.MirrorService{
									{Name: "foo-whoami4-8080"},
									{Name: "baz-whoami6-8080"},
									{Name: "bar-mirrored"},
									{Name: "foo-wrr1"},
								},
							},
						},
						"bar-wrr2": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "foo-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"foo-mirror4": {
							Mirroring: &dynamic.Mirroring{
								Service: "foo-wrr1",
								Mirrors: []dynamic.MirrorService{
									{Name: "foo-whoami4-8080"},
									{Name: "baz-whoami6-8080"},
									{Name: "bar-mirrored"},
									{Name: "foo-wrr1"},
								},
							},
						},
					},
				},
			},
		},
		{
			desc:  "one kube service (== servers lb) in a mirroring",
			paths: []string{"with_mirroring.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-mirror1",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-mirror1": {
							Mirroring: &dynamic.Mirroring{
								Service: "default-whoami5-8080",
								Mirrors: []dynamic.MirrorService{
									{Name: "default-whoami4-8080", Percent: 50},
								},
							},
						},
						"default-whoami4-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:8080",
									},
									{
										URL: "http://10.10.0.2:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "weighted services in a mirroring",
			paths: []string{"with_mirroring2.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-mirror1",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-mirror1": {
							Mirroring: &dynamic.Mirroring{
								Service: "default-wrr1",
								Mirrors: []dynamic.MirrorService{
									{Name: "default-wrr2", Percent: 30},
								},
							},
						},
						"default-wrr1": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami4-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-wrr2": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami5-8080",
										Weight: func(i int) *int { return &i }(1),
									},
								},
							},
						},
						"default-whoami4-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:8080",
									},
									{
										URL: "http://10.10.0.2:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-whoami5-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "One ingress Route with two different services, with weights",
			paths: []string{"services.yml", "with_two_services_weight.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-77c62dfe9517144aeeaa": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-77c62dfe9517144aeeaa",
							Rule:        "Host(`foo.com`) && PathPrefix(`/foo`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-77c62dfe9517144aeeaa": {
							Weighted: &dynamic.WeightedRoundRobin{
								Services: []dynamic.WRRService{
									{
										Name:   "default-whoami-80",
										Weight: Int(10),
									},
									{
										Name:   "default-whoami2-8080",
										Weight: Int(0),
									},
								},
							},
						},
						"default-whoami-80": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
						"default-whoami2-8080": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.3:8080",
									},
									{
										URL: "http://10.10.0.4:8080",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:         "Ingress class",
			paths:        []string{"services.yml", "simple.yml"},
			ingressClass: "tchouk",
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "Route with empty rule value is ignored",
			paths: []string{"services.yml", "with_no_rule_value.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "Route with kind not of a rule type (empty kind) is ignored",
			paths: []string{"services.yml", "with_wrong_rule_kind.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "check rule quoting validity",
			paths: []string{"services.yml", "with_bad_host_rule.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "TLS",
			paths: []string{"services.yml", "with_tls.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Certificates: []*tls.CertAndStores{
						{
							Certificate: tls.Certificate{
								CertFile: tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								KeyFile:  tls.FileOrContent("-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----"),
							},
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS:         &dynamic.RouterTLSConfig{},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with tls options",
			paths: []string{"services.yml", "with_tls_options.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
							CipherSuites: []string{
								"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
								"TLS_RSA_WITH_AES_256_GCM_SHA384",
							},
							ClientAuth: tls.ClientAuth{
								CAFiles: []tls.FileOrContent{
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								},
								ClientAuthType: "VerifyClientCertIfGiven",
							},
							SniStrict:                true,
							PreferServerCipherSuites: true,
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS: &dynamic.RouterTLSConfig{
								Options: "default-foo",
							},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with two default tls options",
			paths: []string{"services.yml", "with_default_tls_options.yml", "with_default_tls_options_default_namespace.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS: &dynamic.RouterTLSConfig{
								Options: "default-foo",
							},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with default tls options",
			paths: []string{"services.yml", "with_default_tls_options.yml"},
			expected: &dynamic.Configuration{
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default": {
							MinVersion: "VersionTLS12",
							CipherSuites: []string{
								"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
								"TLS_RSA_WITH_AES_256_GCM_SHA384",
							},
							ClientAuth: tls.ClientAuth{
								CAFiles: []tls.FileOrContent{
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								},
								ClientAuthType: "VerifyClientCertIfGiven",
							},
							SniStrict:                true,
							PreferServerCipherSuites: true,
						},
					},
				},
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS: &dynamic.RouterTLSConfig{
								Options: "default-foo",
							},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with tls options and specific namespace",
			paths: []string{"services.yml", "with_tls_options_and_specific_namespace.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"myns-foo": {
							MinVersion: "VersionTLS12",
							CipherSuites: []string{
								"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
								"TLS_RSA_WITH_AES_256_GCM_SHA384",
							},
							ClientAuth: tls.ClientAuth{
								CAFiles: []tls.FileOrContent{
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								},
								ClientAuthType: "VerifyClientCertIfGiven",
							},
							SniStrict: true,
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS: &dynamic.RouterTLSConfig{
								Options: "myns-foo",
							},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with bad tls options",
			paths: []string{"services.yml", "with_bad_tls_options.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
							CipherSuites: []string{
								"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
								"TLS_RSA_WITH_AES_256_GCM_SHA384",
							},
							ClientAuth: tls.ClientAuth{
								CAFiles: []tls.FileOrContent{
									tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								},
								ClientAuthType: "VerifyClientCertIfGiven",
							},
							SniStrict: true,
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS: &dynamic.RouterTLSConfig{
								Options: "default-foo",
							},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with unknown tls options",
			paths: []string{"services.yml", "with_unknown_tls_options.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS: &dynamic.RouterTLSConfig{
								Options: "default-unknown",
							},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with unknown tls options namespace",
			paths: []string{"services.yml", "with_unknown_tls_options_namespace.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{
					Options: map[string]tls.Options{
						"default-foo": {
							MinVersion: "VersionTLS12",
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS: &dynamic.RouterTLSConfig{
								Options: "unknown-foo",
							},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with ACME",
			paths: []string{"services.yml", "with_tls_acme.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS:         &dynamic.RouterTLSConfig{},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "Simple Ingress Route, defaulting to https for servers",
			paths: []string{"services.yml", "with_https_default.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"foo"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "https://10.10.0.5:8443",
									},
									{
										URL: "https://10.10.0.6:8443",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "Simple Ingress Route, explicit https scheme",
			paths: []string{"services.yml", "with_https_scheme.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"foo"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "https://10.10.0.7:8443",
									},
									{
										URL: "https://10.10.0.8:8443",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "Simple Ingress Route, with basic auth middleware",
			paths: []string{"services.yml", "with_auth.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{
						"default-basicauth": {
							BasicAuth: &dynamic.BasicAuth{
								Users: dynamic.Users{"test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/", "test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0"},
							},
						},
						"default-digestauth": {
							DigestAuth: &dynamic.DigestAuth{
								Users: dynamic.Users{"test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/", "test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0"},
							},
						},
						"default-forwardauth": {
							ForwardAuth: &dynamic.ForwardAuth{
								Address: "test.com",
								TLS: &dynamic.ClientTLS{
									CA:   "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----",
									Cert: "-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----",
									Key:  "-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----",
								},
							},
						},
					},
					Services: map[string]*dynamic.Service{},
				},
			},
		},
		{
			desc:  "Simple Ingress Route, with error page middleware",
			paths: []string{"services.yml", "with_error_page.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{
						"default-errorpage": {
							Errors: &dynamic.ErrorPage{
								Status:  []string{"404", "500"},
								Service: "default-errorpage-errorpage-service",
								Query:   "query",
							},
						},
					},
					Services: map[string]*dynamic.Service{
						"default-errorpage-errorpage-service": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "Simple Ingress Route, with options",
			paths: []string{"services.yml", "with_options.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"foo"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader:     Bool(false),
								ResponseForwarding: &dynamic.ResponseForwarding{FlushInterval: "10s"},
							},
						},
					},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "TLS with tls store",
			paths: []string{"services.yml", "with_tls_store.yml"},
			expected: &dynamic.Configuration{
				TLS: &dynamic.TLSConfiguration{
					Stores: map[string]tls.Store{
						"default": {
							DefaultCertificate: &tls.Certificate{
								CertFile: tls.FileOrContent("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----"),
								KeyFile:  tls.FileOrContent("-----BEGIN PRIVATE KEY-----\n-----END PRIVATE KEY-----"),
							},
						},
					},
				},
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS:         &dynamic.RouterTLSConfig{},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc:  "TLS with tls store default two times",
			paths: []string{"services.yml", "with_tls_store.yml", "with_default_tls_store.yml"},
			expected: &dynamic.Configuration{
				TLS: &dynamic.TLSConfiguration{
					Stores: map[string]tls.Store{},
				},
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers: map[string]*dynamic.Router{
						"default-test-route-6b204d94623b3df4370c": {
							EntryPoints: []string{"web"},
							Service:     "default-test-route-6b204d94623b3df4370c",
							Rule:        "Host(`foo.com`) && PathPrefix(`/bar`)",
							Priority:    12,
							TLS:         &dynamic.RouterTLSConfig{},
						},
					},
					Middlewares: map[string]*dynamic.Middleware{},
					Services: map[string]*dynamic.Service{
						"default-test-route-6b204d94623b3df4370c": {
							LoadBalancer: &dynamic.ServersLoadBalancer{
								Servers: []dynamic.Server{
									{
										URL: "http://10.10.0.1:80",
									},
									{
										URL: "http://10.10.0.2:80",
									},
								},
								PassHostHeader: Bool(true),
							},
						},
					},
				},
			},
		},
		{
			desc: "port selected by name (TODO)",
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			if test.expected == nil {
				return
			}

			p := Provider{IngressClass: test.ingressClass}
			conf := p.loadConfigurationFromCRD(context.Background(), newClientMock(test.paths...))
			assert.Equal(t, test.expected, conf)
		})
	}
}

func TestLoadIngressRouteUDPs(t *testing.T) {
	testCases := []struct {
		desc         string
		ingressClass string
		paths        []string
		expected     *dynamic.Configuration
	}{
		{
			desc: "Empty",
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "Simple Ingress Route, with foo entrypoint",
			paths: []string{"udp/services.yml", "udp/simple.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers: map[string]*dynamic.UDPRouter{
						"default-test.route-0": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-0",
						},
					},
					Services: map[string]*dynamic.UDPService{
						"default-test.route-0": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
					},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with two different routes",
			paths: []string{"udp/services.yml", "udp/with_two_routes.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers: map[string]*dynamic.UDPRouter{
						"default-test.route-0": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-0",
						},
						"default-test.route-1": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-1",
						},
					},
					Services: map[string]*dynamic.UDPService{
						"default-test.route-0": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.1:8000",
										Port:    "",
									},
									{
										Address: "10.10.0.2:8000",
										Port:    "",
									},
								},
							},
						},
						"default-test.route-1": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.3:8080",
										Port:    "",
									},
									{
										Address: "10.10.0.4:8080",
										Port:    "",
									},
								},
							},
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with two different services",
			paths: []string{"udp/services.yml", "udp/with_two_services.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers: map[string]*dynamic.UDPRouter{
						"default-test.route-0": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-0",
						},
					},
					Services: map[string]*dynamic.UDPService{
						"default-test.route-0": {
							Weighted: &dynamic.UDPWeightedRoundRobin{
								Services: []dynamic.UDPWRRService{
									{
										Name:   "default-test.route-0-whoamiudp-8000",
										Weight: func(i int) *int { return &i }(2),
									},
									{
										Name:   "default-test.route-0-whoamiudp2-8080",
										Weight: func(i int) *int { return &i }(3),
									},
								},
							},
						},
						"default-test.route-0-whoamiudp-8000": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.1:8000",
									},
									{
										Address: "10.10.0.2:8000",
									},
								},
							},
						},
						"default-test.route-0-whoamiudp2-8080": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.3:8080",
									},
									{
										Address: "10.10.0.4:8080",
									},
								},
							},
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:  "One ingress Route with different services namespaces",
			paths: []string{"udp/services.yml", "udp/with_different_services_ns.yml"},
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers: map[string]*dynamic.UDPRouter{
						"default-test.route-0": {
							EntryPoints: []string{"foo"},
							Service:     "default-test.route-0",
						},
					},
					Services: map[string]*dynamic.UDPService{
						"default-test.route-0": {
							Weighted: &dynamic.UDPWeightedRoundRobin{
								Services: []dynamic.UDPWRRService{
									{
										Name:   "default-test.route-0-whoamiudp-8000",
										Weight: func(i int) *int { return &i }(2),
									},
									{
										Name:   "default-test.route-0-whoamiudp2-8080",
										Weight: func(i int) *int { return &i }(3),
									},
									{
										Name:   "default-test.route-0-whoamiudp3-8083",
										Weight: func(i int) *int { return &i }(4),
									},
								},
							},
						},
						"default-test.route-0-whoamiudp-8000": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.1:8000",
									},
									{
										Address: "10.10.0.2:8000",
									},
								},
							},
						},
						"default-test.route-0-whoamiudp2-8080": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.3:8080",
									},
									{
										Address: "10.10.0.4:8080",
									},
								},
							},
						},
						"default-test.route-0-whoamiudp3-8083": {
							LoadBalancer: &dynamic.UDPServersLoadBalancer{
								Servers: []dynamic.UDPServer{
									{
										Address: "10.10.0.7:8083",
									},
									{
										Address: "10.10.0.8:8083",
									},
								},
							},
						},
					},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
		{
			desc:         "Ingress class does not match",
			paths:        []string{"udp/services.yml", "udp/simple.yml"},
			ingressClass: "tchouk",
			expected: &dynamic.Configuration{
				UDP: &dynamic.UDPConfiguration{
					Routers:  map[string]*dynamic.UDPRouter{},
					Services: map[string]*dynamic.UDPService{},
				},
				TCP: &dynamic.TCPConfiguration{
					Routers:  map[string]*dynamic.TCPRouter{},
					Services: map[string]*dynamic.TCPService{},
				},
				HTTP: &dynamic.HTTPConfiguration{
					Routers:     map[string]*dynamic.Router{},
					Middlewares: map[string]*dynamic.Middleware{},
					Services:    map[string]*dynamic.Service{},
				},
				TLS: &dynamic.TLSConfiguration{},
			},
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			if test.expected == nil {
				return
			}

			p := Provider{IngressClass: test.ingressClass}
			conf := p.loadConfigurationFromCRD(context.Background(), newClientMock(test.paths...))
			assert.Equal(t, test.expected, conf)
		})
	}
}

func TestParseServiceProtocol(t *testing.T) {
	testCases := []struct {
		desc          string
		scheme        string
		portName      string
		portNumber    int32
		expected      string
		expectedError bool
	}{
		{
			desc:       "Empty scheme and name",
			scheme:     "",
			portName:   "",
			portNumber: 1000,
			expected:   "http",
		},
		{
			desc:       "h2c scheme and emptyname",
			scheme:     "h2c",
			portName:   "",
			portNumber: 1000,
			expected:   "h2c",
		},
		{
			desc:          "invalid scheme",
			scheme:        "foo",
			portName:      "",
			portNumber:    1000,
			expectedError: true,
		},
		{
			desc:       "Empty scheme and https name",
			scheme:     "",
			portName:   "https-secure",
			portNumber: 1000,
			expected:   "https",
		},
		{
			desc:       "Empty scheme and port number",
			scheme:     "",
			portName:   "",
			portNumber: 443,
			expected:   "https",
		},
		{
			desc:       "https scheme",
			scheme:     "https",
			portName:   "",
			portNumber: 1000,
			expected:   "https",
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			protocol, err := parseServiceProtocol(test.scheme, test.portName, test.portNumber)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expected, protocol)
			}
		})
	}
}
