package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"log"
)

var db *gorm.DB

// Setup initializes the database instance
func Setup() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	CheckTable()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CheckTable() {
	if !db.HasTable("attachment") {
		db.CreateTable(Attachment{})
	} else {
		db.AutoMigrate(Attachment{})
	}
	if !db.HasTable("msg") {
		db.CreateTable(Msg{})
	} else {
		db.AutoMigrate(Msg{})
	}
	if !db.HasTable("msg_tag") {
		db.CreateTable(MsgTag{})
	} else {
		db.AutoMigrate(MsgTag{})
	}
	if !db.HasTable("department") {
		db.CreateTable(Department{})
	} else {
		db.AutoMigrate(Department{})
	}
	if !db.HasTable("user") {
		db.CreateTable(User{})
	} else {
		db.AutoMigrate(User{})
	}
	if !db.HasTable("device") {
		db.CreateTable(Device{})
	} else {
		db.AutoMigrate(Device{})
	}
	if !db.HasTable("devmodify") {
		db.CreateTable(Devmodify{})
	} else {
		db.AutoMigrate(Devmodify{})
	}
	if !db.HasTable("devoperation") {
		db.CreateTable(Devoperation{})
	} else {
		db.AutoMigrate(Devoperation{})
	}
	if !db.HasTable("devstate") {
		db.CreateTable(Devstate{})
	} else {
		db.AutoMigrate(Devstate{})
	}
	if !db.HasTable("devtype") {
		db.CreateTable(Devtype{})
	} else {
		db.AutoMigrate(Devtype{})
	}
}
