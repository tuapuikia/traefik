package haystack

import (
	"github.com/containous/traefik/v2/pkg/log"
)

/*NullLogger does nothing*/
type haystackLogger struct {
	logger log.Logger
}

/*Error prints the error message*/
func (l haystackLogger) Error(format string, v ...interface{}) {
	l.logger.Errorf(format, v)
}

/*Info prints the info message*/
func (l haystackLogger) Info(format string, v ...interface{}) {
	l.logger.Infof(format, v)
}

/*Debug prints the info message*/
func (l haystackLogger) Debug(format string, v ...interface{}) {
	l.logger.Debug(format, v)
}
