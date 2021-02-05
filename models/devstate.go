package models

type Devstate struct {
	ID uint   `gorm:"primary_key"`
	Dm string `json:"dm" gorm:"COMMENT:'设备状态代码'"`
	Mc string `json:"mc" gorm:"COMMENT:'设备状态'"`
}

func GetDevstate() ([]*Devstate, error) {
	var ds []*Devstate
	if err := db.Table("devstate").Find(&ds).Error; err != nil {
		return nil, err
	}
	return ds, nil
}

func GetDevstateByDm(dm string) (*Devstate, error) {
	var dt Devstate
	if err := db.Where("dm=?", dm).First(&dt).Error; err != nil {
		return nil, err
	}
	return &dt, nil
}
