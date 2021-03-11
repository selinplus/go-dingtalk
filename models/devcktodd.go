package models

import (
	"fmt"
	"log"
)

/*推送钉钉盘点任务*/
type Devcktodd struct {
	ID         uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Jsr        string `json:"jsr" gorm:"COMMENT:'消息接收人'"`
	DevcheckID uint   `json:"devcheck_id" gorm:"COMMENT:'盘点任务编码'"`
	Devcheck   Devcheck
	FlagNotice int `json:"flag_notice" gorm:"COMMENT:'0: 未推送 1: 已推送';size:1;default:'0'"`
}

func GetDevCkTaskFlag() ([]*Devcktodd, error) {
	var devcktodds []*Devcktodd
	if err := db.Table("devcktodd").
		Preload("Devcheck").
		Where("flag_notice=0").Find(&devcktodds).Error; err != nil {
		return nil, err
	}
	return devcktodds, nil
}

func UpdateDevCkTaskFlag(id uint) error {
	if err := db.Table("devcktodd").
		Where("id = ? and flag_notice = 0", id).
		Update("flag_notice", 1).Error; err != nil {
		return err
	}
	return nil
}

func AddSendDevCkTasks(checkId uint, ckBz string) {
	var jsrs = make([]string, 0)
	if ckBz == "Y" { //自我盘点，通知所有使用人
		sql1 := fmt.Sprintf(
			`select DISTINCT syr from devckdetail where syr!='' and check_id=%d`, checkId)
		syrs, err := QueryData(sql1)
		if err != nil {
			log.Println(err, "=====", sql1)
			return
		}
		for _, syr := range syrs {
			jsrs = append(jsrs, syr["syr"])
		}
	} else if ckBz == "N" { //非自我盘点，通知所有管理员
		sql2 := `select DISTINCT gly from devdept where gly!=''`
		glys, err := QueryData(sql2)
		if err != nil {
			log.Println(err, "=====", sql2)
			return
		}
		for _, gly := range glys {
			jsrs = append(jsrs, gly["gly"])
		}
	}
	for _, jsr := range jsrs {
		if err := db.Create(&Devcktodd{Jsr: jsr, DevcheckID: checkId}).Error; err != nil {
			log.Println(err, "=====", checkId, "=====", jsr)
		}
	}
}
