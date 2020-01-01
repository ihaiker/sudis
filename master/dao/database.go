package dao

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"time"

	"github.com/ihaiker/gokit/logs"
	"github.com/ihaiker/sudis/conf"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func CreateEngine(config *conf.Database) (err error) {
	datasource := config.Url
	dbType := config.Type
	logs.Debug("connect ", dbType, " :", datasource)

	engine, err = xorm.NewEngine(dbType, datasource)
	if err != nil {
		return err
	}
	engine.SetLogger(coreLogger)
	if logs.IsDebugMode() {
		engine.ShowSQL(true)
		engine.ShowExecTime(true)
	}
	return
}

func createTables() error {
	tables := []interface{}{
		&Node{}, &Program{},
		&Tag{}, &User{}, &Notify{},
	}
	if err := engine.CreateTables(tables...); err != nil {
		return err
	}
	user := &User{
		Name:   "admin",
		Passwd: "12345678",
		Time:   time.Now().Format("2006-01-02 15:03:04"),
	}
	if _, err := engine.InsertOne(user); err != nil {
		return err
	}
	return nil
}

func InitDatabase(config *conf.Database) (err error) {
	if err = CreateEngine(config); err != nil {
		return
	}
	defer engine.Close()
	return createTables()
}

func Timestamp() string {
	return time.Now().Format("2006-01-03 15:04:05")
}
