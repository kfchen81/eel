package db

import (
	"time"
	"fmt"
	
	"github.com/kfchen81/gorm"
	"github.com/kfchen81/eel/log"
	"github.com/opentracing/opentracing-go"
)

type Model struct {
	Id        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeletableModel struct {
	Id        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
	//DeletedAt *time.Time `sql:"index"`
}

//记录所有的model，用于支持syncdb
var models = make([]interface{}, 0)

func RegisterModel(model interface{}) {
	models = append(models, model)
}

func GetRegisteredModels() []interface{} {
	return models
}

func createSubSpan(scope *gorm.Scope, operationName string) {
	if span, ok := scope.Get("rootSpan"); ok {
		rootSpan := span.(opentracing.Span)
		operationName = fmt.Sprintf("db-%s", operationName)
		subSpan := rootSpan.Tracer().StartSpan(
			operationName,
			opentracing.ChildOf(rootSpan.Context()),
		)
		
		scope.Set("subSpan", subSpan)
		
	}
}

func finishSubSpan(scope *gorm.Scope) {
	if span, ok := scope.Get("subSpan"); ok {
		subSpan := span.(opentracing.Span)
		defer subSpan.Finish()
		subSpan.LogKV("sql", scope.SQL)
	}
}

// support open tracing
func beginTracingForCreate(scope *gorm.Scope) {
	createSubSpan(scope, fmt.Sprintf("db-insert-%s", scope.TableName()))
}
func finishTracingForCreate(scope *gorm.Scope) {
	finishSubSpan(scope)
}
func beginTracingForUpdate(scope *gorm.Scope) {
	createSubSpan(scope, fmt.Sprintf("db-update-%s", scope.TableName()))
}
func finishTracingForUpdate(scope *gorm.Scope) {
	finishSubSpan(scope)
}
func beginTracingForDelete(scope *gorm.Scope) {
	createSubSpan(scope, fmt.Sprintf("db-delete-%s", scope.TableName()))
}
func finishTracingForDelete(scope *gorm.Scope) {
	finishSubSpan(scope)
}
func beginTracingForQuery(scope *gorm.Scope) {
	createSubSpan(scope, fmt.Sprintf("db-select-%s", scope.TableName()))
}
func finishTracingForQuery(scope *gorm.Scope) {
	finishSubSpan(scope)
}

func init() {
	log.Logger.Info("register GORM's callback")
	gorm.DefaultCallback.Create().Before("gorm:before_create").Register("eel:before_create", beginTracingForCreate)
	gorm.DefaultCallback.Create().After("gorm:after_create").Register("eel:after_create", finishTracingForCreate)
	
	gorm.DefaultCallback.Update().Before("gorm:before_update").Register("eel:before_update", beginTracingForUpdate)
	gorm.DefaultCallback.Update().After("gorm:after_update").Register("eel:after_update", finishTracingForUpdate)
	
	gorm.DefaultCallback.Delete().Before("gorm:before_delete").Register("eel:before_delete", beginTracingForDelete)
	gorm.DefaultCallback.Delete().After("gorm:after_delete").Register("eel:after_delete", finishTracingForDelete)
	
	gorm.DefaultCallback.Query().Before("gorm:query").Register("eel:before_query", beginTracingForQuery)
	gorm.DefaultCallback.Query().After("gorm:after_query").Register("eel:after_query", finishTracingForQuery)
}