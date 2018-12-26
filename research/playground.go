package main

import (
	"time"
	"fmt"
	"github.com/pkg/errors"
)

type IDable interface {
	GetId() int
}

type Article struct {
}
func (this *Article) GetId() int {
	return 1
}

func GetId2(idable interface{}) (int, error) {
	return idable.(IDable).GetId(), nil
	//switch value := idable.(type) {
	//case IDable:
	//	return value.GetId(), nil
	//}
	//
	//return 0, errors.New("invalid type")
}

func GetId4(idable interface{}) (int, error) {
	switch value := idable.(type) {
	case IDable:
		return value.GetId(), nil
	}
	
	return 0, errors.New("invalid type")
}

func GetId3(article *Article) (int, error) {
	return article.GetId(), nil
}

func main() {
	blog := Article{}
	COUNT := 1000000
	//use switch type
	start := time.Now()
	for i := 0; i < COUNT; i++ {
		GetId2(&blog)
	}
	fmt.Println(time.Since(start))
	
	//directly call
	start4 := time.Now()
	for i := 0; i < COUNT; i++ {
		GetId4(&blog)
	}
	fmt.Println(time.Since(start4))
	
	//directly call
	start2 := time.Now()
	for i := 0; i < COUNT; i++ {
		GetId3(&blog)
	}
	fmt.Println(time.Since(start2))
}