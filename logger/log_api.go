package logger

import (
	"fmt"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
)

var innerCallDepth = 1

func Info(v ...interface{}) {
	Holder.Info(getMessage("", v), zap.String("caller", getCaller(1)))
}

func Infof(format string, v ...interface{}) {
	Holder.Info(getMessage(format, v), zap.String("caller", getCaller(1)))
}

func Debug(v ...interface{}) {
	Holder.Debug(getMessage("", v), zap.String("caller", getCaller(1)))
}

func Debugf(format string, v ...interface{}) {
	Holder.Debug(getMessage(format, v), zap.String("caller", getCaller(1)))
}

func Warning(v ...interface{}) {
	Holder.Warn(getMessage("", v), zap.String("caller", getCaller(1)))
}

func Warningf(format string, v ...interface{}) {
	Holder.Warn(getMessage(format, v), zap.String("caller", getCaller(1)))
}

func Error(v ...interface{}) {
	Holder.Error(getMessage("", v), zap.String("caller", getCaller(1)))
}

func Errorf(format string, v ...interface{}) {
	Holder.Error(getMessage(format, v), zap.String("caller", getCaller(1)))
}

func Fatal(v ...interface{}) {
	Holder.Fatal(getMessage("", v), zap.String("caller", getCaller(1)))
}

func Fatalf(format string, v ...interface{}) {
	Holder.Fatal(getMessage(format, v), zap.String("caller", getCaller(1)))
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func getCaller(callDepth int) string {
	var buf strings.Builder

	_, file, line, ok := runtime.Caller(callDepth + innerCallDepth)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		buf.WriteString(short)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
	}

	return buf.String()
}
