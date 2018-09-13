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

package main

import (
	"log"
	"net/http"
	"time"
	"fmt"
)

const logo string = `
    ________    __
   / ____/ /   / /
  / __/ / /   / /
 / /___/ /___/ /___
/_____/_____/_____/   v0.1
`

var endRunning chan bool

func handler(resp http.ResponseWriter, req *http.Request) {

}

type RestResourceRegister struct {

}

func (p *RestResourceRegister) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println("in RestResourceRegister")
	respStr := "hello eel"
	resp.Write([]byte(respStr))
}

type Service struct {
	Handler *RestResourceRegister
	Server   *http.Server
}

func NewService() *Service {
	resourceRegister := &RestResourceRegister{}
	app := &Service{
		Handler: resourceRegister,
		Server: &http.Server{},
	}
	return app
}

func (this *Service) Run() {
	addr := "0.0.0.0:3130"
	this.Server.Handler = this.Handler
	this.Server.ReadTimeout = 30 * time.Second
	this.Server.WriteTimeout = 10 * time.Second
	this.Server.Addr = addr
	
	log.Printf("http server Running on http://%s\n", this.Server.Addr)
	if err := this.Server.ListenAndServe(); err != nil {
		log.Fatalln("ListenAndServe: ", err)
		time.Sleep(100 * time.Microsecond)
		endRunning <- true
	}
}

func runServer() {
	fmt.Println(logo)
	service := NewService()
	endRunning = make(chan bool, 1)
	go func() {
		service.Run()
	}()
	<-endRunning
}

func main() {
	runServer()
	http.HandleFunc("/", handler)
	log.Println("listen in 3130")
	err := http.ListenAndServe(":3130", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

