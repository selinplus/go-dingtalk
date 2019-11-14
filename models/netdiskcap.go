package models

import "github.com/jinzhu/gorm"

type NetdiskCap struct {
	UserID   string `json:"userid" gorm:"primary_key;column:userid;COMMENT:'用户标识'"`
	Capacity int    `json:"capacity" gorm:"COMMENT:'网盘容量,单位MB'"`
}

func ModNetdiskCap(userid string, cap int) error {
	nc := NetdiskCap{userid, cap}
	if err := db.Save(&nc).Error; err != nil {
		return err
	}
	return nil
}

func GetNetdiskSpareCap(userid string) (int, error) {
	var nc NetdiskCap
	err := db.First(&nc, "userid=?", userid).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return -1, nil
	}
	return nc.Capacity, nil
}
