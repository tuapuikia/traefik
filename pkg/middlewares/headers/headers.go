// Package headers Middleware based on https://github.com/unrolled/secure.
package headers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/unrolled/secure"
)

const (
	typeName = "Headers"
)

type headers struct {
	name    string
	handler http.Handler
}

// New creates a Headers middleware.
func New(ctx context.Context, next http.Handler, config dynamic.Headers, name string) (http.Handler, error) {
	// HeaderMiddleware -> SecureMiddleWare -> next
	logger := log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName))
	logger.Debug("Creating middleware")

	hasSecureHeaders := config.HasSecureHeadersDefined()
	hasCustomHeaders := config.HasCustomHeadersDefined()
	hasCorsHeaders := config.HasCorsHeadersDefined()

	if !hasSecureHeaders && !hasCustomHeaders && !hasCorsHeaders {
		return nil, errors.New("headers configuration not valid")
	}

	var handler http.Handler
	nextHandler := next

	if hasSecureHeaders {
		logger.Debug("Setting up secureHeaders from %v", config)
		handler = newSecure(next, config)
		nextHandler = handler
	}

	if hasCustomHeaders || hasCorsHeaders {
		logger.Debug("Setting up customHeaders/Cors from %v", config)
		handler = NewHeader(nextHandler, config)
	}

	return &headers{
		handler: handler,
		name:    name,
	}, nil
}

func (h *headers) GetTracingInformation() (string, ext.SpanKindEnum) {
	return h.name, tracing.SpanKindNoneEnum
}

func (h *headers) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(rw, req)
}

type secureHeader struct {
	next   http.Handler
	secure *secure.Secure
}

// newSecure constructs a new secure instance with supplied options.
func newSecure(next http.Handler, headers dynamic.Headers) *secureHeader {
	opt := secure.Options{
		BrowserXssFilter:        headers.BrowserXSSFilter,
		ContentTypeNosniff:      headers.ContentTypeNosniff,
		ForceSTSHeader:          headers.ForceSTSHeader,
		FrameDeny:               headers.FrameDeny,
		IsDevelopment:           headers.IsDevelopment,
		SSLRedirect:             headers.SSLRedirect,
		SSLForceHost:            headers.SSLForceHost,
		SSLTemporaryRedirect:    headers.SSLTemporaryRedirect,
		STSIncludeSubdomains:    headers.STSIncludeSubdomains,
		STSPreload:              headers.STSPreload,
		ContentSecurityPolicy:   headers.ContentSecurityPolicy,
		CustomBrowserXssValue:   headers.CustomBrowserXSSValue,
		CustomFrameOptionsValue: headers.CustomFrameOptionsValue,
		PublicKey:               headers.PublicKey,
		ReferrerPolicy:          headers.ReferrerPolicy,
		SSLHost:                 headers.SSLHost,
		AllowedHosts:            headers.AllowedHosts,
		HostsProxyHeaders:       headers.HostsProxyHeaders,
		SSLProxyHeaders:         headers.SSLProxyHeaders,
		STSSeconds:              headers.STSSeconds,
		FeaturePolicy:           headers.FeaturePolicy,
	}

	return &secureHeader{
		next:   next,
		secure: secure.New(opt),
	}
}

func (s secureHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	s.secure.HandlerFuncWithNextForRequestOnly(rw, req, s.next.ServeHTTP)
}

// Header is a middleware that helps setup a few basic security features.
// A single headerOptions struct can be provided to configure which features should be enabled,
// and the ability to override a few of the default values.
type Header struct {
	next             http.Handler
	hasCustomHeaders bool
	hasCorsHeaders   bool
	headers          *dynamic.Headers
}

// NewHeader constructs a new header instance from supplied frontend header struct.
func NewHeader(next http.Handler, headers dynamic.Headers) *Header {
	hasCustomHeaders := headers.HasCustomHeadersDefined()
	hasCorsHeaders := headers.HasCorsHeadersDefined()

	return &Header{
		next:             next,
		headers:          &headers,
		hasCustomHeaders: hasCustomHeaders,
		hasCorsHeaders:   hasCorsHeaders,
	}
}

