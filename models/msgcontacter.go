package models

import (
	"github.com/jinzhu/gorm"
)

type MsgContacter struct {
	ID       uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	BookeID  uint   `json:"bookid" gorm:"column:bookid;COMMENT:'通讯组ID'"`
	UserID   string `json:"userid" gorm:"column:userid;COMMENT:'联系人标识'"`
	DeptName string `json:"deptname" gorm:"column:deptname;COMMENT:'联系人所属部门名称'"`
}

func AddContacter(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteContacter(userid string, bookid uint) error {
	if err := db.Where("userid=? and bookid=?", userid, bookid).
		Delete(&MsgContacter{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteContacters(bookid uint) error {
	if err := db.Where("bookid=?", bookid).
		Delete(&MsgContacter{}).Error; err != nil {
		return err
	}
	return nil
}

func GetContacters(bookid uint) ([]*MsgContacter, error) {
	var msg []*MsgContacter
	err := db.Where("bookid=?", bookid).Find(&msg).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return msg, nil
}
