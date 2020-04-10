package models

import (
	"github.com/jinzhu/gorm"
)

type MsgAddressbook struct {
	ID     uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Name   string `json:"name" gorm:"COMMENT:'通讯组名称'"`
	UserID string `json:"userid" gorm:"column:userid;COMMENT:'通讯组所属人标识'"`
}

func AddAddressbook(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAddressbook(id uint, userID string) error {
	if err := db.Where("id=? and userid=?", id, userID).
		Delete(MsgAddressbook{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateAddressbook(book *MsgAddressbook) error {
	if err := db.Table("msg_addressbook").
		Where("id = ?", book.ID).Updates(book).Error; err != nil {
		return err
	}
	return nil
}

func GetAddressbooks(userID string) ([]*MsgAddressbook, error) {
	var book []*MsgAddressbook
	err := db.Where("userid=?", userID).Find(&book).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return book, nil
}
