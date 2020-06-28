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

func UserSync(user interface{}) error {
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func CountUserSyncNum(t string) (int, error) {
	var userNum int
	if err := db.Table("user").
		Where("sync_time>=?", t).Count(&userNum).Error; err != nil {
		return 0, err
	}
	return userNum, nil
}

func GetUserByDepartmentID(deptId string) ([]*User, error) {
	var (
		users    []*User
		usersAll []*User
	)
	if err := db.Table("user").
		Where("deptId like ?", "%"+deptId+"%").Find(&usersAll).Error; err != nil {
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

func GetUserByMobile(mobile string) (*User, error) {
	var user User
	if err := db.Table("user").
		Where("mobile=?", mobile).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func IsUserExist(userid string, t string) bool {
	var user User
	if err := db.Select("userid").Where("userid =? and sync_time>=?", userid, t).
		First(&user).Error; err != nil {
		return false
	}
	return true
}

func DeleteUser(userid string) error {
	if err := db.Where("userid=?", userid).Delete(User{}).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByUserid(userid string) (*User, error) {
	var user User
	if err := db.Table("user").
		Where("userid=?", userid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByMc(mc string) ([]*User, error) {
	var us []*User
	err := db.Table("user").
		Select("user.userid,user.name,department.name as deptId,user.mobile").
		Joins("join department on user.deptId=department.id").
		Where("user.name like ? ", "%"+mc+"%").Scan(&us).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return us, nil
}

func CleanUpUser() error {
	err := db.Where("DATEDIFF(NOW(),sync_time)>7").Delete(User{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}
