package models

type Procmodify struct {
	ID         uint   `gorm:"primary_key;AUTO_INCREMENT"`
	ProcID     uint   `json:"procid" gorm:"COMMENT:'流程实例ID';column:procid"`
	Dm         string `json:"dm" gorm:"COMMENT:'提报类型代码'"`
	Node       string `json:"node" gorm:"COMMENT:'当前节点代码'"`
	Tsr        string `json:"tsr" gorm:"COMMENT:'推送人'"`
	Czr        string `json:"czr" gorm:"COMMENT:'操作人'"`
	Spyj       string `json:"spyj" gorm:"COMMENT:'审批意见'"`
	Czrq       string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	FlagNotice int    `json:"flag_notice" gorm:"COMMENT:'0: 未推送 1: 已推送';size:1;default:'0'"`
}

func AddProcMod(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProcMod(pm *Procmodify) error {
	//log.Println(pm)
	if err := db.Table("procmodify").
		Where("id=?", pm.ID).Updates(&pm).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProcessModFlag(ID uint) error {
	if err := db.Table("procmodify").
		Where("id = ?", ID).Update("flag_notice", 0).Error; err != nil {
		return err
	}
	return nil
}

func GetProcMod(id uint) (*Procmodify, error) {
	var pm Procmodify
	if err := db.Where("id=?", id).First(&pm).Error; err != nil {
		return nil, err
	}
	return &pm, nil
}

func GetProcMods(procid uint) ([]*Procmodify, error) {
	var pms []*Procmodify
	query := `
select procmodify.id,
       procmodify.procid,
       proctype.mc as dm,
       procmodify.node,
       user.name   as tsr,
       procmodify.czr,
       procmodify.spyj,
       procmodify.czrq
    from
       procmodify
           left join proctype on procmodify.dm = proctype.dm
           left join user on user.mobile = procmodify.tsr
    where
       procmodify.procid = ?
	order by procmodify.id desc`
	if err := db.Raw(query, procid).Scan(&pms).Error; err != nil {
		return nil, err
	}
	return pms, nil
}

func IsProcManualDone(procid uint) bool {
	var p Procmodify
	if err := db.Where("procid=? and (czrq = '' or czrq is null)", procid).
		First(&p).Error; err != nil {
		return false
	}
	return true
}

func GetProcessFlag() ([]*Procmodify, error) {
	var procmodifies []*Procmodify
	if err := db.Table("procmodify").
		Where("flag_notice=0").Find(&procmodifies).Error; err != nil {
		return nil, err
	}
	return procmodifies, nil
}

func UpdateProcessFlag(ID uint) error {
	if err := db.Table("procmodify").
		Where("id = ? and flag_notice = 0", ID).
		Update("flag_notice", 1).Error; err != nil {
		return err
	}
	return nil
}
