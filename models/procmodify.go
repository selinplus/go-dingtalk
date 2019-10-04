package models

type Procmodify struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT"`
	ProcID uint   `json:"procid" gorm:"COMMENT:'流程实例ID';column:procid"`
	Dm     string `json:"dm" gorm:"COMMENT:'提报类型代码'"`
	Node   string `json:"node" gorm:"COMMENT:'当前节点代码'"`
	Tsr    string `json:"tsr" gorm:"COMMENT:'推送人'"`
	Czr    string `json:"czr" gorm:"COMMENT:'操作人'"`
	Spyj   string `json:"spyj" gorm:"COMMENT:'审批意见'"`
	Czrq   string `json:"czrq" gorm:"COMMENT:'操作日期'"`
}

func AddProcMod(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProcMod(pm *Procmodify) error {
	if err := db.Table("procmodify").Where("id=?", pm.ID).Updates(&pm).Error; err != nil {
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
	if err := db.Raw("select procmodify.id,procmodify.procid,proctype.mc as dm,procmodify.node,user.name as tsr,procmodify.czr,procmodify.spyj,procmodify.czrq from procmodify left join proctype on procmodify.dm=proctype.dm left join user on user.mobile=procmodify.tsr where procmodify.procid=?", procid).
		Order("procmodify.id desc").Find(&pms).Error; err != nil {
		return nil, err
	}
	return pms, nil
}
