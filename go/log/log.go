package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/hujinrun-github/tool_functions/go/concurrent"
	str_util "github.com/hujinrun-github/tool_functions/go/string"
)

type LogLevel int

const (
	LEVEL_INFO LogLevel = iota + 1
	LEVEL_DEBUG
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

type LogOutput int

const (
	OUTPUT_CONSOLE LogOutput = iota + 100
	OUTPUT_FILE
)

type logInfoST struct {
	level   LogLevel
	message string
}

type Logger struct {
	level     LogLevel
	output    LogOutput
	fileName  string
	whiteList map[string]struct{}
	file      *os.File
	ch        chan logInfoST

	parallel int
}

func NewLogHandler(opts ...LogOption) *Logger {
	l := &Logger{
		level:    DEFAULT_LOG_LEVEL,
		output:   DEFAULT_LOG_OUTPUT,
		fileName: DEFAULT_LOG_FILE_NAME,
		parallel: DEFAULT_LOG_PARALLEL_HANDLE_SIZE,
	}

	for _, opt := range opts {
		opt(l)
	}

	log.Printf("初始化日志模块：logger:%+v\n", l)

	l.ch = make(chan logInfoST, DEFAULT_LOG_CHAN_SIZE)

	var err error
	if l.output == OUTPUT_FILE {
		l.file, err = os.OpenFile(l.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("打开文件失败：", err)
		}

		log.SetOutput(l.file)
	}

	for i := 0; i < l.parallel; i++ {
		index := i
		go concurrent.Try(func() {
			log.Printf(">>> 启动日志处理协程 %d\n", index)
			for {
				select {
				case info, ok := <-l.ch:
					if !ok {
						return
					}
					log.Println(info.message)
				}
			}
		})
	}

	return l
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level

}

func (l *Logger) write(prefix string, message ...string) {
	strMessageBuild := strings.Builder{}
	strMessageBuild.WriteString(prefix)
	strMessageBuild.WriteByte(' ')
	for _, v := range message {
		strMessageBuild.WriteString(v)
		strMessageBuild.WriteByte(' ')
	}

	l.ch <- logInfoST{
		message: strMessageBuild.String(),
	}
}

// =========> write with level <===============

func (l *Logger) Write(level LogLevel, message ...string) {
	if level < l.level {
		return
	}

	prefix := ""
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		prefix = fmt.Sprintf("[Function:%s][File:%s][Line:%d]", str_util.LastTrimAndRetain(runtime.FuncForPC(pc).Name(), "/"), str_util.LastTrimAndRetain(file, "/"), line)
	}

	l.write(prefix, message...)
}

func (l *Logger) Info(message ...string) {
	l.Write(LEVEL_INFO, message...)
}

func (l *Logger) Debug(message ...string) {
	l.Write(LEVEL_DEBUG, message...)
}

func (l *Logger) Warn(message ...string) {
	l.Write(LEVEL_WARN, message...)
}

func (l *Logger) Error(message ...string) {
	l.Write(LEVEL_ERROR, message...)
}

func (l *Logger) Fatal(message ...string) {
	log.Fatal(message)
}

func (l *Logger) WriteWithWL(id string, level LogLevel, message ...string) {
	_, ok := l.whiteList[id]
	if !ok && l.level > level {
		return
	}
	prefix := ""
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		prefix = fmt.Sprintf("[id:%s][Function:%s][File:%s][Line:%d]", id, str_util.LastTrimAndRetain(runtime.FuncForPC(pc).Name(), "/"), str_util.LastTrimAndRetain(file, "/"), line)
	}

	l.write(prefix, message...)
}

func (l *Logger) WInfo(id string, message ...string) {
	l.WriteWithWL(id, LEVEL_INFO, message...)
}

func (l *Logger) WDebug(id string, message ...string) {
	l.WriteWithWL(id, LEVEL_DEBUG, message...)
}

func (l *Logger) WWarn(id string, message ...string) {
	l.WriteWithWL(id, LEVEL_WARN, message...)
}

func (l *Logger) WError(id string, message ...string) {
	l.WriteWithWL(id, LEVEL_ERROR, message...)
}
