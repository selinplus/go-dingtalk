package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

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
	TypeStatus    string   `json:"type_status" gorm:"COMMENT:'自我推荐精选,0:未推选 1:自我推选 2:审核通过';default:'0'"`
	StarNum       int      `json:"star_num" gorm:"-"` //点赞数
	Star          bool     `json:"star" gorm:"-"`     //点赞标志，仅前台展示，不做数据库存储
	Dm            string   `json:"dm" gorm:"-"`       //发布人所在学习小组
	StudyHltStars []StudyHltStar
}

func AddStudyHlt(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdStudyHlt(hlt *StudyHlt) error {
	if err := db.Table("study_hlt").
		Where("id=?", hlt.ID).Updates(hlt).Error; err != nil {
		return err
	}
	return nil
}

func DelStudyHlt(id string) error {
	if err := db.Where("id=?", id).Delete(&StudyHlt{}).Error; err != nil {
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
		Joins(`left join (
    SELECT
        study_hlt.id, star_num
    FROM
        study_hlt,
        ( SELECT study_hlt_id, count( study_hlt_id ) star_num FROM study_hlt_star GROUP BY study_hlt_id) a
    WHERE
            study_hlt.id = a.study_hlt_id ) b on b.id=study_hlt.id`).
		Joins(`left join study_member on study_hlt.userid = study_member.userid `).
		Where(cond).Order("b.star_num desc").Order("fbrq desc").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&hlts).Error; err != nil {
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

//根据活动id和小组代码获取发布风采的人员名单
func GetStudyActHltUsersByStudyDm(actId, dm string) []string {
	type Userid struct {
		Userid string
	}
	var userids []*Userid
	sql := fmt.Sprintf(`
SELECT DISTINCT study_hlt.userid FROM study_hlt
left join study_member on study_member.userid = study_hlt.userid
WHERE study_hlt.status =1 and study_hlt.act_id = %s and study_member.dm = '%s'`, actId, dm)
	if err := db.Raw(sql).Scan(&userids).Error; err != nil {
		return nil
	}
	if len(userids) > 0 {
		var ids []string
		for _, u := range userids {
			ids = append(ids, u.Userid)
		}
		return ids
	}
	return nil
}

//根据活动id和小组代码获取点赞风采的总数
func CountStudyActHltStarsByStudyDm(actId, dm string) int {
	type Cnt struct {
		Cnt int
	}
	var cnt Cnt
	sql := fmt.Sprintf(`
SELECT sum(star_num) cnt FROM study_hlt
left join (
    SELECT study_hlt.id,  star_num
    FROM study_hlt,
        ( SELECT study_hlt_id, count( study_hlt_id ) star_num FROM study_hlt_star GROUP BY study_hlt_id) a
    WHERE
            study_hlt.id = a.study_hlt_id ) b on b.id=study_hlt.id
left join study_member on study_member.userid = study_hlt.userid
WHERE study_hlt.status =1 and study_hlt.act_id = %s and study_member.dm = '%s'`, actId, dm)
	if err := db.Raw(sql).Scan(&cnt).Error; err != nil {
		return 0
	}
	return cnt.Cnt
}
