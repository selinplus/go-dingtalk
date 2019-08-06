package models

import (
	"github.com/jinzhu/gorm"
)

/*消息*/
type Msg struct {
	ID          uint   `gorm:"primary_key;COMMENT:'消息标识';size:11;AUTO_INCREMENT"`
	FromName    string `json:"from_name" gorm:"COMMENT:'发件人名'"`
	FromID      string `json:"from_id" gorm:"COMMENT:'发件人id'"`
	ToName      string `json:"to_name" gorm:"COMMENT:'收件人名'"`
	ToID        string `json:"to_id" gorm:"COMMENT:'收件人id'"`
	Title       string `json:"title" gorm:"COMMENT:'消息标题'"`
	Content     string `json:"content" gorm:"COMMENT:'消息内容';size:1000"`
	Time        string `json:"time" gorm:"COMMENT:'发送时间'"`
	FlagNotice  int    `json:"flag_notice" gorm:"COMMENT:'0: 未推送 1: 已推送';size:1;default:'0'"`
	Attachments []Attachment
}

func AddMsgSend(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GetMsgs(userID, tag uint, pageNum, pageSize int) ([]*Msg, error) {
	var mgs []*Msg
	err := db.Table("msg").
		Select("msg.*, msg_tag.tag").
		Joins("msg_tag ON msg.id = msg_tag.msg_id ").
		Where("msg_tag.owner_id=? and msg_tag.tag=?", userID, tag).
		Order("msg.time").
		Offset(pageNum).Limit(pageSize).Find(mgs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return mgs, nil
}
func GetMsgCount(userID, tag uint) (int, error) {
	var cnt int
	err := db.Table("msg").
		Select("msg.*, msg_tag.tag").
		Joins("msg_tag ON msg.id = msg_tag.msg_id ").
		Where("msg_tag.owner_id=? and msg_tag.tag=?", userID, tag).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return cnt, nil
}
func GetMsgByID(id, userID, tag uint) (*Msg, error) {
	var msg Msg
	if err := db.Preload("Attachments").
		Table("msg").
		Select("msg.*, msg_tag.tag").
		Joins("msg_tag ON msg.id = msg_tag.msg_id ").
		Where("msg.id = ? and msg_tag.owner_id=? and msg_tag.tag=?", id, userID, tag).
		Find(&msg).Error; err != nil {
		return nil, err
	}
	if msg.ID > 0 {
		return &msg, nil
	}
	return nil, nil
}
