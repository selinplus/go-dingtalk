package models

import "github.com/jinzhu/gorm"

//党员风采
type StudyHlt struct {
	ID            uint     `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	StudyActID    uint     `json:"act_id" gorm:"column:act_id;COMMENT:'活动id'"`
	UserID        string   `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	Title         string   `json:"title" gorm:"COMMENT:'风采主题'"`
	Content       string   `json:"content" gorm:"COMMENT:'风采内容';size:65535"`
	HltUrl        string   `json:"hlt_url" gorm:"COMMENT:'风采文件真实路径';size:65535"`
	HltUrls       []string `json:"hlt_urls" gorm:"-"` //返回前台[]Url
	Fbrq          string   `json:"fbrq" gorm:"COMMENT:'发布日期'"`
	Xgrq          string   `json:"xgrq" gorm:"COMMENT:'修改日期'"`
	Flag          string   `json:"flag" gorm:"COMMENT:'0:图文 1:视频';default:'0'"`
	Status        string   `json:"status" gorm:"COMMENT:'状态,0:未审核 1:审核通过(发布) 2:撤销发布 3:审核驳回';default:'0'"`
	StarNum       int      `json:"star_num" gorm:"-"` //点赞数
	Star          bool     `json:"star" gorm:"-"`     //点赞标志，仅前台展示，不做数据库存储
	StudyHltStars []StudyHltStar
}

func AddStudyHlt(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdStudyHlt(activity *StudyHlt) error {
	if err := db.Table("study_hlt").
		Where("id=?", activity.ID).Updates(activity).Error; err != nil {
		return err
	}
	return nil
}

func DelStudyHlt(id string) error {
	if err := db.Where("id=?", id).Delete(StudyHlt{}).Error; err != nil {
		return err
	}
	return nil
}

func GetStudyHlt(id string) (*StudyHlt, error) {
	var hlt StudyHlt
	if err := db.
		Preload("StudyHltStars", func(db *gorm.DB) *gorm.DB {
			return db.Order("study_hlt_star.stime")
		}).
		Where("id=?", id).First(&hlt).Error; err != nil {
		return nil, err
	}
	return &hlt, nil
}

func GetStudyHlts(cond string, pageNo, pageSize int) ([]*StudyHlt, error) {
	var hlts []*StudyHlt
	if err := db.
		Preload("StudyHltStars", func(db *gorm.DB) *gorm.DB {
			return db.Order("study_hlt_star.stime")
		}).
		Where(cond).Order("fbrq desc").Limit(pageSize).Offset(pageSize * (pageNo - 1)).
		Find(&hlts).Error; err != nil {
		return nil, err
	}
	return hlts, nil
}

func GetStudyHltsCnt(cond string) (cnt int) {
	err := db.Table("study_hlt").Where(cond).Count(&cnt).Error
	if err != nil {
		cnt = 0
	}
	return cnt
}

func GetStudyHltsByUserid(userid, status, flag string, pageNo, pageSize int) ([]*StudyHlt, error) {
	var hlts []*StudyHlt
	if err := db.
		Preload("StudyHltStars", func(db *gorm.DB) *gorm.DB {
			return db.Order("study_hlt_star.stime")
		}).
		Where("userid=? and status like ? and flag like ?", userid, status+"%", flag+"%").
		Order("fbrq desc").Limit(pageSize).Offset(pageSize * (pageNo - 1)).
		Find(&hlts).Error; err != nil {
		return nil, err
	}
	return hlts, nil
}

func GetStudyHltsCntByUserid(userid, status, flag string) (cnt int) {
	err := db.Table("study_hlt").
		Where("userid=? and status like ? and flag like ?",
			userid, status+"%", flag+"%").Count(&cnt).Error
	if err != nil {
		cnt = 0
	}
	return cnt
}
