package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"log"
	"strconv"
)

type StudyGroup struct {
	ID   uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Dm   string `json:"dm" gorm:"COMMENT:'学习小组代码'"`
	Mc   string `json:"mc" gorm:"COMMENT:'学习小组名称'"`
	Sjdm string `json:"sjdm" gorm:"COMMENT:'上级学习小组代码'"`
	Gly  string `json:"gly" gorm:"COMMENT:'管理员代码'"`
	Lrr  string `json:"lrr" gorm:"COMMENT:'录入人代码'"`
	Lrrq string `json:"lrrq" gorm:"COMMENT:'录入日期'"`
	Xgr  string `json:"xgr" gorm:"COMMENT:'修改人代码'"`
	Xgrq string `json:"xgrq" gorm:"COMMENT:'修改日期'"`
}

//根据sjdm生成小组dm
func GenGroupDmBySjjgdm(sjdm string) (string, error) {
	var group StudyGroup
	err := db.Table("study_group").
		Where("sjdm=?", sjdm).Limit(1).Order("id desc").First(&group).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	if err == gorm.ErrRecordNotFound {
		return sjdm + "01", nil
	}
	dm, err := strconv.Atoi(group.Dm[len(sjdm) : len(sjdm)+2])
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

func AddStudyGroup(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdStudyGroup(devd *StudyGroup) error {
	if err := db.Table("study_group").
		Where("dm=?", devd.Dm).Updates(devd).Error; err != nil {
		return err
	}
	return nil
}

func DelStudyGroupGly(devd map[string]interface{}) error {
	if err := db.Table("study_group").
		Where("dm=?", devd["dm"]).Updates(devd).Error; err != nil {
		return err
	}
	return nil
}

func IsStudySjjg(dm string) bool {
	var d StudyGroup
	if err := db.Where("sjdm=? ", dm).First(&d).Error; err != nil {
		return false
	}
	return true
}

func DelStudyGroup(dm string) error {
	if err := db.Where("dm=?", dm).Delete(StudyGroup{}).Error; err != nil {
		return err
	}
	return nil
}

func GetStudyGroup(dm string) (*StudyGroup, error) {
	var dd StudyGroup
	if err := db.Where("dm=?", dm).First(&dd).Error; err != nil {
		return nil, err
	}
	return &dd, nil
}

func GetStudyGroupBySjdm(dm string) ([]*StudyGroup, error) {
	var groups []*StudyGroup
	if err := db.Where("sjdm=?", dm).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

type StudyGroupTree struct {
	Dm          string `json:"dm"`
	Mc          string `json:"mc"`
	Sjdm        string `json:"sjdm"`
	Gly         string `json:"gly"`
	Disabled    bool   `json:"disabled"`
	ScopedSlots `json:"scopedSlots"`
	Children    []*StudyGroupTree `json:"children"`
}

//获取学习小组列表
func GetStudyGroupTree(dm string) ([]StudyGroupTree, error) {
	ytsw, err := GetStudyGroup(dm)
	if err != nil {
		return nil, err
	}
	perms := make([]StudyGroupTree, 0)
	ytsw.Gly = "超级管理员"
	child := StudyGroupTree{
		Dm:          ytsw.Dm,
		Mc:          ytsw.Mc,
		ScopedSlots: ScopedSlots{Title: "custom"},
		Gly:         ytsw.Gly,
		Children:    []*StudyGroupTree{},
	}
	if err := getStudyGroupTreeNode(dm, &child); err != nil {
		return nil, err
	}
	perms = append(perms, child)
	return perms, nil
}

//递归获取子节点
func getStudyGroupTreeNode(sjdm string, tree *StudyGroupTree) error {
	var perms []*StudyGroup
	err := db.Where("sjdm=?", sjdm).Find(&perms).Error //根据父结点Id查询数据表，获取相应的子结点信息
	if err != nil {
		return err
	}
	for i := 0; i < len(perms); i++ {
		if perms[i].Gly != "" {
			u, _ := GetUserByUserid(perms[i].Gly)
			perms[i].Gly = u.Name
		}
		child := StudyGroupTree{
			Dm:          perms[i].Dm,
			Mc:          perms[i].Mc,
			Sjdm:        perms[i].Sjdm,
			Gly:         perms[i].Gly,
			ScopedSlots: ScopedSlots{Title: "custom"},
			Children:    []*StudyGroupTree{},
		}
		tree.Children = append(tree.Children, &child)
		err = getStudyGroupTreeNode(perms[i].Dm, &child)
	}
	return err
}

func initGrouproot() {
	err := db.Create(&StudyGroup{Dm: "00", Mc: "福山品牌党建", Gly: "fsdj_admin"}).Error
	if err != nil {
		log.Println("initGrouproot err:", err)
		return
	}
}
