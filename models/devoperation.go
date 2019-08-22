package models

type Devoperation struct {
	ID uint   `gorm:"primary_key;"`
	Dm string `json:"dm" gorm:"COMMENT:'操作类型代码';"`
	Mc string `json:"mc" gorm:"COMMENT:'操作类型';"`
}

func GetDevOp() ([]*Devoperation, error) {
	var ds []*Devoperation
	if err := db.Table("devoperation").Find(&ds).Error; err != nil {
		return nil, err
	}
	return ds, nil
}

func EditDevOp(data interface{}) error {
	if err := db.Model(&Devoperation{}).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
