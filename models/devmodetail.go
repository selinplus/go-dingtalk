package models

type Devmodetail struct {
	ID    uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Lsh   string `json:"lsh" gorm:"COMMENT:'流水号'"`
	Czlx  string `json:"czlx" gorm:"COMMENT:'操作类型'"`
	Czrq  string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Lx    string `json:"lx" gorm:"COMMENT:'设备类型'"`
	DevID string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"` //devinfo ID
	Zcbh  string `json:"zcbh" gorm:"COMMENT:'资产编号'"`
	Syr   string `json:"syr" gorm:"COMMENT:'设备使用人代码'"`
	Bz    string `json:"bz" gorm:"COMMENT:'备注'"`
}

type DevmodResp struct {
	*Devmodetail
	Idstr string `json:"idstr"` //6位短编号
	Sbbh  uint   `json:"sbbh"`  //主键
	Mc    string `json:"mc"`
	Zt    string `json:"zt"`
	Rkrq  string `json:"rkrq"`
	Ly    string `json:"ly"`
	Czr   string `json:"czr"`
	Jgdm  string `json:"jgdm"`
	Jgmc  string `json:"jgmc"`
}

func GetDevModetails(lsh string, pageNo, pageSize int) ([]*DevmodResp, error) {
	var devs []*DevmodResp
	offset := (pageNo - 1) * pageSize
	query := `	select devmodetail.id,devmodetail.lsh,devmodetail.devid,devoperation.mc as czlx,devmodetail.syr,
				devinfo.mc,devstate.mc as zt,devmodetail.czrq,devtype.mc as lx,devmodetail.devid,devmodetail.zcbh, 
				userdemo.name as czr,devmodetail.bz,devinfo.sbbh
				from devmodetail 
				left join devinfo on devinfo.id=devmodetail.devid 
				left join devoperation on devoperation.dm=devmodetail.czlx 
				left join devtype on devtype.dm=devmodetail.lx  
				left join devstate on devstate.dm=devinfo.zt 
				left join devmod on devmod.lsh=devmodetail.lsh 
				left join userdemo on userdemo.userid=devmod.czr 
				where devmodetail.lsh = ? order by devmodetail.id LIMIT ?,?`
	if err := db.Raw(query, lsh, offset, pageSize).Scan(&devs).Error; err != nil {
		return nil, err
	}
	return devs, nil
}

func GetDevModsByDevid(devid string) ([]*DevmodResp, error) {
	var devs []*DevmodResp
	query := `select devmodetail.id,devmodetail.lsh,devmodetail.devid,devmodetail.czrq,devmodetail.zcbh,
			devinfo.mc,devinfo.rkrq,devinfo.ly,devmod.jgdm,devdept.jgmc,userdemo.name as czr,devmodetail.syr,			
			devoperation.mc as czlx,devtype.mc as lx,devmodetail.bz
			from devmodetail 
			left join devinfo on devinfo.id=devmodetail.devid 
			left join devmod on devmod.lsh=devmodetail.lsh 
			left join devoperation on devoperation.dm=devmodetail.czlx 
			left join devtype on devtype.dm=devmodetail.lx 
			left join devdept on devdept.jgdm=devmod.jgdm 
			left join userdemo on userdemo.userid=devmod.czr 
			where devmodetail.devid = ? order by devmodetail.id desc`
	if err := db.Raw(query, devid).Scan(&devs).Error; err != nil {
		return nil, err
	}
	return devs, nil
}
