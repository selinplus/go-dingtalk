package models

import "github.com/jinzhu/gorm"

type Ydksdata struct {
	ID   uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Lb   string `json:"lb" gorm:"COMMENT:'业务类别'"`
	Data string `json:"data" gorm:"COMMENT:'json数据';size:65535"`
	Rq   string `json:"rq" gorm:"COMMENT:'插入日期'"`
}

func AddYdksdata(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GetYdksdata(rq, lb string) ([]*Ydksdata, error) {
	var list []*Ydksdata
	err := db.Table("ydksdata").Where("rq like ? and lb=?", rq+"%", lb).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return list, nil
}
