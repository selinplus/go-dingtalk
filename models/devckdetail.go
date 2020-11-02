package models

type Devckdetail struct {
	ID      uint   `gorm:"primary_key;AUTO_INCREMENT"`
	CheckID uint   `json:"check_id" gorm:"COMMENT:'盘点任务编码'"`
	Cktime  string `json:"cktime" gorm:"COMMENT:'盘点时间'"`
	Pdr     string `json:"pdr" gorm:"COMMENT:'盘点人代码'"`
	CkBz    int    `json:"ck_bz" gorm:"COMMENT:'0: 未盘点 1: 已盘点';size:1;default:'0'"`

	DevinfoID string `json:"devinfo_id" gorm:"COMMENT:'设备编号'"`
	Zcbh      string `json:"zcbh" gorm:"COMMENT:'资产编号'"`
	Lx        string `json:"lx" gorm:"COMMENT:'设备类型'"`
	Mc        string `json:"mc" gorm:"COMMENT:'设备名称'"`
	Xh        string `json:"xh" gorm:"COMMENT:'设备型号'"`
	Xlh       string `json:"xlh" gorm:"COMMENT:'序列号'"`
	Ly        string `json:"ly" gorm:"COMMENT:'设备来源'"`
	Gys       string `json:"gys" gorm:"COMMENT:'供应商'"`
	Jg        string `json:"jg" gorm:"COMMENT:'价格'"`
	Scs       string `json:"scs" gorm:"COMMENT:'生产商'"`
	Scrq      string `json:"scrq" gorm:"COMMENT:'生产日期'"`
	Grrq      string `json:"grrq" gorm:"COMMENT:'购入日期'"`
	Bfnx      string `json:"bfnx" gorm:"COMMENT:'设备报废年限'"`
	Rkrq      string `json:"rkrq" gorm:"COMMENT:'入库日期'"`
	Czr       string `json:"czr" gorm:"COMMENT:'操作人'"`
	Czrq      string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Zt        string `json:"zt" gorm:"COMMENT:'设备状态'"`
	Jgdm      string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	Syr       string `json:"syr" gorm:"COMMENT:'设备使用人代码'"`
	SyrJgdm   string `json:"syr_jgdm" gorm:"COMMENT:'使用人员所在设备管理机构代码'"`
	Cfwz      string `json:"cfwz" gorm:"COMMENT:'存放位置'"`
	Sx        string `json:"sx" gorm:"COMMENT:'设备属性'"`
}

func DevCheck(checkId uint, devinfoId string, ck interface{}) error {
	if err := db.Table("devckdetail").
		Where("check_id=? and devinfo_id=?", checkId, devinfoId).
		Updates(ck).Error; err != nil {
		return err
	}
	return nil
}

func IsChecked(CheckID uint, DevinfoID string) bool {
	var ck Devckdetail
	if err := db.Table("devckdetail").
		Where("check_id=? and devinfo_id=? and ck_bz=1", CheckID, DevinfoID).
		First(&ck).Error; err != nil {
		return false
	}
	return true
}

func CheckSyrSelf(CheckID uint, DevinfoID, syr string) bool {
	var ck Devckdetail
	if err := db.Table("devckdetail").
		Where("check_id=? and devinfo_id=? and syr=?", CheckID, DevinfoID, syr).
		First(&ck).Error; err != nil {
		return false
	}
	return true
}

func GetDevckdetails(cond string, pageNo, pageSize int) ([]*Devckdetail, error) {
	var devckdetails []*Devckdetail
	if err := db.Table("devckdetail").Where(cond).
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&devckdetails).Error; err != nil {
		return nil, err
	}
	return devckdetails, nil
}

func GetDevckdetailsCnt(cond string) (cnt int) {
	err := db.Table("devckdetail").Where(cond).Count(&cnt).Error
	if err != nil {
		cnt = 0
	}
	return cnt
}
