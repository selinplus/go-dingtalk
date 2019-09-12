package models

type Process struct {
	ID     uint   `gorm:"primary_key;COMMENT:'流程实例ID'"`
	Dm     string `json:"dm" gorm:"COMMENT:'提报类型代码'"`
	Tbr    string `json:"tbr" gorm:"COMMENT:'提报人'"`
	Mobile string `json:"mobile" gorm:"COMMENT:'提报人员电话'"`
	DevID  string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"`
	Xq     string `json:"xq" gorm:"COMMENT:'详细描述'"`
	Zp     string `json:"zp" gorm:"COMMENT:'设备照片'"`
	Tbsj   string `json:"tbsj" gorm:"COMMENT:'提报日期'"`
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

type ProcResponse struct {
	ID     uint
	Dm     string `json:"dm"`
	Tbr    string `json:"tbr"`
	Mobile string `json:"mobile"`
	DevID  string `json:"devid"`
	Xq     string `json:"xq"`
	Zp     string `json:"zp"`
	Tbsj   string `json:"tbsj"`
	Node   string `json:"node"`
	Czr    string `json:"czr"`
	Zt     string `json:"zt"`
}

func GetProcDetail(procid uint) (*ProcResponse, error) {
	var pr ProcResponse
	if err := db.Raw("select process.id,proctype.mc as dm,process.tbr,process.mobile,process.devid,process.xq,process.zp,process.tbsj,procmodify.node,user.name as czr from process left join procmodify on process.id=procmodify.procid left join user on user.mobile=procmodify.czr where process.id = ? and (procmodify.czrq ='' or procmodify.czrq is null)", procid).
		Scan(&pr).Error; err != nil {
		return nil, err
	}
	if pr.ID > 0 {
		return &pr, nil
	}
	return nil, nil
}

func GetProcTodoList(czr string) ([]*ProcResponse, error) {
	var pr []*ProcResponse
	if err := db.Raw("select process.id,proctype.mc as dm,process.tbr,process.mobile,process.devid,process.xq,process.zp,process.tbsj,procmodify.node,user.name as czr from process left join procmodify on process.id=procmodify.procid left join user on user.mobile=procmodify.czr where procmodify.czr = ? and (procmodify.czrq ='' or procmodify.czrq is null)", czr).
		Scan(&pr).Error; err != nil {
		return nil, err
	}
	if len(pr) > 0 {
		return pr, nil
	}
	return nil, nil
}

func GetProcDoneList(czr string) ([]*ProcResponse, error) {
	var pr []*ProcResponse
	if err := db.Raw("select process.id,proctype.mc as dm,process.tbr,process.mobile,process.devid,process.xq,process.zp,process.tbsj,procmodify.node,user.name as czr from process left join procmodify on process.id=procmodify.procid left join user on user.mobile=procmodify.czr where procmodify.czr = ? and (procmodify.czrq !='' and procmodify.czrq is not null)", czr).
		Scan(&pr).Error; err != nil {
		return nil, err
	}
	if len(pr) > 0 {
		return pr, nil
	}
	return nil, nil
}
