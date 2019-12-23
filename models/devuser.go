package models

type Devuser struct {
	ID  uint   `gorm:"primary_key"`
	Dm  string `json:"dm" gorm:"COMMENT:'设备管理机构代码'"`
	Syr string `json:"glr" gorm:"COMMENT:'设备使用人代码'"`
}

func AddDevuser(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateDevuser(devu *Devuser) error {
	if err := db.Table("devuser").Where("id=?", devu.ID).Updates(devu).Error; err != nil {
		return err
	}
	return nil
}

func DeleteDevuser(id uint) error {
	if err := db.Where("id=?", id).Delete(Devuser{}).Error; err != nil {
		return err
	}
	return nil
}

//TODO
func IsDevuserExist(userid string, id int) bool {
	var nt Devuser
	if err := db.
		Where("userid =? and id=?", userid, id).First(&nt).Error; err != nil {
		return false
	}
	return true
}
