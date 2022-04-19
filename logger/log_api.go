package logger

func Info(v ...interface{}) {
	Sugar.Info(v...)
}

func Infof(format string, v ...interface{}) {
	Sugar.Infof(format, v...)
}

func Debug(v ...interface{}) {
	Sugar.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	Sugar.Infof(format, v...)
}

func Warning(v ...interface{}) {
	Sugar.Info(v...)
}

func Warningf(format string, v ...interface{}) {
	Sugar.Warnf(format, v...)
}

func Error(v ...interface{}) {
	Sugar.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	Sugar.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	Sugar.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	Sugar.Fatalf(format, v...)
}
