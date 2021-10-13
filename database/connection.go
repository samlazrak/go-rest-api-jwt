package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jwt/models"
)

// * is pointer not &
var DB *gorm.DB

func Connect() {
	dsn := "root:admin@tcp(127.0.0.1:3306)/gojwt?charset=utf8mb4&parseTime=True&loc=Local"
	c, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	// I believe this works similar to c - references is as one would guess
	// a reference to the ~ while not putting the & does not.
	// pointers? no thats *

	DB = c

	c.AutoMigrate(&models.User{})
}