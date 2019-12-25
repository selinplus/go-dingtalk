package models

type Devuser struct {
	ID   uint   `gorm:"primary_key"`
	Jgdm string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	Syr  string `json:"syr" gorm:"COMMENT:'设备使用人代码'"`
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

func GetDevuser(jgdm string) ([]*Devuser, error) {
	var dus []*Devuser
	if err := db.Where("jgdm=?", jgdm).Find(&dus).Error; err != nil {
		return nil, err
	}
	return dus, nil
}

func DeleteDevuser(id uint) error {
	if err := db.Where("id=?", id).Delete(Devuser{}).Error; err != nil {
		return err
	}
	return nil
}

func IsDevuserExist(jgdm, syr string) bool {
	var du Devuser
	if err := db.Where("jgdm=? and syr=?", jgdm, syr).First(&du).Error; err != nil {
		return false
	}
	return true
}

func IsDevdeptUserExist(jgdm string) bool {
	var du Devuser
	if err := db.Where("jgdm=?", jgdm).First(&du).Error; err != nil {
		return false
	}
	return true
}

func CreateDevuser(data interface{}) error {
	tx := db.Begin()
	defer tx.Close()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Create(data).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
