package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	goauth "github.com/abbot/go-http-auth"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/middlewares/accesslog"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	basicTypeName = "BasicAuth"
)

type basicAuth struct {
	next         http.Handler
	auth         *goauth.BasicAuth
	users        map[string]string
	headerField  string
	removeHeader bool
	name         string
}

// NewBasic creates a basicAuth middleware.
func NewBasic(ctx context.Context, next http.Handler, authConfig dynamic.BasicAuth, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, basicTypeName)).Debug("Creating middleware")
	users, err := getUsers(authConfig.UsersFile, authConfig.Users, basicUserParser)
	if err != nil {
		return nil, err
	}

	ba := &basicAuth{
		next:         next,
		users:        users,
		headerField:  authConfig.HeaderField,
		removeHeader: authConfig.RemoveHeader,
		name:         name,
	}

	realm := defaultRealm
	if len(authConfig.Realm) > 0 {
		realm = authConfig.Realm
	}
	ba.auth = goauth.NewBasicAuthenticator(realm, ba.secretBasic)

	return ba, nil
}

func (b *basicAuth) GetTracingInformation() (string, ext.SpanKindEnum) {
	return b.name, tracing.SpanKindNoneEnum
}

func (b *basicAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := log.FromContext(middlewares.GetLoggerCtx(req.Context(), b.name, basicTypeName))

	if username := b.auth.CheckAuth(req); username == "" {
		logger.Debug("Authentication failed")
		tracing.SetErrorWithEvent(req, "Authentication failed")
		b.auth.RequireAuth(rw, req)
	} else {
		logger.Debug("Authentication succeeded")
		req.URL.User = url.User(username)

		logData := accesslog.GetLogData(req)
		if logData != nil {
			logData.Core[accesslog.ClientUsername] = username
		}

		if b.headerField != "" {
			req.Header[b.headerField] = []string{username}
		}

		if b.removeHeader {
			logger.Debug("Removing authorization header")
			req.Header.Del(authorizationHeader)
		}
		b.next.ServeHTTP(rw, req)
	}
}

func (b *basicAuth) secretBasic(user, realm string) string {
	if secret, ok := b.users[user]; ok {
		return secret
	}

	return ""
}

func basicUserParser(user string) (string, string, error) {
	split := strings.Split(user, ":")
	if len(split) != 2 {
		return "", "", fmt.Errorf("error parsing BasicUser: %v", user)
	}
	return split[0], split[1], nil
}
