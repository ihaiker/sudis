package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ihaiker/gokit/files"
	"github.com/ihaiker/sudis/libs/config"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"strings"
	"time"

	"github.com/ihaiker/gokit/logs"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func CreateEngine(path string, config *config.Database) (err error) {
	datasource := config.Url
	dbType := config.Type
	if !strings.HasPrefix(datasource, "/") {
		datasource = filepath.Join(path, datasource)
	}
	logs.Infof("connect %s %s", dbType, datasource)

	notExists := files.NotExist(datasource)
	engine, err = xorm.NewEngine(dbType, datasource)
	if err != nil {
		return err
	}
	engine.SetLogger(coreLogger)
	if logs.IsDebugMode() {
		engine.ShowSQL(true)
	}
	if notExists {
		err = createTables()
	}
	return
}

func createTables() error {
	tables := []interface{}{
		&Node{}, &Tag{},
		&User{}, &Notify{},
	}
	if err := engine.CreateTables(tables...); err != nil {
		return err
	}
	user := &User{
		Name: "admin", Passwd: "12345678", Time: Timestamp(),
	}
	if _, err := engine.InsertOne(user); err != nil {
		return err
	}
	return nil
}

func Timestamp() string {
	return time.Now().Format("2006-01-03 15:04:05")
}
