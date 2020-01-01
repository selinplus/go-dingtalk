package models

type Devmodetail struct {
	ID    uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Lsh   string `json:"lsh" gorm:"COMMENT:'流水号''"`
	Czlx  string `json:"czlx" gorm:"COMMENT:'操作类型'"`
	Czrq  string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Lx    string `json:"lx" gorm:"COMMENT:'设备类型'"`
	DevID string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"` //devinfo ID
	Zcbh  string `json:"zcbh" gorm:"COMMENT:'资产编号'"`
}

func AddDevModetail(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func GetDevModetails(devid string, pageNo, pageSize int) ([]*Devmodetail, error) {
	var devs []*Devmodetail
	offset := (pageNo - 1) * pageSize
	query := `	select devmodetail.id,devmodetail.lsh,devmodetail.devid,devoperation.mc as czlx,
				devtype.mc,devmodetail.devid,devmodetail.zcbh from devmodetail 
				left join devinfo on device.id=devmodetail.devid 
				left join devoperation on devoperation.dm=devmodetail.czlx 
				left join devtype on devtype.id=devmodetail.lx  
				where devmodetail.devid = ? order by devmodetail.id desc LIMIT ?,?`
	if err := db.Raw(query, devid, offset, pageSize).Scan(&devs).Error; err != nil {
		return nil, err
	}
	return devs, nil
}
