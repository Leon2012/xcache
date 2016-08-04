package log

import "testing"

func execLogger(logger *Logger) {
	logger.Debug("debug")
	logger.Info("info")
	logger.Notice("notice")
	logger.Warning("warning")
	logger.Error("error")
	logger.Fatal("fatal")
}

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger("test", 0)
	if err != nil {
		t.Fail()
	}
	defer logger.Close()
	execLogger(logger)
}

func TestSetLevel(t *testing.T) {
	logger, err := NewLogger("test", 0)
	if err != nil {
		t.Fail()
	}
	defer logger.Close()
	SetLevel(ERROR)
	execLogger(logger)
}

func TestNewFileLogger(t *testing.T) {
	logFile := "./app.log"
	logger, err := NewFileLogger("test", 0, logFile)
	if err != nil {
		t.Fail()
	}
	defer logger.Close()
	SetLevel(ERROR)
	execLogger(logger)
}

func TestDefaultLogger(t *testing.T) {
	Info("info")
	Debug("debug")
	Notice("notice")
	Warning("warning")
	Error("error")
	Fatal("fatal")
}
