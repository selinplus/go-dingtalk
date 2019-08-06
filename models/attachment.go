package models

/*附件*/
type Attachment struct {
	ID       uint   `gorm:"primary_key;COMMENT:'文件标识';size:11;AUTO_INCREMENT"`
	MsgID    uint   `json:"msg_id" gorm:"COMMENT:'消息标识';size:11"`
	FileName string `json:"file_name" gorm:"COMMENT:'文件原始名'"`
	FileUrl  string `json:"file_url" gorm:"COMMENT:'文件真实路径'"`
	FileSize int    `json:"file_size" gorm:"COMMENT:'文件大小';size:20"`
	FileType string `json:"file_type" gorm:"COMMENT:'文件类型';size:20"`
	Time     string `json:"time" gorm:"COMMENT:'插入时间'"`
}
