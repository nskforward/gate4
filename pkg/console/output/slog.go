package output

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	logPadding = "    "
	timeFormat = "2006-01-02 15:04:05"
)

var (
	logDebugStr = FormatText("DEBUG", Black, Bold)
	logInfoStr  = FormatText("INFO", Green, Bold)
	logWarnStr  = FormatText("WARN", Yellow, Bold)
	logErrorStr = FormatText("ERROR", Red, Bold)

	logTimePrefix = BuildPrefix(Cyan)
	logMsgPrefix  = BuildPrefix(White, Bold)
	logKeyPrefix  = BuildPrefix(Gray100)
)

var (
	bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 2048))
		},
	}
)

type LoggerHandler struct {
	w     io.Writer
	level slog.Level
	attrs []slog.Attr
	group string
	wd    string
}

func NewLogger(w io.Writer, level slog.Level) *slog.Logger {
	wd, _ := os.Getwd()
	return slog.New(&LoggerHandler{
		w:     w,
		level: level,
		wd:    wd + "/",
	})
}

func (h *LoggerHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *LoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs), len(h.attrs)+1)
	copy(newAttrs, h.attrs)
	if h.group != "" {
		newAttrs = append(newAttrs, slog.Attr{
			Key:   h.group,
			Value: slog.GroupValue(attrs...),
		})
	} else {
		newAttrs = append(newAttrs, attrs...)
	}
	return &LoggerHandler{w: h.w, wd: h.wd, level: h.level, attrs: newAttrs}
}

func (h *LoggerHandler) WithGroup(name string) slog.Handler {
	return &LoggerHandler{w: h.w, wd: h.wd, level: h.level, attrs: h.attrs, group: name}
}

func (h *LoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	buf.WriteString(logTimePrefix)
	buf.WriteString(r.Time.Format(timeFormat))
	buf.WriteString(Reset)
	buf.WriteByte(' ')
	buf.WriteString(formatLevel(r.Level))
	buf.WriteByte(' ')
	buf.WriteString(logMsgPrefix)
	buf.WriteString(r.Message)
	buf.WriteString(Reset)
	buf.WriteByte('\n')

	attrs := make([]slog.Attr, 0, len(h.attrs)+r.NumAttrs())

	if r.PC != 0 {
		file := strings.TrimPrefix(r.Source().File, h.wd)
		srcBuf := make([]byte, 0, len(file)+6) // максимум 5 цифр на номер строки
		srcBuf = append(srcBuf, file...)
		srcBuf = append(srcBuf, ':')
		srcBuf = strconv.AppendInt(srcBuf, int64(r.Source().Line), 10)
		attrs = append(attrs, slog.Attr{Key: "src", Value: slog.StringValue(string(srcBuf))})
	}

	attrs = append(attrs, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	writeAttrs(buf, attrs, 1)
	buf.WriteByte('\n')

	_, err := buf.WriteTo(h.w)
	return err
}

func writeAttrs(buf *bytes.Buffer, attrs []slog.Attr, padding int) {
	maxKey := 0
	for _, a := range attrs {
		if len(a.Key) > maxKey {
			maxKey = len(a.Key)
		}
	}

	for _, a := range attrs {
		for range padding {
			buf.WriteString(logPadding)
		}
		buf.WriteString(logKeyPrefix)
		buf.WriteString(a.Key)
		buf.WriteByte(':')
		buf.WriteString(Reset)
		for range maxKey - len(a.Key) + 1 {
			buf.WriteByte(' ')
		}

		switch a.Value.Kind() {
		case slog.KindString:
			buf.WriteString(a.Value.String())
			buf.WriteByte('\n')
		case slog.KindGroup:
			buf.WriteByte('\n')
			writeAttrs(buf, a.Value.Group(), padding+1)
		default:
			buf.WriteString(a.Value.String())
			buf.WriteByte('\n')
		}
	}
}

func formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return logDebugStr
	case slog.LevelInfo:
		return logInfoStr
	case slog.LevelWarn:
		return logWarnStr
	case slog.LevelError:
		return logErrorStr
	default:
		return level.String()
	}
}
