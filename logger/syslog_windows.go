//go:build windows
// +build windows

package logger

import (
	"fmt"
)

// NewSyslogLogger 在Windows平台上提供一个兼容API，但始终返回错误
// Windows不支持syslog
func NewSyslogLogger(server, tag string, options ...Option) (LoggerWithWriter, error) {
	return nil, fmt.Errorf("syslog not supported on Windows")
}
