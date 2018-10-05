// Copyright 2018 eel Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eel

import (
	"net/http"
	"time"
	"os"
	
	"github.com/kfchen81/eel/log"
	"github.com/kfchen81/eel/config"
	"fmt"
	"github.com/kfchen81/eel/router"
	"github.com/kfchen81/eel/handler"
	"github.com/kfchen81/eel/util"
)

type Request handler.Request

const logo = `
    ________    __
   / ____/ /   / /
  / __/ / /   / /
 / /___/ /___/ /___
/_____/_____/_____/  for speed & efficiency. v0.1


`

var endRunning chan bool
//
//func handler(resp http.ResponseWriter, req *http.Request) {
//
//}

// export inner type
type Context = handler.Context
type RestResource = handler.RestResource
type Map = handler.Map
type Middleware = handler.Middleware

func RegisterResource(resource handler.RestResourceInterface) {
	router.DoRegisterResource(resource)
}

func RegisterMiddleware(middleware handler.MiddlewareInterface) {
	router.DoRegisterMiddleware(middleware)
}

func NewBusinessError(code string, msg string) *util.BusinessError{
	return util.NewBusinessError(code, msg)
}


type Service struct {
	Handler *router.RestResourceRegister
	Server   *http.Server
}

func NewService() *Service {
	app := &Service{
		Handler: router.NewRestResourceRegister(),
		Server: &http.Server{},
	}
	return app
}

func (this *Service) run() {
	host := config.ServiceConfig.String("service::HOST")
	httpPort := config.ServiceConfig.String("service::HTTP_PORT")
	addr := fmt.Sprintf("%s:%s", host, httpPort)

	this.Server.Handler = this.Handler
	readTimeout := time.Duration(config.ServiceConfig.DefaultInt("service::READ_TIMEOUT", 30))
	writeTimeout := time.Duration(config.ServiceConfig.DefaultInt("service::WRITE_TIMEOUT", 10))
	readTimeout = 30
	this.Server.ReadTimeout = readTimeout * time.Second
	this.Server.WriteTimeout = writeTimeout * time.Second
	this.Server.Addr = addr
	
	log.Infof("http server Running on http://%s\n", this.Server.Addr)
	if err := this.Server.ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe: ", err)
		time.Sleep(100 * time.Microsecond)
		endRunning <- true
	}
}

func RunService() {
	os.Stderr.Write([]byte(logo))
	//fmt.Println(logo)
	service := NewService()
	endRunning = make(chan bool, 1)
	go func() {
		service.run()
	}()
	<-endRunning
}


