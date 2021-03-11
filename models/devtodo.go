package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Devtodo struct {
	ID         uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Czlx       string `json:"czlx" gorm:"COMMENT:'操作类型'"`
	Lsh        string `json:"lsh" gorm:"COMMENT:'流水号'"`
	Czr        string `json:"czr" gorm:"COMMENT:'操作人'"`
	Czrq       string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Jgdm       string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	DstJgdm    string `json:"dst_jgdm" gorm:"COMMENT:'更改后设备管理机构代码'"`
	SrcJgksdm  string `json:"src_jgksdm" gorm:"COMMENT:'设备所属机构代码'"`
	DstJgksdm  string `json:"dst_jgksdm" gorm:"COMMENT:'更改后设备所属机构代码'"`
	SrcCfwz    string `json:"src_cfwz" gorm:"COMMENT:'存放位置'"`
	DstCfwz    string `json:"dst_cfwz" gorm:"COMMENT:'更改后存放位置'"`
	Bz         string `json:"bz" gorm:"COMMENT:'待办备注'"`
	DevID      string `json:"devid" gorm:"COMMENT:'设备编号';column:devid"` //devinfo ID
	Done       int    `json:"done" gorm:"COMMENT:'0: 待办 1: 已办';size:1;default:'0'"`
	FlagNotice int    `json:"flag_notice" gorm:"COMMENT:'0: 未推送 1: 已推送';size:1;default:'0'"`
}

type DevtodoResp struct {
	Devtodo
	Czrid   string `json:"czrid"`
	Gly     string `json:"gly"`
	Bgr     string `json:"bgr"`
	Zcbh    string `json:"zcbh"`
	Mc      string `json:"mc"`
	Zt      string `json:"zt"`
	Num     int    `json:"num"`
	SrcJgdm string `json:"src_jgdm"`
	Jgmc    string `json:"jgmc"`     //调整前管理机构名称
	DstJgmc string `json:"dst_jgmc"` //调整后管理机构名称
	SrcKsmc string `json:"src_ksmc"` //调整前所属科室名称
	DstKsmc string `json:"dst_ksmc"` //调整后所属科室名称
}

