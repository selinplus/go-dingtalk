package models

import (
	"github.com/jinzhu/gorm"
	"strings"
)

type User struct {
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
			if err := db.Model(&User{}).Save(&user).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
func CountUserSyncNum(t string) (int, error) {
	var userNum int
	if err := db.Table("user").Where("sync_time>=?", t).Count(&userNum).Error; err != nil {
		return 0, err
	}
	return userNum, nil
}
func GetUserByDepartmentID(deptId string) ([]*User, error) {
	var (
		users    []*User
		usersAll []*User
	)
	dp := "%" + deptId + "%"
	if err := db.Table("user").Where("deptId like ?", dp).Find(&usersAll).Error; err != nil {
		return nil, err
	}
	for _, user := range usersAll {
		for _, DepartmentId := range strings.Split(user.Department, ",") {
			if DepartmentId == deptId {
				users = append(users, user)
			}
		}
	}
	return users, nil
}
func GetUseridByMobile(mobile string) (string, error) {
	var user User
	if err := db.Table("user").Where("mobile=?", mobile).First(&user).Error; err != nil {
		return "", err
	}
	return user.UserID, nil
}
func GetDeptIdByMobile(mobile string) (string, error) {
	var user User
	if err := db.Table("user").Where("mobile=?", mobile).First(&user).Error; err != nil {
		return "", err
	}
	return user.Department, nil
}
func UserDetailSync(data interface{}) error {
	if err := db.Model(&User{}).Save(&data).Error; err != nil {
		return err
	}
	return nil
}
func IsUseridExist(userid string) (bool, error) {
	var user User
	err := db.Select("userid").Where("userid = ? ", userid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if len(user.UserID) > 0 {
		return true, nil
	}
	return false, nil
}
func DeleteUser(userid string) error {
	if err := db.Where("userid=?", userid).Delete(User{}).Error; err != nil {
		return err
	}
	return nil
}
