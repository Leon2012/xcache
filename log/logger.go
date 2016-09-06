package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	colors map[string]string
	logNo  uint64
)

type LEVEL int32

var logLevel LEVEL = 1

const (
	Black = (iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

//Logger Level
const (
	ALL LEVEL = iota
	DEBUG
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
	OFF
)

type Worker struct {
	Minion  *log.Logger
	Color   int
	LogFile *os.File
	mu      *sync.Mutex
}

type Infoer struct {
	Id      uint64
	Time    string
	Module  string
	Level   string
	Message string
	format  string
}

type Logger struct {
	Module string
	Worker *Worker
}

func SetLevel(l LEVEL) {
	logLevel = l
}

func NewLogger(module string, color int) (*Logger, error) {
	initColors()
	newWorker := NewConsoleWorker("", 0, color)
	return &Logger{Module: module, Worker: newWorker}, nil
}

func NewFileLogger(module string, color int, logFile string) (*Logger, error) {
	fileHandler, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	} else {
		initColors()
		newWorker := NewFileWorker("", 0, color, fileHandler)
		return &Logger{Module: module, Worker: newWorker}, nil
	}
}

func NewDailyLogger(module string, color int, logPath string) (*Logger, error) {

	var logFile string
	const layout = "2006-01-02"
	now := time.Now()
	fileName := now.Format(layout)
	if len(logPath) == 0 {
		logFile = "./" + fileName + ".log"
	} else {
		logFile = logPath + "/" + fileName + ".log"
	}
	fileHandler, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	} else {
		initColors()
		newWorker := NewFileWorker("", 0, color, fileHandler)
		return &Logger{Module: module, Worker: newWorker}, nil
	}
}

func (l *Logger) Log(lvl string, message string) {
	//var formatString string = "#%d %s ▶ %.3s %s"
	var formatString string = "%s [%s] %.4s %s"
	info := &Infoer{
		Id:      atomic.AddUint64(&logNo, 1),
		Time:    time.Now().Format("2006-01-02 15:04:05.000"),
		Module:  l.Module,
		Level:   lvl,
		Message: message,
		format:  formatString,
	}
	l.Worker.Log(lvl, 2, info)
}

func (l *Logger) SetModule(m string) {
	l.Module = m
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	if logLevel <= FATAL {
		message := fmt.Sprintf(format, v...)
		l.Log("FATAL", message)
		os.Exit(1)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if logLevel <= ERROR {
		message := fmt.Sprintf(format, v...)
		l.Log("ERROR", message)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	if logLevel <= WARN {
		message := fmt.Sprintf(format, v...)
		l.Log("WARNING", message)
	}

}

func (l *Logger) Notice(format string, v ...interface{}) {
	if logLevel <= NOTICE {
		message := fmt.Sprintf(format, v...)
		l.Log("NOTICE", message)
	}

}

func (l *Logger) Info(format string, v ...interface{}) {
	if logLevel <= INFO {
		message := fmt.Sprintf(format, v...)
		l.Log("INFO", message)
	}

}

func (l *Logger) Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		message := fmt.Sprintf(format, v...)
		l.Log("DEBUG", message)
	}

}

func (l *Logger) Panic(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.Log("FATAL", message)
	panic(message)
}

func (l *Logger) Critical(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.Log("FATAL", message)
}

func (l *Logger) Strack(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	message += "\n"
	buf := make([]byte, 1024*1024)
	n := runtime.Stack(buf, true)
	message += string(buf[:n])
	message += "\n"
	l.Log("FATAL", message)
}

func (l *Logger) Close() {
	if l.Worker.LogFile != nil {
		l.Worker.LogFile.Close()
	}
}

func NewWorker(prefix string, flag int, color int, out io.Writer) *Worker {
	return &Worker{
		Minion:  log.New(out, prefix, flag),
		Color:   color,
		LogFile: nil,
		mu:      new(sync.Mutex)}
}

func NewConsoleWorker(prefix string, flag int, color int) *Worker {
	return NewWorker(prefix, flag, color, os.Stdout)
}

func NewFileWorker(prefix string, flag int, color int, logFile *os.File) *Worker {
	return &Worker{
		Minion:  log.New(logFile, prefix, flag),
		Color:   color,
		LogFile: logFile,
		mu:      new(sync.Mutex),
	}
}

func (w *Worker) Log(level string, calldepth int, info *Infoer) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Color != 0 {
		buf := &bytes.Buffer{}
		buf.Write([]byte(colors[level]))
		buf.Write([]byte(info.Output()))
		buf.Write([]byte("\033[0m"))
		return w.Minion.Output(calldepth+1, buf.String())
	} else {
		return w.Minion.Output(calldepth+1, info.Output())
	}
}

func colorString(color int) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

func initColors() {
	colors = map[string]string{
		"FATAL":   colorString(Magenta),
		"ERROR":   colorString(Red),
		"WARNING": colorString(Yellow),
		"NOTICE":  colorString(Green),
		"DEBUG":   colorString(Cyan),
		"INFO":    colorString(White),
	}
}

func (i *Infoer) Output() string {
	// depthOfFunctionCaller := 1
	// pc, _, _, _ := runtime.Caller(depthOfFunctionCaller)
	// fn := runtime.FuncForPC(pc)
	// elems := strings.Split(fn.Name(), ".")
	// fi := elems[len(elems)-1]
	msg := fmt.Sprintf(i.format, i.Module, i.Time, i.Level, i.Message)
	return msg
}
