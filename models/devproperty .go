package models

type Devproperty struct {
	ID uint   `gorm:"primary_key"`
	Dm string `json:"dm" gorm:"COMMENT:'设备属性代码'"`
	Mc string `json:"mc" gorm:"COMMENT:'设备属性'"`
}

func GetDevproperty() ([]*Devproperty, error) {
	var ds []*Devproperty
	if err := db.Table("devproperty").Find(&ds).Error; err != nil {
		return nil, err
	}
	return ds, nil
}
