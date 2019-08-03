package models

/*消息标签*/
type MsgTag struct {
	ID       uint   `gorm:"primary_key;type:int(11);AUTO_INCREMENT"`
	Tag      uint   `json:"tag" gorm:"COMMENT:'0：已删除 1: 收件箱 2: 发件箱 3: 草稿箱';type:varchar(255);type:int(11);default:'1'"`
	MsgID    uint   `json:"msg_id" gorm:"COMMENT:'消息标识';type:int(11)"`
	OwnerID  string `json:"owner_id" gorm:"type:varchar(100)"`
	FlagRead uint   `json:"flag_read" gorm:"COMMENT:'0: 未读 1: 已读';type:int(1);default:'0'"`
}
