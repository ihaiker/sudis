package dao

import (
	"github.com/ihaiker/gokit/logs"
	"xorm.io/xorm/log"
)

var logger = logs.GetLogger("dao")

var coreLogger *XormLogger

func init() {
	coreLogger = new(XormLogger)
}

type XormLogger struct {
	showsql bool
}

func (self *XormLogger) Level() log.LogLevel {
	return log.LogLevel(int(logger.Level()))
}

func (self *XormLogger) SetLevel(l log.LogLevel) {
	logger.SetLevel(logs.Level(int(l)))
}

func (self *XormLogger) Debug(v ...interface{}) {
	logger.Debug(v...)
}

func (self *XormLogger) Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func (self *XormLogger) Error(v ...interface{}) {
	logger.Error(v...)
}

func (self *XormLogger) Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func (self *XormLogger) Info(v ...interface{}) {
	logger.Info(v...)
}

func (self *XormLogger) Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func (self *XormLogger) Warn(v ...interface{}) {
	logger.Warn(v...)
}

func (self *XormLogger) Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func (self *XormLogger) ShowSQL(show ...bool) {
	self.showsql = show[0]
}

func (self *XormLogger) IsShowSQL() bool {
	return self.showsql
}
