package models

import "github.com/jinzhu/gorm"

type Ydksworkrecord struct {
	ID         uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Lb         string `json:"lb" gorm:"COMMENT:'业务类别'"`
	Req        string `json:"req" gorm:"COMMENT:'钉钉推送请求json';size:65535"`
	Crrq       string `json:"crrq" gorm:"COMMENT:'插入日期'"`
	Tsrq       string `json:"tsrq" gorm:"COMMENT:'推送日期'"`
	FlagNotice int    `json:"flag_notice" gorm:"COMMENT:'1: 未推送 2: 已推送';size:1;default:'1'"`
}

func AddWorkrecord(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GetWorkrecordFlag() ([]*Ydksworkrecord, error) {
	var ods []*Ydksworkrecord
	if err := db.Table("ydksworkrecord").Where("flag_notice=1").Find(&ods).Error; err != nil {
		return nil, err
	}
	return ods, nil
}

func GetWorkrecordSendCnt() (cnt int) {
	if err := db.Table("ydksworkrecord").Where("flag_notice=2").Count(&cnt).Error; err != nil {
		return 0
	}
	return cnt
}

func UpdateWorkrecordFlag(id uint, upd map[string]interface{}) error {
	if err := db.Table("ydksworkrecord").
		Where("id = ? and flag_notice = 1", id).Updates(upd).Error; err != nil {
		return err
	}
	return nil
}

func GetYtstworkrecords(rq string, flag int) ([]*Ydksworkrecord, error) {
	var records []*Ydksworkrecord
	err := db.Table("ydksworkrecord").
		Where("crrq like ? and flag_notice=?", rq, flag).
		Order("tsrq desc").Find(&records).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return records, nil
}
