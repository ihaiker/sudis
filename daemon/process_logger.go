package daemon

import (
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/logs/appenders"
	"io"
)

const LOGGER_LINES_NUM = 100

type ProcessLogger struct {
	file    io.Writer
	lines   []string
	lineIdx *atomic.AtomicInt
}

func NewLogger(loggerFile string) (logger *ProcessLogger, err error) {
	logger = new(ProcessLogger)
	writers := make([]io.Writer, 0)
	if loggerFile != "" {
		if logger.file, err = appenders.NewDailyRollingFileOut(loggerFile); err != nil {
			return nil, err
		} else {
			writers = append(writers, logger.file)
		}
	}
	logger.lineIdx = atomic.NewAtomicInt(0)
	logger.lines = make([]string, LOGGER_LINES_NUM)
	return logger, nil
}

func (self *ProcessLogger) Write(p []byte) (int, error) {
	self.lines[self.lineIdx.GetAndDecrement(1)%LOGGER_LINES_NUM] = string(p)
	if self.file != nil {
		return self.file.Write(p)
	}
	return len(p), nil
}

func (self *ProcessLogger) Close() error {
	if self.file != nil {
		return (self.file.(io.Closer)).Close()
	}
	return nil
}
