package models

import "github.com/jinzhu/gorm"

type Devmodify struct {
	ID    uint   `gorm:"primary_key;AUTO_INCREMENT"`
	DevID string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"`
	Czlx  string `json:"czlx" gorm:"COMMENT:'操作类型'"`
	Sydw  string `json:"sydw" gorm:"COMMENT:'使用单位'"`
	Syks  string `json:"syks" gorm:"COMMENT:'使用科室'"`
	Syr   string `json:"syr" gorm:"COMMENT:'使用人'"`
	Cfwz  string `json:"cfwz" gorm:"COMMENT:'存放位置'"`
	Czrq  string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Czr   string `json:"czr" gorm:"COMMENT:'操作人'"`
	Qsrq  string `json:"qsrq" gorm:"COMMENT:'起始日期'"`
	Zzrq  string `json:"zzrq" gorm:"COMMENT:'终止日期'"`
}

func AddDevModify(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func IsLastModifyZzrqExist(devid string) (bool, error) {
	var dev Devmodify
	err := db.Table("devmodify").Where("devid=?", devid).Order("id desc").Limit(1).First(&dev).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil
}

func ModifyZzrq(devid, t string) error {
	if err := db.Table("devmodify").Where("devid=?", devid).Order("id desc").Limit(1).Update("zzrq", t).Error; err != nil {
		return err
	}
	return nil
}

func GetDevModifes(devid string, pageNo, pageSize int) ([]*Devmodify, error) {
	var devs []*Devmodify
	offset := (pageNo - 1) * pageSize
	query := `select devmodify.id,devmodify.devid,devoperation.mc as czlx,department.name as sydw,
			department.name as syks,devmodify.syr,devmodify.cfwz,devmodify.czrq,devmodify.czr,
			devmodify.qsrq,devmodify.zzrq from devmodify 
			left join device on device.id=devmodify.devid 
			left join devoperation on devoperation.dm=devmodify.czlx 
			left join  department on department.id=devmodify.syks  
			where devmodify.devid = ? order by devmodify.id desc LIMIT ?,?`
	if err := db.Raw(query, devid, offset, pageSize).Scan(&devs).Error; err != nil {
		return nil, err
	}
	return devs, nil
}
