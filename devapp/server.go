package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"database/sql"
	"fmt"
)

func readThroughSql() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/bacchus")
	
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	
	var (
		id int
		name string
	)
	
	rows, err := db.Query("select name, id from user_mobile_user")
	defer rows.Close()
	
	for rows.Next() {
		err := rows.Scan(&name, &id)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(name, id)
	}
	
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID int
	Name string
}
func (u User)TableName() string {
	return "user_mobile_user"
}

func readThroughGROM() {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/bacchus?charset=utf8&parseTime=True")
	defer db.Close()
	
	if err != nil {
		log.Fatal(err)
	}
	
	users := make([]User, 10)
	db.Find(&users)
	fmt.Println(len(users))
	
	for _, user := range users {
		fmt.Println(user)
	}
}

func main() {
	readThroughGROM()
}