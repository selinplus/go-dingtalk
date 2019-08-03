package models

import "github.com/selinplus/go-dingtalk/pkg/util"

/*附件*/
type Attachment struct {
	ID       uint          `gorm:"primary_key;COMMENT:'文件标识';type:int(11);AUTO_INCREMENT"`
	MsgID    uint          `json:"msg_id" gorm:"COMMENT:'消息标识';type:int(11)"`
	FileName string        `json:"file_name" gorm:"COMMENT:'文件原始名';type:varchar(255)"`
	FileUrl  string        `json:"file_url" gorm:"COMMENT:'文件真实路径';type:varchar(255)"`
	FileSize string        `json:"file_size" gorm:"COMMENT:'文件大小';type:varchar(20)"`
	FileType string        `json:"file_type" gorm:"COMMENT:'文件类型';type:varchar(20)"`
	Time     util.JSONTime `json:"time" gorm:"COMMENT:'插入时间'"`
}
