/*
 * @Author: GG
 * @Date: 2022-11-21 11:48:33
 * @LastEditTime: 2022-11-23 14:54:29
 * @LastEditors: GG
 * @Description:
 * @FilePath: \ginchat\test\testGorm.go
 *
 */
package main

import (
	"ginchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("ginchat:ginchat@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移
	db.AutoMigrate(&models.Message{}, &models.Contact{}, &models.GroupBasic{})

	// Create
	// user := &models.UserBasic{}
	// user.Name = "admin"
	// db.Create(user)

	// // Read
	// fmt.Println(db.First(user, 1))

	// db.Model(user).Update("PassWord", 1234)
}
