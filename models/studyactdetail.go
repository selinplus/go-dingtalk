package models

import "fmt"

//党建活动参与详情
type StudyActdetail struct {
	ID         uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	StudyActID uint   `json:"study_act_id" gorm:"COMMENT:'活动id'"`
	UserID     string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	Bmrq       string `json:"bmrq" gorm:"COMMENT:'报名日期'"`
	Status     string `json:"status" gorm:"COMMENT:'状态,0:未审核 1:审核通过';default:'0'"`
	Sprq       string `json:"sprq" gorm:"COMMENT:'审批日期'"`
	StudyAct   StudyAct
}

func AddStudyActdetail(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdStudyActdetail(activity *StudyActdetail) error {
	if err := db.Table("study_actdetail").
		Where("id=?", activity.ID).Updates(activity).Error; err != nil {
		return err
	}
	return nil
}

func GetApproveStudyActs(cond string) ([]*StudyActdetail, error) {
	var actdetails []*StudyActdetail
	if err := db.Preload("StudyAct").
		Where("status='0'").Where(cond).
		Find(&actdetails).Error; err != nil {
		return nil, err
	}
	return actdetails, nil
}

func GetStudyActdetails(actId uint, status string) ([]*StudyActdetail, error) {
	var actdetails []*StudyActdetail
	if err := db.Preload("StudyAct").
		Where("study_act_id=? and status like ?", actId, status+"%").
		Find(&actdetails).Error; err != nil {
		return nil, err
	}
	return actdetails, nil
}

func IsJoinStrudyAct(actId uint, userid string) string {
	var actdetail StudyActdetail
	if err := db.Where("study_act_id=? and userid = ?", actId, userid).
		First(&actdetail).Error; err != nil {
		return "N"
	}
	return "Y"
}

type CountActResp struct {
	ActID  uint   `json:"act_id"`
	Title  string `json:"title"`
	UserID string `json:"userid"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Dm     string `json:"dm"`
	Mc     string `json:"mc"`
}

func CountStudyAct(actId, dm string) ([]*CountActResp, error) {
	var car []*CountActResp
	sql := fmt.Sprintf(`
select study_act.id, study_act.title, study_actdetail.userid, user.name,user.mobile, study_member.dm, study_group.mc
from study_act, study_actdetail
         left join study_member on study_member.userid = study_actdetail.userid
         left join study_group on study_group.dm = study_member.dm
         left join user on study_member.userid = user.userid
where study_act.id = study_actdetail.study_act_id
  and study_act.id = %s and study_member.dm = '%s'`, actId, dm)
	err := db.Raw(sql).Scan(&car).Error
	if err != nil {
		return nil, err
	}
	return car, nil
}
