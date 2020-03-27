package models

type Devmod struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Lsh    string `json:"lsh" gorm:"COMMENT:'流水号'"`
	PreLsh string `json:"pre_lsh" gorm:"COMMENT:'前置流水号'"`
	Czlx   string `json:"czlx" gorm:"COMMENT:'操作类型'"`
	Czrq   string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Num    int    `json:"num" gorm:"COMMENT:'设备数'"`
	Czr    string `json:"czr" gorm:"COMMENT:'操作人'"`
	Jgdm   string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
}

func GetDevMods(czr string) ([]*Devmod, error) {
	var devs []*Devmod
	query := `select devmod.id,devmod.lsh,devmod.pre_lsh,devoperation.mc as czlx,devmod.czrq,
			devmod.num,user.name as czr,devdept.jgmc as jgdm from devmod 
			left join devoperation on devoperation.dm=devmod.czlx 
			left join devdept on devdept.jgdm=devmod.jgdm  
			left join user on user.userid=devmod.czr  
			where devmod.czr=? order by devmod.czrq desc`
	if err := db.Raw(query, czr).Scan(&devs).Error; err != nil {
		return nil, err
	}
	return devs, nil
}
