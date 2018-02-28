package controllers

import (
	"apilab/models"
	"log"

	"github.com/jinzhu/gorm"
	//连接数据库
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//DB 提供给main函数关闭数据库连接
//
var DB, db *gorm.DB

func init() {
	// initDB()
}

func initDB() {
	var err error
	db, err = gorm.Open(
		"mysql",
		"root:key000000@tcp(192.168.3.3:3306)/evolver?charset=utf8mb4&parseTime=True&loc=Local",
	)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	db.LogMode(true)
	db.DropTableIfExists(&models.User{})
	db.AutoMigrate(&models.User{}, &models.Role{})

	roles := []models.Role{models.Role{RoleName: "supervisor"}, models.Role{RoleName: "staff"}}
	var user models.User
	user.UserName = "admin"
	user.Password = "admin"
	user.Roles = roles
	db.Create(&user)
	//导出DB
	DB = db
}
