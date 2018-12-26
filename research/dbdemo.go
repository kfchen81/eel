package main

import (
	"fmt"
	"github.com/kfchen81/eel/config"
	"github.com/kfchen81/gorm"
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/tracing"
)

var Db *gorm.DB

type Blog struct {
	eel.Model
	UserId int
	Title string
	Content string `gorm:"size:1024"`
}
func (this *Blog) TableName() string {
	return "blog_blog"
}

func main() {
	fmt.Println("hell gorm")
	
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
	defer Db.Close()
	
	Db := Db.Begin()
	defer Db.Commit()
	
	rootSpan := tracing.Tracer.StartSpan("db2")
	defer tracing.Closer.Close()
	Db.InstantSet("rootSpan", rootSpan)
	
	Db.LogMode(true)
	
	model := Blog{
		UserId: 123,
		Title: "robert1",
		Content: "content1",
	}
	
	//tx := Db.Begin()
	//tx.Create(&model)
	Db.Create(&model)
	fmt.Println(model.Id)
	
	//select
	//var model1 Blog
	//Db.Take(&model1)
	//fmt.Println(model1)
	
	//Db.Where("id in (?)", []int{69, 70}).Find(&models)
	//Db.Filter(map[string]interface{}{
	//	"id__in": []int{69, 70},
	//}).All(&models)
	
	var models []*Blog
	Db.Filter("title__contains", "robert").All(&models)
	fmt.Println(models)
	
	Db.QueryTable(&Blog{}).Filter(gorm.Map{
		"id": model.Id,
	}).Update(gorm.UpdateParam{
		"title": "updated title 1",
	})
	
	//delete
	Db.QueryTable(&Blog{}).Filter("id", model.Id).Delete()
}