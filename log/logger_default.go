package log

const (
	DEFAULT_LOGGER_MODULE = "DEFAULT"
)

var defaultLogger *Logger

func init() {
	defaultLogger, _ = NewLogger(DEFAULT_LOGGER_MODULE, 0)
}

func SetModule(m string) {
	defaultLogger.SetModule(m)
}

func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Notice(format string, v ...interface{}) {
	defaultLogger.Notice(format, v...)
}

func Warning(format string, v ...interface{}) {
	defaultLogger.Warning(format, v...)
}

func Debug(format string, v ...interface{}) {
	defaultLogger.Debug(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	defaultLogger.Fatal(format, v...)
}
