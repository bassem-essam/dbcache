package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var bucket string = "default"

var dbname = "db_cache_test"
var username = "username"
var password = "password"

func init_db() {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, dbname)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db.AutoMigrate(&Number{})
}

func loadcache() {
	var numbers []Number
	db.Where("bucket = ?", bucket).Find(&numbers)
	for _, n := range numbers {
		num := new(Number)
		*num = n
		fmt.Println("NumberCache.LoadCache(", n.Value, ", ", num, ")")
		NumberCache.Store(n.Value, num)
	}

	fmt.Println("NumberCache loaded for bucket: ", bucket)
}
