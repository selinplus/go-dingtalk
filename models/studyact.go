package models

import (
	"github.com/jinzhu/gorm"
)

//党建活动
type StudyAct struct {
	ID              uint     `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	TopicImage      string   `json:"topic_image" gorm:"COMMENT:'主题图片'"`
	Title           string   `json:"title" gorm:"COMMENT:'活动主题'"`
	Content         string   `json:"content" gorm:"COMMENT:'活动内容';size:65535"`
	ImageUrl        string   `json:"image_url" gorm:"COMMENT:'活动图片真实路径';size:65535"`
	ImageUrls       []string `json:"image_urls" gorm:"-"` //返回前台[]Url
	Fbrq            string   `json:"fbrq" gorm:"COMMENT:'发布日期'"`
	Xgrq            string   `json:"xgrq" gorm:"COMMENT:'修改日期'"`
	Deadline        string   `json:"deadline" gorm:"COMMENT:'活动期限'"`
	Share           string   `json:"share" gorm:"COMMENT:'分享标志,0:否 1:分享';default:'0'"`
	Status          string   `json:"status" gorm:"COMMENT:'状态,0:未审核 1:审核通过(发布) 2:撤销发布';default:'0'"`
	Type            string   `json:"type" gorm:"COMMENT:'状态,1:普通活动 2:学习笔记';default:'0'"`
	JoinNum         int      `json:"join_num" gorm:"-"` //参加人数
	Joined          bool     `json:"joined" gorm:"-"`   //参与标志
	StudyActdetails []StudyActdetail
	StudyHlts       []StudyHlt
}

func AddStudyAct(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdStudyAct(activity *StudyAct) error {
	if err := db.Table("study_act").
		Where("id=?", activity.ID).Updates(activity).Error; err != nil {
		return err
	}
	return nil
}

func DelStudyAct(id string) error {
	if err := db.Where("id=?", id).Delete(StudyAct{}).Error; err != nil {
		return err
	}
	return nil
}

func GetStudyAct(id string) (*StudyAct, error) {
	var activity StudyAct
	if err := db.
		Preload("StudyActdetails", func(db *gorm.DB) *gorm.DB {
			return db.Where("study_actdetail.status='1'").
				Order("study_actdetail.bmrq")
		}).
		Preload("StudyHlts", func(db *gorm.DB) *gorm.DB {
			return db.Where("study_hlt.status='1'").
				Order("study_hlt.fbrq")
		}).
		Where("id=?", id).First(&activity).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

func GetStudyActs(share, status, tp, deadline string, pageNo, pageSize int) ([]*StudyAct, error) {
	var acts []*StudyAct
	err := db.
		Preload("StudyActdetails", func(db *gorm.DB) *gorm.DB {
			return db.Where("study_actdetail.status='1'").
				Order("study_actdetail.bmrq")
		}).
		Preload("StudyHlts", func(db *gorm.DB) *gorm.DB {
			return db.Where("study_hlt.status='1'").
				Order("study_hlt.fbrq")
		}).
		Preload("StudyHlts.StudyHltStars", func(db *gorm.DB) *gorm.DB {
			return db.Order("study_hlt_star.stime")
		}).
		Where("share like ? and status like ? and type like ?",
			share+"%", status+"%", tp+"%").Where(deadline).
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&acts).Error
	if err != nil {
		return nil, err
	}
	return acts, nil
}

func GetStudyActsCnt(share, status, tp, deadline string) (cnt int) {
	err := db.Model(&StudyAct{}).
		Where("share like ? and status like ? and type like ?",
			share+"%", status+"%", tp+"%").
		Where(deadline).Count(&cnt).Error
	if err != nil {
		cnt = 0
	}
	return cnt
}
