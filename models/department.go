package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Department struct {
	ID        int    `gorm:"primary_key;COMMENT:'部门id'"`
	Name      string `json:"name" gorm:"COMMENT:'部门名称'"`
	Parentid  int    `json:"parentid" gorm:"COMMENT:'父部门id，根部门为1'"`
	OuterDept bool   `json:"outerDept" gorm:"column:outerDept;COMMENT:'是否本部门的员工仅可见员工自己，为true时，本部门员工默认只能看到员工自己'"`
	SyncTime  string `json:"sync_time" gorm:"COMMENT:'同步时间'"`
}

func DepartmentSync(data interface{}) error {
	if err := db.Model(&Department{}).Save(data).Error; err != nil {
		return err
	}
	return nil
}
func CountDepartmentSyncNum() (int, error) {
	var depidsNum int
	t := time.Now().Format("2006-01-02") + " 00:00:00"
	if err := db.Table("department").Where("sync_time>?", t).Count(&depidsNum).Error; err != nil {
		return 0, err
	}
	return depidsNum, nil
}
func GetDepartmentByParentID(ParentID int) ([]*Department, error) {
	var departments []*Department
	err := db.Table("department").Where("parentid=?", ParentID).Find(&departments).Error
	if err != nil {
		return nil, err
	}
	return departments, nil
}
func GetDepartmentByID(id int) (*Department, error) {
	var dt Department
	if err := db.Table("department").Where("id=?", id).First(&dt).
		Error; err != nil {
		return nil, err
	}
	if dt.ID > 0 {
		return &dt, nil
	}
	return nil, nil
}
func IsDeptIdExist(deptId int) (bool, error) {
	var dt Department
	err := db.Select("id").Where("id =?", deptId).First(&dt).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if dt.ID > 0 {
		return true, nil
	}
	return false, nil
}
func DeleteDepartment(depId int) error {
	if err := db.Where("id=?", depId).Delete(Department{}).Error; err != nil {
		return err
	}
	return nil
}
