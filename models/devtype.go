package models

type Devtype struct {
	ID uint   `gorm:"primary_key;"`
	Dm string `json:"dm" gorm:"COMMENT:'设备类型代码';"`
	Mc string `json:"mc" gorm:"COMMENT:'设备类型';"`
}

func GetDevtype() ([]*Devtype, error) {
	var ds []*Devtype
	if err := db.Table("devtype").Find(&ds).Error; err != nil {
		return nil, err
	}
	return ds, nil
}

func EditDevtype(data interface{}) error {
	if err := db.Model(&Devtype{}).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
