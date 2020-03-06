package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"strconv"
)

type Devdept struct {
	ID     uint   `gorm:"primary_key"`
	Jgdm   string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	Jgmc   string `json:"jgmc" gorm:"COMMENT:'设备管理机构名称'"`
	Sjjgdm string `json:"sjjgdm" gorm:"COMMENT:'上级设备管理机构代码'"`
	Gly    string `json:"gly" gorm:"COMMENT:'设备管理员代码'"`
	Lrr    string `json:"lrr" gorm:"COMMENT:'录入人代码'"`
	Lrrq   string `json:"lrrq" gorm:"COMMENT:'录入日期'"`
	Xgr    string `json:"xgr" gorm:"COMMENT:'修改人代码'"`
	Xgrq   string `json:"xgrq" gorm:"COMMENT:'修改日期'"`
}

func GenDevdeptDmBySjjgdm(sjjgdm string) (string, error) {
	var ddt Devdept
	err := db.Table("devdept").
		Where("sjjgdm=?", sjjgdm).Limit(1).Order("id desc").First(&ddt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	if err == gorm.ErrRecordNotFound {
		return sjjgdm + "01", nil
	}
	dm, err := strconv.Atoi(ddt.Jgdm[len(sjjgdm) : len(sjjgdm)+2])
	if err != nil {
		return "", err
	}
	if dm+1 < 10 {
		return sjjgdm + "0" + strconv.Itoa(dm+1), nil
	}
	if dm+1 > 99 {
		return "", errors.New("机构代码超过99")
	}
	return sjjgdm + strconv.Itoa(dm+1), nil
}

func AddDevdept(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateDevdept(devd *Devdept) error {
	if err := db.Table("devdept").
		Where("jgdm=?", devd.Jgdm).Updates(devd).Error; err != nil {
		return err
	}
	return nil
}
func DelDevdeptGly(devd map[string]interface{}) error {
	if err := db.Table("devdept").
		Where("jgdm=?", devd["jgdm"]).Updates(devd).Error; err != nil {
		return err
	}
	return nil
}

func IsSjjg(jgdm string) bool {
	var d Devdept
	if err := db.Where("sjjgdm=? ", jgdm).First(&d).Error; err != nil {
		return false
	}
	return true
}

func IsDevdeptGylExist(jgdm string) bool {
	var du Devuser
	if err := db.Where("jgdm=? and gyl !=''", jgdm).First(&du).Error; err != nil {
		return false
	}
	return true
}

func DeleteDevdept(jgdm string) error {
	if err := db.Where("jgdm=?", jgdm).Delete(Devdept{}).Error; err != nil {
		return err
	}
	return nil
}

func GetDevdept(jgdm string) (*Devdept, error) {
	var dd Devdept
	if err := db.Where("jgdm=?", jgdm).First(&dd).Error; err != nil {
		return nil, err
	}
	return &dd, nil
}

func GetDevGly(gly string) ([]*Devdept, error) {
	var dds []*Devdept
	if err := db.Where("gly=?", gly).Find(&dds).Error; err != nil {
		return nil, err
	}
	return dds, nil
}

func GetDevdeptBySjjgdm(jgdm string) ([]*Devdept, error) {
	var dds []*Devdept
	if err := db.Where("sjjgdm=?", jgdm).Find(&dds).Error; err != nil {
		return nil, err
	}
	return dds, nil
}

func IsLeafDevdept(jgdm string) bool {
	var dt Devdept
	err := db.Select("id").Where("sjjgdm =?", jgdm).First(&dt).Error
	if err == gorm.ErrRecordNotFound {
		return true
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if dt.ID > 0 {
		return false
	}
	return true
}

type ScopedSlots struct {
	Title string `json:"title"`
}

type DevdeptTree struct {
	Jgdm        string `json:"jgdm"`
	Jgmc        string `json:"jgmc"`
	Sjjgdm      string `json:"sjjgdm"`
	Gly         string `json:"gly"`
	Disabled    bool   `json:"disabled"`
	ScopedSlots `json:"scopedSlots"`
	Children    []*DevdeptTree `json:"children"`
}

//获取设备管理机构列表
func GetDevdeptTree(jgdm, bz string) ([]DevdeptTree, error) {
	ytsw, err := GetDevdept(jgdm)
	if err != nil {
		return nil, err
	}
	perms := make([]DevdeptTree, 0)
	child := DevdeptTree{
		Jgdm:        ytsw.Jgdm,
		Jgmc:        ytsw.Jgmc,
		ScopedSlots: ScopedSlots{Title: "custom"},
		Gly:         ytsw.Gly,
		Children:    []*DevdeptTree{},
	}
	if bz == "0" {
		child.Disabled = false
	}
	if bz == "1" {
		if ytsw.Gly == "" {
			child.Disabled = true
		}
	}
	if err := getDevdeptTreeNode(jgdm, bz, &child); err != nil {
		return nil, err
	}
	perms = append(perms, child)
	return perms, nil
}

//递归获取子节点
func getDevdeptTreeNode(sjjgdm, bz string, tree *DevdeptTree) error {
	var perms []*Devdept
	err := db.Where("sjjgdm=?", sjjgdm).Find(&perms).Error //根据父结点Id查询数据表，获取相应的子结点信息
	if err != nil {
		return err
	}
	for i := 0; i < len(perms); i++ {
		child := DevdeptTree{
			Jgdm:        perms[i].Jgdm,
			Jgmc:        perms[i].Jgmc,
			Sjjgdm:      perms[i].Sjjgdm,
			Gly:         perms[i].Gly,
			ScopedSlots: ScopedSlots{Title: "custom"},
			Children:    []*DevdeptTree{},
		}
		if bz == "0" {
			child.Disabled = false
		}
		if bz == "1" {
			if perms[i].Gly == "" {
				child.Disabled = true
			}
		}
		tree.Children = append(tree.Children, &child)
		err = getDevdeptTreeNode(perms[i].Jgdm, bz, &child)
	}
	return err
}

func GetSyrDepts(userid string) ([]*Devdept, error) {
	var syrDepts []*Devdept
	query := `select distinct devuser.jgdm, devdept.jgmc
			from devdept,devuser
			where devdept.jgdm = devuser.jgdm and devuser.syr=?`
	err := db.Raw(query, userid).Find(&syrDepts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return syrDepts, nil
}

func GetGlyDepts(userid string) ([]*Devdept, error) {
	var syrDepts []*Devdept
	err := db.Where("gly=?", userid).Find(&syrDepts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return syrDepts, nil
}
