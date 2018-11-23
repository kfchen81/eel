package config

import (
	_ "github.com/go-sql-driver/mysql"
	"path/filepath"
	"os"
	"github.com/kfchen81/eel/log"
	"github.com/kfchen81/eel/utils"
	_ "github.com/kfchen81/eel/config/env"
)

var (
	// AppPath is the absolute path to the app
	AppPath string
	
	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "ini"
	
	ServiceConfig *serviceConfig
	Runtime *runtime
)

func parseConfig(appConfigPath string) (err error) {
	ServiceConfig, err = newServiceConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	
	return nil
	//return assignConfig(AppConfig)
}

func init() {
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	
	var filename = "app.conf"
	if os.Getenv("_SERVICE_MODE") != "" {
		filename = os.Getenv("_SERVICE_MODE") + ".app.conf"
	}
	
	appConfigPath = filepath.Join(workPath, "conf", filename)
	if !utils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "conf", filename)
		if !utils.FileExists(appConfigPath) {
			log.Logger.Panicf("no app.conf in %s directory", appConfigPath)
			panic("no app.conf in ./conf directory")
		}
	}
	if err = parseConfig(appConfigPath); err != nil {
		panic(err)
	}
}