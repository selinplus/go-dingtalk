package models

import "github.com/jinzhu/gorm"

type Process struct {
	ID      uint   `gorm:"primary_key;COMMENT:'流程实例ID'"`
	Dm      string `json:"dm" gorm:"COMMENT:'提报类型代码'"`
	Tbr     string `json:"tbr" gorm:"COMMENT:'提报人'"`
	Mobile  string `json:"mobile" gorm:"COMMENT:'提报人员电话'"`
	DevID   string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"`
	Zp      string `json:"zp" gorm:"COMMENT:'设备照片'"`
	SyrName string `json:"syr_name" gorm:"COMMENT:'使用人姓名'"`
	Syr     string `json:"syr" gorm:"COMMENT:'使用人手机号'"`
	Cfwz    string `json:"cfwz" gorm:"COMMENT:'存放位置'"`
	Title   string `json:"title" gorm:"COMMENT:'提报事项'"`
	Xq      string `json:"xq" gorm:"COMMENT:'详细描述'"`
	Tbsj    string `json:"tbsj" gorm:"COMMENT:'提报日期'"`
	Zfbz    string `json:"zfbz" gorm:"COMMENT:'作废标志,0:未作废,1:作废';default:'0'"`
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
	SyrName  string `json:"syr_name"`
	Syr      string `json:"syr"`
	Cfwz     string `json:"cfwz"`
	Zp       string `json:"zp"`
	Title    string `json:"title"`
	Xq       string `json:"xq"`
	Tbsj     string `json:"tbsj"`
	Modifyid uint   `json:"modifyid"`
	Node     string `json:"node"`
	Czr      string `json:"czr"`
	Zt       string `json:"zt"`
}

func GetProcDetail(procid uint) (*ProcResponse, error) {
	var pr ProcResponse
	sql := `select process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,process.title,
       			   process.xq,process.zp,process.tbsj,procmodify.id as modifyid,procmodify.node,user.name as czr,
			       process.syr_name,process.syr,process.cfwz
			from process
         	left join procmodify on process.id = procmodify.procid
         	left join proctype on process.dm = proctype.dm
         	left join user on user.mobile = procmodify.czr
			where process.id = ? order by procmodify.id desc`
	if err := db.Raw(sql, procid).Limit(1).Scan(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func GetProcTodoList(czr string) ([]*ProcResponse, error) {
	var pr []*ProcResponse
	sql := `select process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,process.title,
       			   process.xq,process.zp,process.tbsj,procmodify.id as modifyid,procmodify.node,user.name as czr,
				   process.syr_name,process.syr,process.cfwz
			from process
         	left join procmodify on process.id = procmodify.procid
         	left join proctype on process.dm = proctype.dm
         	left join user on user.mobile = procmodify.czr
			where procmodify.czr = ?  and (procmodify.czrq = '' or procmodify.czrq is null)
 			order by process.id desc`
	err := db.Raw(sql, czr).Scan(&pr).Error
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
	sql := `select process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,process.title,
       			   process.xq,process.zp,process.tbsj,procmodify.id as modifyid,procmodify.node,user.name as czr,
				   process.syr_name,process.syr,process.cfwz
			from process
         	left join procmodify on process.id = procmodify.procid
         	left join proctype on process.dm = proctype.dm
         	left join user on user.mobile = procmodify.czr
			where process.mobile = ?  and (procmodify.node is null or procmodify.node = '')`
	err := db.Raw(sql, tbr).Scan(&pr).Error
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
	sql := `select distinct process.id,process.dm,proctype.mc as dmmc,process.tbr,process.mobile,process.devid,
			       process.title,process.xq,process.xq,process.zp,process.tbsj,process.syr_name,process.syr,process.cfwz
		  	from process
          	left join procmodify on process.id = procmodify.procid
          	left join proctype on process.dm = proctype.dm
 		  	where procmodify.czr = ?  and procmodify.czrq !='' and procmodify.czrq is not null 
			order by process.id desc`
	err := db.Raw(sql, czr).Scan(&pr).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return pr, nil
}
