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

func GenDevdeptDmBySjjgdm(sjdm string) (string, error) {
	var ddt *Devdept
	err := db.Table("devdept").
		Where("sjdm=?", sjdm).Limit(1).Order("id desc").First(&ddt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	if err == gorm.ErrRecordNotFound {
		return sjdm + "01", nil
	}
	dm, err := strconv.Atoi(ddt.Jgdm[2:4])
	if err != nil {
		return "", err
	}
	if dm+1 < 10 {
		return sjdm + "0" + strconv.Itoa(dm+1), nil
	}
	if dm+1 > 99 {
		return "", errors.New("机构代码超过99")
	}
	return sjdm + strconv.Itoa(dm+1), nil
}

func AddDevdept(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateDevdept(devd *Devdept) error {
	if err := db.Table("devdept").
		Where("id=?", devd.ID).Updates(devd).Error; err != nil {
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
func DeleteDevdept(jgdm string) error {
	if err := db.Where("jgdm=?", jgdm).Delete(Devdept{}).Error; err != nil {
		return err
	}
	return nil
}

type Slots struct {
	Icon string `json:"icon"`
}

type DevdeptTree struct {
	Jgdm     string `json:"jgdm"`
	Jgmc     string `json:"jgmc"`
	Sjjgdm   string `json:"sjjgdm"`
	Gly      string `json:"gly"`
	Slots    `json:"slots"`
	Children []*DevdeptTree `json:"children"`
}

//获取设备管理机构列表
func GetDevdeptTree() ([]DevdeptTree, error) {
	perms := make([]DevdeptTree, 0)
	child := DevdeptTree{
		Jgdm:     "00",
		Jgmc:     "烟台市税务局",
		Slots:    Slots{Icon: "icon"},
		Children: []*DevdeptTree{},
	}
	err := getDevdeptTreeNode("00", &child)
	if err != nil {
		return nil, err
	}
	perms = append(perms, child)
	return perms, nil
}

//递归获取子节点
func getDevdeptTreeNode(sjjgdm string, tree *DevdeptTree) error {
	var perms []*Devdept
	err := db.Where("sjjgdm=?", sjjgdm).Find(&perms).Error //根据父结点Id查询数据表，获取相应的子结点信息
	if err != nil {
		return err
	}
	for i := 0; i < len(perms); i++ {
		child := DevdeptTree{
			Jgdm:     perms[i].Jgdm,
			Jgmc:     perms[i].Jgmc,
			Sjjgdm:   perms[i].Sjjgdm,
			Gly:      perms[i].Gly,
			Slots:    Slots{Icon: "icon"},
			Children: []*DevdeptTree{},
		}
		tree.Children = append(tree.Children, &child)
		err = getDevdeptTreeNode(perms[i].Jgdm, &child)
	}
	return err
}
