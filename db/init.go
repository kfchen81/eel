package db

import (
	"time"
)

type Model struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeletableModel struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
	DeletedAt *time.Time `sql:"index"`
}

//记录所有的model，用于支持syncdb
var models = make([]interface{}, 0)

func RegisterModel(model interface{}) {
	models = append(models, model)
}

func GetRegisteredModels() []interface{} {
	return models
}

func init() {

}