func (s *Header) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Handle Cors headers and preflight if configured.
	if isPreflight := s.processCorsHeaders(rw, req); isPreflight {
		return
	}

	if s.hasCustomHeaders {
		s.modifyCustomRequestHeaders(req)
	}

	// If there is a next, call it.
	if s.next != nil {
		s.next.ServeHTTP(rw, req)
	}
}

// modifyCustomRequestHeaders sets or deletes custom request headers.
func (s *Header) modifyCustomRequestHeaders(req *http.Request) {
	// Loop through Custom request headers
	for header, value := range s.headers.CustomRequestHeaders {
		if value == "" {
			req.Header.Del(header)
		} else {
			req.Header.Set(header, value)
		}
	}
}

// preRequestModifyCorsResponseHeaders sets during request processing time,
// all the CORS response headers that we already know that are supposed to be set,
// and which do not depend on a later state of the response.
// One notable example of a header that can only be modified later on is "Vary",
// And this is set in the post-response response modifier method
func (s *Header) preRequestModifyCorsResponseHeaders(rw http.ResponseWriter, req *http.Request) {
	originHeader := req.Header.Get("Origin")
	allowOrigin := s.getAllowOrigin(originHeader)

	if allowOrigin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", allowOrigin)
	}

	if s.headers.AccessControlAllowCredentials {
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if len(s.headers.AccessControlExposeHeaders) > 0 {
		exposeHeaders := strings.Join(s.headers.AccessControlExposeHeaders, ",")
		rw.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
	}
}

// PostRequestModifyResponseHeaders set or delete response headers.
// This method is called AFTER the response is generated from the backend
// and can merge/override headers from the backend response.
func (s *Header) PostRequestModifyResponseHeaders(res *http.Response) error {
	// Loop through Custom response headers
	for header, value := range s.headers.CustomResponseHeaders {
		if value == "" {
			res.Header.Del(header)
		} else {
			res.Header.Set(header, value)
		}
	}
	if !s.headers.AddVaryHeader {
		return nil
	}

	varyHeader := res.Header.Get("Vary")
	if varyHeader == "Origin" {
		return nil
	}

	if varyHeader != "" {
		varyHeader += ","
	}
	varyHeader += "Origin"

	res.Header.Set("Vary", varyHeader)
	return nil
}

// processCorsHeaders processes the incoming request,
// and returns if it is a preflight request.
// If not a preflight, it handles the preRequestModifyCorsResponseHeaders.
func (s *Header) processCorsHeaders(rw http.ResponseWriter, req *http.Request) bool {
	if !s.hasCorsHeaders {
		return false
	}

	reqAcMethod := req.Header.Get("Access-Control-Request-Method")
	originHeader := req.Header.Get("Origin")

	if reqAcMethod != "" && originHeader != "" && req.Method == http.MethodOptions {
		// If the request is an OPTIONS request with an Access-Control-Request-Method header,
		// and Origin headers, then it is a CORS preflight request,
		// and we need to build a custom response: https://www.w3.org/TR/cors/#preflight-request
		if s.headers.AccessControlAllowCredentials {
			rw.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		allowHeaders := strings.Join(s.headers.AccessControlAllowHeaders, ",")
		if allowHeaders != "" {
			rw.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		}

		allowMethods := strings.Join(s.headers.AccessControlAllowMethods, ",")
		if allowMethods != "" {
			rw.Header().Set("Access-Control-Allow-Methods", allowMethods)
		}

		allowOrigin := s.getAllowOrigin(originHeader)

		if allowOrigin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}

		rw.Header().Set("Access-Control-Max-Age", strconv.Itoa(int(s.headers.AccessControlMaxAge)))
		return true
	}

	s.preRequestModifyCorsResponseHeaders(rw, req)
	return false
}

func (s *Header) getAllowOrigin(header string) string {
	switch s.headers.AccessControlAllowOrigin {
	case "origin-list-or-null":
		if len(header) == 0 {
			return "null"
		}
		return header
	case "*":
		return "*"
	default:
		return ""
	}
}
