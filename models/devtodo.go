package models

import (
	"github.com/jinzhu/gorm"
)

type Devtodo struct {
	ID         uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Czlx       string `json:"czlx" gorm:"COMMENT:'操作类型'"`
	Lsh        string `json:"lsh" gorm:"COMMENT:'流水号'"`
	Czr        string `json:"czr" gorm:"COMMENT:'操作人'"`
	Czrq       string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Jgdm       string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	DevID      string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"` //devinfo ID
	Done       int    `json:"done" gorm:"COMMENT:'0: 待办 1: 已办';size:1;default:'0'"`
	FlagNotice int    `json:"flag_notice" gorm:"COMMENT:'0: 未推送 1: 已推送';size:1;default:'0'"`
}

type DevtodoResp struct {
	*Devtodo
	Gly     string `json:"gly"`
	Zcbh    string `json:"zcbh"`
	Mc      string `json:"mc"`
	Zt      string `json:"zt"`
	Lsh     string `json:"lsh"`
	Czlx    string `json:"czlx"`
	Num     int    `json:"num"`
	SrcJgdm string `json:"src_jgdm"`
	Jgmc    string `json:"jgmc"`
}

func GetDevTodosOrDones(done int) ([]*DevtodoResp, error) {
	var dtos []*DevtodoResp
	err := db.Table("devtodo").
		Select("devtodo.id,devtodo.czlx,devtodo.lsh,user.name as czr,devtodo.czrq,devtodo.jgdm,devdept.gly,devtodo.devid,devinfo.zcbh,devinfo.mc,devinfo.zt,devtodo.done").
		Joins("left join devdept on devdept.jgdm=devtodo.jgdm").
		Joins("left join devinfo on devinfo.id=devtodo.devid").
		Joins("left join user on user.userid=devtodo.czr").
		Where("devtodo.done=?", done).Scan(&dtos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return dtos, nil
}

func GetUpDevTodosOrDones(done int) ([]*DevtodoResp, error) {
	var dtos []*DevtodoResp
	err := db.Table("devtodo").
		Select("devtodo.id,devtodo.czlx,devtodo.lsh,user.name as czr,devtodo.czrq,devtodo.devid,devdept.gly,devtodo.done,devmod.jgdm as src_jgdm,devdept.jgmc,devmod.num").
		Joins("left join user on user.userid=devtodo.czr").
		Joins("left join devmod on devmod.lsh=devtodo.lsh").
		Joins("left join devdept on devdept.jgdm=devmod.jgdm").
		Where("devtodo.done=?", done).Scan(&dtos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return dtos, nil
}

func GetDevFlag() ([]*DevtodoResp, error) {
	var dtos []*DevtodoResp
	if err := db.Table("devtodo").
		Select("devtodo.id,devtodo.czlx,devtodo.czrq,devtodo.jgdm,devdept.gly,devtodo.devid,devtodo.done,devmod.czlx,devmod.num,devmod.jgdm as src_jgdm").
		Joins("left join devdept on devdept.jgdm=devtodo.jgdm").
		Joins("left join devmod on devmod.lsh=devtodo.lsh").
		Where("devtodo.flag_notice=0").Scan(&dtos).Error; err != nil {
		return nil, err
	}
	return dtos, nil
}

func UpdateDevtodoFlag(id uint) error {
	if err := db.Table("devtodo").Where("id = ? and flag_notice = 0", id).
		Update("flag_notice", 1).Error; err != nil {
		return err
	}
	return nil
}
