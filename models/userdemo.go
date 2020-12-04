package models

// 设备管理用户信息参照表
type Userdemo struct {
	UserID     string `json:"userid" gorm:"primary_key;column:userid;COMMENT:'用户标识'"`
	Name       string `json:"name" gorm:"COMMENT:'名称'"`
	Department string `json:"deptId" gorm:"column:deptId;COMMENT:'部门id'"`
	Mobile     string `json:"mobile" gorm:"COMMENT:'手机号'"`
	IsAdmin    bool   `json:"isAdmin" gorm:"column:IsAdmin;COMMENT:'是否是企业的管理员，true表示是，false表示不是'"`
	Active     bool   `json:"active" gorm:"COMMENT:'是否激活'"`
	Avatar     string `json:"avatar" gorm:"COMMENT:'头像url'"`
	Remark     string `json:"remark" gorm:"COMMENT:'备注'"`
	SyncTime   string `json:"sync_time" gorm:"COMMENT:'同步时间'"`
}

func SaveUserdemo(userdemo interface{}) error {
	if err := db.Save(userdemo).Error; err != nil {
		return err
	}
	return nil
}

func GetUserdemoByMobile(mobile string) (*Userdemo, error) {
	var userdemo Userdemo
	if err := db.Table("userdemo").
		Where("mobile=?", mobile).First(&userdemo).Error; err != nil {
		return nil, err
	}
	return &userdemo, nil
}

func GetUserdemoByUserid(userid string) (*Userdemo, error) {
	var userdemo Userdemo
	if err := db.Table("userdemo").
		Where("userid=?", userid).First(&userdemo).Error; err != nil {
		return nil, err
	}
	return &userdemo, nil
}
