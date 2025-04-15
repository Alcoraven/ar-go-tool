package arlog

import (
	"fmt"
	"runtime"
	"time"
)

type LogDet struct {
	Title   string
	Content any
}

func Det(title string, content any) LogDet {
	return LogDet{title, content}
}

func getPrefix(timeLoc *time.Location, depth int, level LogLevel) string {
	var prefix string
	var logTypeTxt string

	switch level {
	case LogLevelDebug:
		logTypeTxt = "DEBUG"
	case LogLevelInfo:
		logTypeTxt = "INFO "
	case LogLevelWarn:
		logTypeTxt = "WARN "
	case LogLevelError:
		logTypeTxt = "ERROR"
	case LogLevelFatal:
		logTypeTxt = "FATAL"
	}

	prefix += fmt.Sprintf(
		"[%s][%s]",
		logTypeTxt,
		time.Now().In(timeLoc).Format("2006-01-02 15:04:05"),
	)
	{
		if pc, _, line, ok := runtime.Caller(depth + 3); ok {
			var fullName = runtime.FuncForPC(pc).Name()
			// var fileName = filepath.Base(file)

			prefix += fmt.Sprintf(
				"[%s:%d]",
				fullName,
				line,
			)
		}
	}

	return prefix
}

func formatMsg(list ...interface{}) (msg string) {
	for i, v := range list {
		if i > 0 {
			// msg += " "
		}
		msg += formatWord(v)
	}

	return
}

func formatWord(v interface{}) (msg string) {
	switch v.(type) {
	case string:
		msg += fmt.Sprintf("%s", v)
	case int, int32, int64, uint, uint32, uint64:
		msg += fmt.Sprintf("%d", v)
	case float32, float64:
		msg += fmt.Sprintf("%.2f", v)
	case LogDet:
		det := v.(LogDet)
		msg += fmt.Sprintf("[%s:%s]", det.Title, formatWord(det.Content))
	case error:
		msg += fmt.Sprintf("%s", v.(error).Error())
	default:
		msg += fmt.Sprintf("%v", v)
	}
	return
}

func printStd(level LogLevel, list ...string) {
	var colorPrefix string
	switch level {
	case LogLevelDebug:
		colorPrefix += "\033[32m"
	case LogLevelInfo:
		colorPrefix += "\033[37m"
	case LogLevelWarn:
		colorPrefix += "\033[33m"
	case LogLevelError:
		colorPrefix += "\033[31m"
	case LogLevelFatal:
		colorPrefix += "\033[35m"
	}

	for _, v := range list {
		fmt.Println(colorPrefix + v + "\033[0m")
	}
}
