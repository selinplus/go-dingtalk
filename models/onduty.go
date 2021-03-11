package models

type Onduty struct {
	ID         uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	UserID     string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	Content    string `json:"content" gorm:"COMMENT:'内容';size:65535"`
	Tsrq       string `json:"tsrq" gorm:"COMMENT:'推送日期'"`
	FlagNotice int    `json:"flag_notice" gorm:"COMMENT:'1: 未推送 2: 已推送';size:1;default:'1'"`
}

func AddOnduty(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOndutyFlag(id uint) error {
	if err := db.Table("onduty").
		Where("id = ? and flag_notice = 1", id).
		Update("flag_notice", 2).Error; err != nil {
		return err
	}
	return nil
}

func GetOndutyFlag() ([]*Onduty, error) {
	var ods []*Onduty
	if err := db.Table("onduty").Where("flag_notice=1").
		Find(&ods).Error; err != nil {
		return nil, err
	}
	return ods, nil
}
