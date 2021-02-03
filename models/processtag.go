package models

type ProcessTag struct {
	ID         uint `gorm:"primary_key;AUTO_INCREMENT"`
	ProcID     uint `json:"proc_id" gorm:"COMMENT:'流程实例ID'"`
	FlagNotice int  `json:"flag_notice" gorm:"COMMENT:'0: 未推送 1: 已推送';size:1;default:'0'"`
}

func AddProcessBcms(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProcessBcmsFlag(id uint) error {
	if err := db.Table("process_tag").
		Where("id = ? and flag_notice = 0", id).Update("flag_notice", 1).Error; err != nil {
		return err
	}
	return nil
}

func GetProcessBcmsFlag() ([]*ProcessTag, error) {
	var msgs []*ProcessTag
	if err := db.Table("process_tag").Where("flag_notice=0").Scan(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}
