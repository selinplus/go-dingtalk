package models

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

func AddDevMod(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func ModifyZzrq(devid, t string) error {
	if err := db.Table("devmodify").Where("devid=?", devid).
		Order("id desc").Limit(1).Update("zzrq=?", t).Error; err != nil {
		return err
	}
	return nil
}
