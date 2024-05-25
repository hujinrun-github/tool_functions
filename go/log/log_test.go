package log

import (
	"testing"
	"time"
)

// use case 1: logger level: Debug, print level: Error
// expection: print log
func TestLoggerDebugPrintErr(t *testing.T) {
	logger := NewLogHandler(
		WithLevel(LEVEL_DEBUG),
	)

	time.Sleep(time.Second)
	logger.Error("test")
}

// use case 1: logger level: Error, print level: Debug
// expection: print noting
func TestLoggerDebugPrintInfo(t *testing.T) {
	logger := NewLogHandler(
		WithLevel(LEVEL_DEBUG),
	)

	time.Sleep(time.Second)
	logger.Info("test")
}

// use case 3: logger level: Debug, print level: Error
// expection: print log
func TestLoggerDebugPrintErrWithFile(t *testing.T) {
	logger := NewLogHandler(
		WithLevel(LEVEL_DEBUG),
		WithOutput(OUTPUT_FILE),
		WithFileName("log.log"),
	)

	time.Sleep(time.Second)
	logger.Error("test")
	// os.Remove("log.log")
}

// use case 4: logger level: Error, print level: Debug
// expection: print noting
func TestLoggerDebugPrintInfoWithFile(t *testing.T) {
	logger := NewLogHandler(
		WithLevel(LEVEL_DEBUG),
		WithOutput(OUTPUT_FILE),
		WithFileName("log.log"),
	)

	time.Sleep(time.Second)
	logger.Info("test")
	// os.Remove("log.log")
}

// use case 5: logger level: Error, print level: Debug, with white list
// expection: print log
func TestLoggerDebugPrintInfoWithWhiteList(t *testing.T) {
	logger := NewLogHandler(
		WithLevel(LEVEL_ERROR),
		WithWhiteList([]string{"test"}),
	)

	time.Sleep(time.Second)
	logger.WDebug("test", "this is 1st test")
	// os.Remove("log.log")
}

// use case 6: logger level: Error, print level: Debug, with white list
// expection: print log
func TestLoggerDebugPrintInfoWithWhiteListWrong(t *testing.T) {
	logger := NewLogHandler(
		WithLevel(LEVEL_ERROR),
		WithWhiteList([]string{"tes"}),
	)

	time.Sleep(time.Second)
	logger.WDebug("test", "this is 1st test")
	// os.Remove("log.log")
}
