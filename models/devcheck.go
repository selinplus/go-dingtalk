package models

import "log"

type Devcheck struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT;COMMENT:'盘点任务编码'"`
	Beg    string `json:"beg" gorm:"COMMENT:'时间起'"`
	End    string `json:"end" gorm:"COMMENT:'时间止'"`
	Fqr    string `json:"fqr" gorm:"COMMENT:'发起人代码'"`
	Sbdl   int    `json:"sbdl" gorm:"default:1;COMMENT:'设备大类,1计算机类设备 2非计算类设备'"`
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
	if err := tx.Table("devinfo").
		Where("sbdl=?", ckTask.Sbdl).Find(&devs).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, dev := range devs {
		var ck = &Devckdetail{
			CheckID:   ckTask.ID,
			Sbbh:      ConvSbbhToIdstr(dev.Sbbh),
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
			SyrJgdm:   dev.Jgksdm,
			Cfwz:      dev.Cfwz,
			Sx:        dev.Sx,
		}
		if dev.Syr == "" {
			switch dev.Sx {
			case "1": //在库设备
				devdept, err := GetDevdept(dev.Jgdm)
				if err != nil {
					log.Printf("获取在库设备管理机构[代码%s]失败:%v", dev.Jgdm, err)
					tx.Rollback()
					return err
				}
				ck.Syr = devdept.Gly
				if dev.Jgdm == "00" && dev.Sbdl == 2 {
					ck.Syr = devdept.Gly2
				}
			case "3": //共用设备
				devdept, err := GetDevdept(dev.Jgksdm)
				if err != nil {
					log.Printf("获取共用设备所属机构[代码%s]失败:%v", dev.Jgdm, err)
					tx.Rollback()
					return err
				}
				ck.Syr = devdept.Bgr
			}
		}
		if err := tx.Table("devckdetail").Create(ck).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func GetDevCheckTask(cond string, pageNo, pageSize int) ([]*Devcheck, error) {
	var devchecks []*Devcheck
	if err := db.Table("devcheck").Where(cond).
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).
		Find(&devchecks).Error; err != nil {
		return nil, err
	}
	return devchecks, nil
}

func GetDevCheckTasksCnt(cond string) (cnt int) {
	err := db.Table("devcheck").Where(cond).Count(&cnt).Error
	if err != nil {
		cnt = 0
	}
	return cnt
}
