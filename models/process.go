package models

import "github.com/jinzhu/gorm"

type Process struct {
	ID     uint   `gorm:"primary_key;COMMENT:'流程实例ID'"`
	Dm     string `json:"dm" gorm:"COMMENT:'提报类型代码'"`
	Tbr    string `json:"tbr" gorm:"COMMENT:'提报人'"`
	Mobile string `json:"mobile" gorm:"COMMENT:'提报人员电话'"`
	DevID  string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"`
	Xq     string `json:"xq" gorm:"COMMENT:'详细描述'"`
	Zp     string `json:"zp" gorm:"COMMENT:'设备照片'"`
	Tbsj   string `json:"tbsj" gorm:"COMMENT:'提报日期'"`
	Zfbz   string `json:"zfbz" gorm:"COMMENT:'作废标志,0:未作废,1:作废';default:'0'"`
}

func IsProcExist(data interface{}) (bool, uint, error) {
	var p Process
	err := db.First(&p, data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return false, 0, nil
	}
	return true, p.ID, nil
}

func AddProc(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProc(proc *Process) error {
	if err := db.Table("process").Where("id=?", proc.ID).Updates(proc).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProc(procid uint) error {
	if err := db.Where("id=?", procid).Delete(Process{}).Error; err != nil {
		return err
	}
	return nil
}

type ProcResponse struct {
	ID       uint
	Dm       string `json:"dm"`
	Dmmc     string `json:"dmmc"`
	Tbr      string `json:"tbr"`
	Mobile   string `json:"mobile"`
	Devid    string `json:"devid"`
	Xq       string `json:"xq"`
	Zp       string `json:"zp"`
	Tbsj     string `json:"tbsj"`
	Modifyid uint   `json:"modifyid"`
	Node     string `json:"node"`
	Czr      string `json:"czr"`
	Zt       string `json:"zt"`
}

func GetProcDetail(procid uint) (*ProcResponse, error) {
	var pr ProcResponse
	if err := db.Raw("select process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,process.xq,process.zp,process.tbsj,procmodify.id as modifyid,procmodify.node,user.name as czr from process left join procmodify on process.id=procmodify.procid left join proctype on process.dm=proctype.dm left join user on user.mobile=procmodify.czr where process.id = ?", procid).
		Order("procmodify.id desc").Limit(1).Scan(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func GetProcTodoList(czr string) ([]*ProcResponse, error) {
	var pr []*ProcResponse
	err := db.Raw("select process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,process.xq,process.zp,process.tbsj,procmodify.id as modifyid,procmodify.node,user.name as czr from process left join procmodify on process.id=procmodify.procid left join proctype on process.dm=proctype.dm left join user on user.mobile=procmodify.czr where procmodify.czr = ? and (procmodify.czrq ='' or procmodify.czrq is null)", czr).
		Scan(&pr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return pr, nil
}

func GetProcSaveList(tbr string) ([]*ProcResponse, error) {
	var pr []*ProcResponse
	err := db.Raw("select process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,process.xq,process.zp,process.tbsj,procmodify.id as modifyid,procmodify.node,user.name as czr from process left join procmodify on process.id=procmodify.procid left join proctype on process.dm=proctype.dm left join user on user.mobile=procmodify.czr where process.mobile = ? and (procmodify.node is null or procmodify.node = '')", tbr).
		Scan(&pr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return pr, nil
}

func GetProcDoneList(czr string) ([]*ProcResponse, error) {
	var pr []*ProcResponse
	err := db.Raw("select distinct process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,process.xq,process.zp,process.tbsj from process left join procmodify on process.id=procmodify.procid left join proctype on process.dm=proctype.dm where procmodify.czr = ?", czr).
		Order("process.tbsj desc").Scan(&pr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return pr, nil
}