func GetDevTodosOrDones(done int) ([]DevtodoResp, error) {
	var dtos []DevtodoResp
	sql := fmt.Sprintf(`
select devtodo.id,devtodo.czlx,devtodo.lsh,devtodo.czr as czrid,userdemo.name as czr,devtodo.czrq,
	devtodo.src_cfwz,devtodo.dst_cfwz,devtodo.jgdm,a.jgmc,devtodo.dst_jgdm,b.jgmc as dst_jgmc,
	devtodo.src_jgksdm,c.jgmc as src_ksmc,devtodo.dst_jgksdm,d.jgmc as dst_ksmc,
	case when devtodo.dst_jgdm='00' and devinfo.sbdl=2 then b.gly2 else b.gly end gly,
	devtodo.devid,devinfo.zcbh,devinfo.mc,devinfo.zt,devtodo.done,devtodo.bz
	from devtodo
	left join devdept a on a.jgdm=devtodo.jgdm
	left join devdept b on b.jgdm=devtodo.dst_jgdm
	left join devdept c on c.jgdm=devtodo.src_jgksdm
	left join devdept d on d.jgdm=devtodo.dst_jgksdm
	left join devinfo on devinfo.id=devtodo.devid
	left join userdemo on userdemo.userid=devtodo.czr 
	where devtodo.done=%d
	order by devtodo.czrq desc`, done)
	err := db.Raw(sql).Scan(&dtos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return dtos, nil
}

func GetDevTodosOrDonesByToid(id uint, done int) (*DevtodoResp, error) {
	var devtodoResp DevtodoResp
	sql := fmt.Sprintf(`
select devtodo.id,devtodo.czlx,devtodo.lsh,devtodo.czr as czrid,userdemo.name as czr,devtodo.czrq,
	devtodo.src_cfwz,devtodo.dst_cfwz,devtodo.jgdm,a.jgmc,devtodo.dst_jgdm,b.jgmc as dst_jgmc,
	devtodo.src_jgksdm,c.jgmc as src_ksmc,devtodo.dst_jgksdm,d.jgmc as dst_ksmc,
	case when devtodo.dst_jgdm='00' and devinfo.sbdl=2 then b.gly2 else b.gly end gly,
	devtodo.devid,devinfo.zcbh,devinfo.mc,devinfo.zt,devtodo.done,devtodo.bz
	from devtodo
	left join devdept a on a.jgdm=devtodo.jgdm
	left join devdept b on b.jgdm=devtodo.dst_jgdm
	left join devdept c on c.jgdm=devtodo.src_jgksdm
	left join devdept d on d.jgdm=devtodo.dst_jgksdm
	left join devinfo on devinfo.id=devtodo.devid
	left join userdemo on userdemo.userid=devtodo.czr 
	where devtodo.done=%d and devtodo.id=%d
	order by devtodo.czrq desc`, done, id)
	err := db.Raw(sql).Scan(&devtodoResp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &devtodoResp, nil
}

func GetUpDevTodosOrDones(done int, userid string) ([]DevtodoResp, error) {
	var dtos []DevtodoResp
	sql := fmt.Sprintf(`
select * from (
    select devtodo.id,devtodo.czlx,devtodo.lsh,userdemo.name as czr,
           devtodo.czrq,devtodo.jgdm,a.jgmc,devtodo.dst_jgdm,
           case when devtodo.dst_jgdm='00' and devinfo.sbdl=2 then b.gly2
                else b.gly end gly,
           devtodo.done,devmod.num,devinfo.sbdl
        from
           devtodo
               left join devinfo on devinfo.id = devtodo.devid
               left join userdemo on userdemo.userid = devtodo.czr
               left join devmod on devmod.lsh = devtodo.lsh
               left join devdept a on a.jgdm = devtodo.jgdm
               left join devdept b on b.jgdm = devtodo.dst_jgdm
        where devtodo.done = %d
) t where gly = '%s' order by czrq desc`, done, userid)
	if err := db.Raw(sql).Scan(&dtos).Error; err != nil {
		return nil, err
	}
	return dtos, nil
}

func GetDevFlag() ([]*DevtodoResp, error) {
	var dtos []*DevtodoResp
	sql := fmt.Sprintf(`
    select devtodo.id,devtodo.czlx,devtodo.lsh,devtodo.czr,devtodo.czrq,
           devtodo.src_cfwz,devtodo.dst_cfwz,devtodo.jgdm,a.jgmc,devtodo.dst_jgdm,b.jgmc as dst_jgmc,
           devtodo.src_jgksdm,c.jgmc as src_ksmc,devtodo.dst_jgksdm,d.jgmc as dst_ksmc,devmod.num,
           case when devtodo.dst_jgdm='00' and devinfo.sbdl=2 then b.gly2
                else b.gly end gly,
           devtodo.devid,devinfo.zcbh,devinfo.mc,devinfo.zt,devtodo.done,devtodo.bz
        from devtodo
               left join devdept a on a.jgdm=devtodo.jgdm
               left join devdept b on b.jgdm=devtodo.dst_jgdm
               left join devdept c on c.jgdm=devtodo.src_jgksdm
               left join devdept d on d.jgdm=devtodo.dst_jgksdm
               left join devinfo on devinfo.id=devtodo.devid
               left join devmod on devmod.lsh = devtodo.lsh
        where devtodo.flag_notice=0
        order by devtodo.czrq desc`)
	if err := db.Raw(sql).Scan(&dtos).Error; err != nil {
		return nil, err
	}
	return dtos, nil
}

func UpdateDevtodoFlag(id uint) error {
	if err := db.Table("devtodo").Where("id = ? and flag_notice = 0", id).
		Update("flag_notice", 1).Error; err != nil {
		return err
	}
	return nil
}
