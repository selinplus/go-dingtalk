package models

//党员风采点赞
type StudyHltStar struct {
	ID         uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	StudyHltID uint   `json:"study_hlt_id" gorm:"COMMENT:'党员风采id'"`
	UserID     string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	Stime      string `json:"stime" form:"COMMENT:'点赞时间'"`
}

func CreateStudyHltStar(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func IsHltStar(star *StudyHltStar) bool {
	var cnt int
	if err := db.Model(&StudyHltStar{}).Where(star).
		Count(&cnt).Error; err != nil {
		return false
	}
	if cnt > 0 {
		return true
	}
	return false
}

func CancelStudyHltStar(hltId uint, userid string) error {
	if err := db.Where("study_hlt_id=? and userid=?", hltId, userid).
		Delete(&StudyHltStar{}).Error; err != nil {
		return err
	}
	return nil
}

func IsStudyHltStar(hltId uint, userid string) bool {
	var cnt int
	if err := db.Model(&StudyHltStar{}).
		Where("study_hlt_id=? and userid=?", hltId, userid).
		Count(&cnt).Error; err != nil {
		return false
	}
	if cnt > 0 {
		return true
	}
	return false
}
