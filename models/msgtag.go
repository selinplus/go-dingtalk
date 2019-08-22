package models

import (
	"strings"
)

/*消息标签*/
type MsgTag struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Tag      uint   `json:"tag" gorm:"COMMENT:'0：已删除 1: 收件箱 2: 发件箱 3: 草稿箱';type:varchar(255);type:int(11);default:'1'"`
	MsgID    uint   `json:"msg_id" gorm:"COMMENT:'消息标识';type:int(11)"`
	OwnerID  string `json:"owner_id" gorm:"type:varchar(100)"`
	FlagRead uint   `json:"flag_read" gorm:"COMMENT:'0: 未读 1: 已读';type:int(1);default:'0'"`
}

func AddMsgTag(msgID uint, ToID, FromID string) error {
	ToIDs := strings.Split(ToID, ",")
	for _, owner_id := range ToIDs {
		msgtag := MsgTag{OwnerID: owner_id, MsgID: msgID, Tag: 1}
		if err := db.Create(&msgtag).Error; err != nil {
			return err
		}
	}
	FromIDs := strings.Split(FromID, ",")
	for _, owner_id := range FromIDs {
		msgtag := MsgTag{OwnerID: owner_id, MsgID: msgID, Tag: 2}
		if err := db.Create(&msgtag).Error; err != nil {
			return err
		}
	}
	return nil
}
func DeleteMsg(userid string, msg_id, tag uint) error {
	if tag == 0 {
		if err := db.Table("msg_tag").
			Where("owner_id = ? and msg_id = ? and tag=?", userid, msg_id, tag).
			Update("tag", 4).Error; err != nil {
			return err
		}
	} else {
		if err := db.Table("msg_tag").
			Where("owner_id = ? and msg_id = ? and tag=?", userid, msg_id, tag).
			Update("tag", 0).Error; err != nil {
			return err
		}
	}
	return nil
}
