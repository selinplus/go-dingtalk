package models

import "github.com/jinzhu/gorm"

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

func AddSendMsg(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateMsgFlag(msgID uint) error {
	if err := db.Table("msg").
		Where("id = ? and flag_notice = 0", msgID).Update("flag_notice", 1).Error; err != nil {
		return err
	}
	return nil
}

func GetMsgs(userID string, tag uint, pageNum, pageSize int) ([]*Msg, error) {
	var msg []*Msg
	err := db.Raw("SELECT msg.* FROM msg LEFT JOIN msg_tag ON msg.id = msg_tag.msg_id WHERE msg_tag.owner_id = ? AND msg_tag.tag = ? ORDER BY msg.time DESC LIMIT ?,?", userID, tag, pageNum, pageSize).
		Scan(&msg).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return msg, nil
}

func GetMsgByID(id, tag uint, userID string) (*Msg, error) {
	var msg Msg
	if err := db.Preload("Attachments").Find(&msg, "id=?", id).
		Joins("msg_tag ON msg.id = msg_tag.msg_id ").
		Where("msg_tag.owner_id=? and msg_tag.tag=?", userID, tag).
		Error; err != nil {
		return nil, err
	}
	if msg.ID > 0 {
		return &msg, nil
	}
	return nil, nil
}

func GetMsgFlag() ([]*Msg, error) {
	var msgs []*Msg
	if err := db.Table("msg").Where("flag_notice=0").Scan(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}
