//go:build !windows
// +build !windows

package logger

import (
	"log/syslog"
	"strings"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

// NewSyslogLogger 创建一个连接到syslog的日志器
// 只在非Windows平台（Unix/Linux/Mac等）上可用
func NewSyslogLogger(server, tag string, options ...Option) (LoggerWithWriter, error) {
	l := newLogger(options...)
	var syslogLevel syslog.Priority

	level := slog.LevelByName(l.level)
	if level == slog.DebugLevel {
		syslogLevel = syslog.LOG_DEBUG
		// 当 log level 为 debug 时开启 caller，方便快速定位打印日志位置
		// logTemplate = "{{datetime}} {{level}} {{message}} [{{caller}}]\n"
	} else {
		syslogLevel = syslog.LOG_INFO
		// logTemplate = "{{datetime}} {{level}} {{message}}\n"
	}

	// custom log format
	logFormatter := genLogFormatter(l.timeFormat)

	if len(server) == 0 {
		// Use local syslog socket by default.
		server = "/dev/log"
	}

	if strings.HasPrefix(server, "/") {
		h, err := handler.NewSysLogHandler(syslogLevel|syslog.LOG_MAIL, tag)
		if err != nil {
			return nil, err
		}
		h.SetFormatter(logFormatter)
		l.sl.AddHandler(h)
	} else {
		w, err := syslog.Dial("tcp", server, syslogLevel|syslog.LOG_MAIL, tag)
		if err != nil {
			return nil, err
		}
		h := handler.NewBufferedHandler(w, l.bufferSize, level)
		h.SetFormatter(logFormatter)
		l.sl.AddHandler(h)
	}

	l.sl.DoNothingOnPanicFatal()

	return l, nil
}
