package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}
type Student struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Age   uint
	Grade string
}

type Account struct {
	ID      uint `gorm:"primaryKey"`
	Balance float64
}

type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID uint
	ToAccountID   uint
	Amount        float64
}

func ConnectDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("SQlite连接失败：", err)
		return nil
	}
	fmt.Println("SQlite连接成功")
	return db
}

func main() {
	//连接数据库
	db := ConnectDb()
	if db == nil {

		fmt.Println("数据库连接失败，程序退出")
		return
	}
	//自动迁移
	db.AutoMigrate(&Account{}, &Transaction{})

	// //插入数据
	// db.Create(&User{Name: "ml", Email: "ml@qq.com"})
	// db.Create(&Student{})
	// var users []User
	// //查询数据
	// db.Debug().Find(&users)
	// fmt.Println("查询成功：", users)
}
