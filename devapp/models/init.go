package models

import (
	_ "github.com/kfchen81/gorm/dialects/mysql"
	"github.com/kfchen81/gorm"
	"fmt"
	
	_ "github.com/kfchen81/eel/devapp/models/blog"
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/config"
)

var Db *gorm.DB

func init() {
	host := config.ServiceConfig.String("db::DB_HOST")
	port := config.ServiceConfig.String("db::DB_PORT")
	db := config.ServiceConfig.String("db::DB_NAME")
	user := config.ServiceConfig.String("db::DB_USER")
	password := config.ServiceConfig.String("db::DB_PASSWORD")
	charset := config.ServiceConfig.String("db::DB_CHARSET")
	mysqlURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Asia%%2FShanghai", user, password, host, port, db, charset)
	
	var err error
	Db, err = gorm.Open("mysql", mysqlURL)
	
	if err != nil {
		eel.Logger.Errorw("[db] connect to mysql fail!!", "error", err.Error())
	} else {
		eel.Logger.Infof("[db] connect to mysql %s success!", mysqlURL)
	}
	
	eel.Runtime.DB = Db
}
