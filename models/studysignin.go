package models

//党建每日签到
type StudySignin struct {
	ID     uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	UserID string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	Qdrq   string `json:"qdrq" gorm:"COMMENT:'签到日期'"`
}

func AddStudySignin(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func IsSinin(signin *StudySignin) bool {
	var cnt int
	if err := db.Model(&StudySignin{}).Where(signin).
		Count(&cnt).Error; err != nil {
		return false
	}
	if cnt > 0 {
		return true
	}
	return false
}

func GetSigninByUserid(userid string) ([]*StudySignin, error) {
	var signins []*StudySignin
	if err := db.Where("userid=?", userid).Find(&signins).Error; err != nil {
		return nil, err
	}
	return signins, nil
}

func GetSigninsByQdrq(qdrq string) ([]*StudySignin, error) {
	var signins []*StudySignin
	if err := db.Where("qdrq=?", qdrq).Find(&signins).Error; err != nil {
		return nil, err
	}
	return signins, nil
}
