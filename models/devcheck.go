package models

type Devcheck struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT;COMMENT:'盘点任务编码'"`
	Beg    string `json:"beg" gorm:"COMMENT:'时间起'"`
	End    string `json:"end" gorm:"COMMENT:'时间止'"`
	Fqr    string `json:"fqr" gorm:"COMMENT:'发起人代码'"`
	Ckself string `json:"ckself" gorm:"COMMENT:'是否自我盘点（Y,N）'"`
}

func AddDevCheckTask(ckTask *Devcheck) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("devcheck").Create(ckTask).Error; err != nil {
		tx.Rollback()
		return err
	}
	var devs []*Devinfo
	if err := tx.Table("devinfo").Find(&devs).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, dev := range devs {
		var syrJgdm string
		u, err := GetUserByUserid(dev.Syr)
		if err != nil {
			syrJgdm = ""
		} else {
			syrJgdm = u.Department
		}
		var ck = &Devckdetail{
			CheckID:   ckTask.ID,
			DevinfoID: dev.ID,
			Zcbh:      dev.Zcbh,
			Lx:        dev.Lx,
			Mc:        dev.Mc,
			Xh:        dev.Xh,
			Xlh:       dev.Xlh,
			Ly:        dev.Ly,
			Gys:       dev.Gys,
			Jg:        dev.Jg,
			Scs:       dev.Scs,
			Scrq:      dev.Scrq,
			Grrq:      dev.Grrq,
			Bfnx:      dev.Bfnx,
			Rkrq:      dev.Rkrq,
			Czr:       dev.Czr,
			Czrq:      dev.Czrq,
			Zt:        dev.Zt,
			Jgdm:      dev.Jgdm,
			Syr:       dev.Syr,
			SyrJgdm:   syrJgdm,
			Cfwz:      dev.Cfwz,
			Sx:        dev.Sx,
		}
		err = tx.Table("devckdetail").Create(ck).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Error
}
