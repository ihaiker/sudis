package daemon

import (
	"github.com/ihaiker/gokit/concurrent/atomic"
	"github.com/ihaiker/gokit/errors"
	"github.com/ihaiker/gokit/logs"
	"io"
)

const LOGGER_LINES_NUM = 100

type TailLogger func(id, line string)

type ProcessLogger struct {
	file    io.Writer
	lines   []string
	lineIdx *atomic.AtomicInt
	tails   map[string]TailLogger
}

func NewLogger(loggerFile string) (logger *ProcessLogger, err error) {
	logger = new(ProcessLogger)
	if loggerFile != "" {
		if logger.file, err = logs.NewDailyRolling(loggerFile); err != nil {
			return nil, err
		}
	}
	logger.lineIdx = atomic.NewAtomicInt(0)
	logger.lines = make([]string, LOGGER_LINES_NUM)
	logger.tails = map[string]TailLogger{}
	return
}

func (self *ProcessLogger) Tail(id string, tail TailLogger, lines int) {
	logger.Debug("logger tail : ", id)
	self.tails[id] = tail
	if lines > LOGGER_LINES_NUM {
		lines = LOGGER_LINES_NUM
	}
	idx := self.lineIdx.Get()
	if idx < lines {
		idx = 0
	} else if idx < LOGGER_LINES_NUM {
		idx -= lines
	}
	for i := idx; i < idx+lines; i++ {
		line := self.lines[(idx+i)%LOGGER_LINES_NUM]
		if line != "" {
			tail(id, line)
		}
	}
}

func (self *ProcessLogger) CtrlC(id string) {
	logger.Debug("stop logger tail: ", id)
	delete(self.tails, id)
}

func (self *ProcessLogger) Write(p []byte) (int, error) {
	line := string(p)
	self.lines[self.lineIdx.GetAndIncrement(1)%LOGGER_LINES_NUM] = line
	for id, tail := range self.tails {
		errors.Try(func() {
			tail(id, line)
		})
	}
	if self.file != nil {
		return self.file.Write(p)
	}
	return len(p), nil
}

func (self *ProcessLogger) Close() error {
	if self.file != nil {
		if closer, match := self.file.(io.Closer); match {
			return closer.Close()
		}
	}
	return nil
}
