package main

import (
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/devapp/models"
)

func main() {
	for _, model := range eel.GetRegisteredModels() {
		eel.Logger.Infof("[db] migrate table %s", model.(eel.IModel).TableName())
	}
	models.Db.AutoMigrate(eel.GetRegisteredModels()...)
}

