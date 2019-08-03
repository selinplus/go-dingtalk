package models

import "github.com/selinplus/go-dingtalk/pkg/util"

/*消息*/
type Msg struct {
	ID         uint          `gorm:"primary_key;COMMENT:'消息标识';type:int(11);AUTO_INCREMENT"`
	FromName   string        `json:"from_name" gorm:"COMMENT:'发件人名';type:varchar(255)"`
	FromID     string        `json:"from_id" gorm:"COMMENT:'发件人id';type:varchar(100)"`
	ToName     string        `json:"to_name" gorm:"COMMENT:'收件人名';type:varchar(255)"`
	ToID       string        `json:"to_id" gorm:"COMMENT:'收件人id';type:varchar(255)"`
	Title      string        `json:"title" gorm:"COMMENT:'消息标题';type:varchar(255)"`
	Content    string        `json:"content" gorm:"COMMENT:'消息内容';type:varchar(1000)"`
	Time       util.JSONTime `json:"time" gorm:"COMMENT:'发送时间'"`
	FlagNotice uint          `json:"flag_notice" gorm:"COMMENT:'0: 未推送 1: 已推送';type:int(1);default:'0'"`
}
