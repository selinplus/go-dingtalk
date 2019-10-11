package models

type Token struct {
	ID    uint   `gorm:"primary_key"`
	Token string `json:"token"`
}

func AddToken(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteToken(token string) error {
	if err := db.Where("token=?", token).Delete(Process{}).Error; err != nil {
		return err
	}
	return nil
}

func IsTokenExist(token string) bool {
	var t Token
	if err := db.Where("token =?", token).First(&t).Error; err != nil {
		return false
	}
	return true
}
