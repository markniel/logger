package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	LevelTrace  = slog.Level(-8)
	LevelStats  = slog.Level(-2)
	LevelNotice = slog.Level(2)
)

type Logger struct {
	*slog.Logger
}

func NewLogger(logType string, dest *os.File, options *slog.HandlerOptions) *Logger {
	switch logType {
	case "JSON":
		return &Logger{
			Logger: slog.New(slog.NewJSONHandler(dest, options)),
		}
	case "Text":
		return &Logger{
			Logger: slog.New(slog.NewTextHandler(dest, options)),
		}
	default:
		return nil
	}
}

func NewLogOptions(level slog.Level, includeSource bool) *slog.HandlerOptions {
	var LevelNames = map[slog.Leveler]string{
		LevelTrace:  "TRACE",
		LevelStats:  "STATS",
		LevelNotice: "NOTICE",
	}
	logHandlerOpts := &slog.HandlerOptions{
		AddSource: includeSource,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				src, ok := a.Value.Any().(*slog.Source)
				if ok {
					elements := strings.Split(src.File, "/")
					var shortPath string
					if len(elements) > 1 {
						shortPath = fmt.Sprintf("%s/%s:%d", elements[len(elements)-2], elements[len(elements)-1], src.Line)
					} else {
						shortPath = fmt.Sprintf("%s:%d", elements[len(elements)-1], src.Line)
					}
					a.Value = slog.StringValue(shortPath)
				}
			}
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	}
	return logHandlerOpts
}

func (logger *Logger) Stats(args ...interface{}) {
	if !logger.Enabled(context.Background(), LevelStats) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	function := runtime.FuncForPC(pcs[0]).Name()
	functionElements := strings.Split(function, ".")
	functionName := functionElements[len(functionElements)-1] + "()"
	fmt.Printf("%s\n", functionName)
	r := slog.NewRecord(time.Now(), LevelStats, fmt.Sprintf("%s", functionName), pcs[0])
	r.Add(args[:]...)
	_ = logger.Handler().Handle(context.Background(), r)
}

func (logger *Logger) Trace(args ...interface{}) {
	if !logger.Enabled(context.Background(), LevelTrace) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), LevelTrace, fmt.Sprintf("%s", args[0]), pcs[0])
	r.Add(args[1:]...)
	_ = logger.Handler().Handle(context.Background(), r)
}

func (logger *Logger) Notice(args ...interface{}) {
	if !logger.Enabled(context.Background(), LevelNotice) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), LevelNotice, fmt.Sprintf("%s", args[0]), pcs[0])
	r.Add(args[1:]...)
	_ = logger.Handler().Handle(context.Background(), r)
}
