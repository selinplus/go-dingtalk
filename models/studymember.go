package models

type StudyMember struct {
	ID     uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Dm     string `json:"dm" gorm:"COMMENT:'学习小组代码'"`
	UserID string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
}

func AddStudyMember(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func IsMemberExist(userid string) bool {
	var u StudyMember
	err := db.Where("userid=?", userid).First(&u).Error
	if err != nil {
		return false
	}
	return true
}

func UpdStudyMember(u *StudyMember) error {
	if err := db.Table("study_member").
		Where("id=?", u.ID).Updates(u).Error; err != nil {
		return err
	}
	return nil
}

func GetStudyMembers(dm string) ([]*StudyMember, error) {
	var dus []*StudyMember
	if err := db.Where("dm=?", dm).Find(&dus).Error; err != nil {
		return nil, err
	}
	return dus, nil
}

func DelStudyMember(id uint) error {
	if err := db.Where("id=?", id).Delete(StudyMember{}).Error; err != nil {
		return err
	}
	return nil
}

func IsStudyMemberExist(dm, userid string) bool {
	var du StudyMember
	if err := db.Where("dm=? and userid=?", dm, userid).First(&du).Error; err != nil {
		return false
	}
	return true
}

func IsNullStudyGroup(dm string) bool {
	var du StudyMember
	if err := db.Where("dm=?", dm).First(&du).Error; err != nil {
		return false
	}
	return true
}
