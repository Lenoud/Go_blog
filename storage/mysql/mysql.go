package mysql

import (
	"fmt"
	"log"

	"blog/config"
	"blog/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.GlobalConfig.MySQL.Username,
		config.GlobalConfig.MySQL.Password,
		config.GlobalConfig.MySQL.Host,
		config.GlobalConfig.MySQL.Port,
		config.GlobalConfig.MySQL.Database,
		config.GlobalConfig.MySQL.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(config.GlobalConfig.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.MySQL.MaxOpenConns)

	// 自动迁移数据库表
	err = DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Post{})
	if err != nil {
		return err
	}

	log.Println("MySQL connected successfully")
	return nil
}
