package models

import "github.com/jinzhu/gorm"

type Devtype struct {
	ID   uint   `gorm:"primary_key"`
	Dm   string `json:"dm" gorm:"COMMENT:'设备类型代码'"`
	Sjdm string `json:"sjdm" gorm:"COMMENT:'上级设备类型代码'"`
	Mc   string `json:"mc" gorm:"COMMENT:'设备类型'"`
}

func GetDevtype() ([]*Devtype, error) {
	var ds []*Devtype
	if err := db.Table("devtype").Find(&ds).Error; err != nil {
		return nil, err
	}
	return ds, nil
}

func IsDevtypeCorrect(dm string) bool {
	var ds Devtype
	if err := db.Table("devtype").Where("dm=?", dm).First(&ds).Error; err != nil {
		return false
	}
	return true
}

func GetDevtypeByDm(dm string) (*Devtype, error) {
	var dt Devtype
	if err := db.Where("dm=?", dm).First(&dt).Error; err != nil {
		return nil, err
	}
	return &dt, nil
}

func GetDevtypeBySjdm(sjdm string) ([]*Devtype, error) {
	var dts []*Devtype
	if err := db.Where("sjdm=?", sjdm).Find(&dts).Error; err != nil {
		return nil, err
	}
	return dts, nil
}

func IsLeafDevtype(dm string) bool {
	var dt Devtype
	err := db.Select("id").Where("sjdm =?", dm).First(&dt).Error
	if err == gorm.ErrRecordNotFound {
		return true
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if dt.ID > 0 {
		return false
	}
	return true
}
