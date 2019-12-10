package models

type Proctype struct {
	ID      uint   `gorm:"primary_key"`
	Dm      string `json:"dm" gorm:"COMMENT:'提报类型代码'"`
	Mc      string `json:"mc" gorm:"COMMENT:'提报类型'"`
	Checked string `json:"checked" gorm:"COMMENT:'选中'"`
}

//信息中心用，含手工提报
func GetProctypeAll() ([]*Proctype, error) {
	var pt []*Proctype
	if err := db.Table("proctype").Find(&pt).Error; err != nil {
		return nil, err
	}
	return pt, nil
}

//返回除手工提报类型外记录
func GetProctype() ([]*Proctype, error) {
	var pt []*Proctype
	if err := db.Table("proctype").Not("dm", "0").Find(&pt).Error; err != nil {
		return nil, err
	}
	return pt, nil
}
