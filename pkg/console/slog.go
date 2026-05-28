package console

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
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
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
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
	fmt.Fprintln(h.w,
		FormatText(r.Time.Format(timeFormat), Blue),
		formatLevel(r.Level),
		FormatText(r.Message, White, Bold),
	)

	attrs := make([]slog.Attr, 0, len(h.attrs)+r.NumAttrs()+1)
	attrs = append(attrs, slog.Attr{
		Key:   "src",
		Value: slog.StringValue(fmt.Sprintf("%s:%d", strings.TrimPrefix(r.Source().File, h.wd), r.Source().Line)),
	})
	attrs = append(attrs, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	printAttrs(attrs, 1)
	fmt.Println()

	return nil
}

func printAttrs(attrs []slog.Attr, padding int) {
	maxKey := 0
	for _, a := range attrs {
		if len(a.Key) > maxKey {
			maxKey = len(a.Key)
		}
	}
	for _, a := range attrs {
		fmt.Print(strings.Repeat("    ", padding), FormatText(fmt.Sprintf("%s:", a.Key), Gray100), strings.Repeat(" ", maxKey-len(a.Key)+1))
		printValue(a.Value, padding)
	}
}

func printValue(v slog.Value, padding int) {
	if v.Kind() == slog.KindString {
		//fmt.Print("\"", v.String(), "\"\n")
		fmt.Print(v.String(), "\n")
		return
	}
	if v.Kind() == slog.KindGroup {
		fmt.Println()
		printAttrs(v.Group(), padding+1)
		return
	}
	fmt.Print(v.String(), "\n")
}

func formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return FormatText("DEBUG", Black, Bold)
	case slog.LevelInfo:
		return FormatText("INFO", Green, Bold)
	case slog.LevelWarn:
		return FormatText("WARN", Yellow, Bold)
	case slog.LevelError:
		return FormatText("ERROR", Red, Bold)
	default:
		return level.String()
	}
}
