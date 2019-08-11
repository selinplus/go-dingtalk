package models

import (
	"github.com/jinzhu/gorm"
)

/*用户*/
type User struct {
	//ID         uint   `gorm:"primary_key"`
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

func UserSync(users *[]User) error {
	for _, user := range *users {
		if user.UserID != "" {
			if err := db.Model(&User{}).Save(user).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
func UserDetailSync(data interface{}) error {
	if err := db.Model(&User{}).Save(data).Error; err != nil {
		return err
	}
	return nil
}
func AddUser(user *User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
func EditUser(user *User) error {
	if err := db.Model(&User{}).Where("userid=?", user.UserID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}
func IsUseridExist(userid string) bool {
	var user User
	err := db.Select("userid").Where("userid = ? ", userid).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if len(user.UserID) > 0 {
		return true
	}
	return false
}
