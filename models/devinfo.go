package models

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/qrcode"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

//new devinfo info
type Devinfo struct {
	Sbbh   uint   `json:"sbbh" gorm:"primary_key;AUTO_INCREMENT;COMMENT:'根据主键id生成6位设备编号'"`
	ID     string `json:"ID" gorm:"COMMENT:'设备编号'"`
	Zcbh   string `json:"zcbh" gorm:"COMMENT:'资产编号'"`
	Lx     string `json:"lx" gorm:"COMMENT:'设备类型'"`
	Mc     string `json:"mc" gorm:"COMMENT:'设备名称'"`
	Xh     string `json:"xh" gorm:"COMMENT:'设备型号'"`
	Xlh    string `json:"xlh" gorm:"COMMENT:'序列号'"`
	Ly     string `json:"ly" gorm:"COMMENT:'设备来源'"`
	Gys    string `json:"gys" gorm:"COMMENT:'供应商'"`
	Jg     string `json:"jg" gorm:"COMMENT:'价格'"`
	Scs    string `json:"scs" gorm:"COMMENT:'生产商'"`
	Scrq   string `json:"scrq" gorm:"COMMENT:'生产日期'"`
	Grrq   string `json:"grrq" gorm:"COMMENT:'购入日期'"`
	Bfnx   string `json:"bfnx" gorm:"COMMENT:'设备报废年限'"`
	QrUrl  string `json:"qrurl" gorm:"COMMENT:'二维码URL';column:qrurl"`
	Rkrq   string `json:"rkrq" gorm:"COMMENT:'入库日期'"`
	Czr    string `json:"czr" gorm:"COMMENT:'操作人'"`
	Czrq   string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Zt     string `json:"zt" gorm:"COMMENT:'设备状态'"`
	Jgdm   string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	Jgksdm string `json:"jgksdm" gorm:"COMMENT:'设备所属机构代码'"`
	Syr    string `json:"syr" gorm:"COMMENT:'设备使用人代码'"`
	Cfwz   string `json:"cfwz" gorm:"COMMENT:'存放位置'"`
	Sx     string `json:"sx" gorm:"COMMENT:'设备属性'"`
	Pnum   int    `json:"pnum" gorm:"default:0;COMMENT:'二维码打印次数'"`
	Img    string `json:"img" gorm:"COMMENT:'设备照片URL'"`
	Sbdl   int    `json:"sbdl" gorm:"default:1;COMMENT:'设备大类,1计算机类设备 2非计算类设备'"`
	Zw     int    `json:"zw" gorm:"default:1;COMMENT:'账外标志,1为账内 2为账外'"`
}

//生成设备编号
func GenerateSbbh(lx, xlh string) string {
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	sbbh := lx + xlh + timeStamp[:13]
	return sbbh
}

//判断序列号是否存在
func IsDevXlhExist(xlh string) bool {
	var dev Devinfo
	err := db.Table("devinfo").Where("xlh=?", xlh).First(&dev).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func IsUserDevExist(userid string) bool {
	var d Devinfo
	if err := db.Where("syr=?", userid).First(&d).Error; err != nil {
		return false
	}
	return true
}

//设备初始入库
func AddDevinfo(dev *Devinfo) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("devinfo").Create(dev).Error; err != nil {
		tx.Rollback()
		return err
	}
	lsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	t := time.Now().Format("2006-01-02 15:04:05")
	dm := &Devmod{
		Lsh:  lsh,
		Czrq: t,
		Czlx: "1",
		Num:  1,
		Czr:  dev.Czr,
		Jgdm: "00",
	}
	if err := tx.Table("devmod").Create(dm).Error; err != nil {
		tx.Rollback()
		return err
	}
	dmd := &Devmodetail{
		Lsh:   lsh,
		Czlx:  "1",
		Czrq:  t,
		Lx:    dev.Lx,
		DevID: dev.ID,
		Zcbh:  dev.Zcbh,
	}
	if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

type OpForm struct {
	Ids       []string `json:"ids"`
	Dms       []string `json:"dms"` //交回&批量收回
	SrcJgdm   string   `json:"src_jgdm"`
	DstJgdm   string   `json:"dst_jgdm"`   //分配,下发,上交
	SrcJgksdm string   `json:"src_jgksdm"` //调整前所属科室代码
	DstJgksdm string   `json:"dst_jgksdm"` //调整后所属科室代码
	SrcCfwz   string   `json:"src_cfwz"`   //调整前存放位置
	DstCfwz   string   `json:"dst_cfwz"`   //调整后存放位置
	Lsh       string   `json:"lsh"`        //上交时,用于修改devtodo表done
	Czr       string   `json:"czr"`        //inner传递操作人mobile
	Syr       string   `json:"syr"`        //inner传递使用人mobile
	CuserID   string   `json:"cuserid"`    //epp传递操作人userid
	SuserID   string   `json:"suserid"`    //epp传递使用人userid
	Cfwz      string   `json:"cfwz"`
	Czlx      string   `json:"czlx"`
	Agree     string   `json:"agree"` //Y:同意 N:不同意
}

//设备下发&上交
func DevIssued(form OpForm, czr, czlx string) error {
	var (
		ids     = form.Ids
		srcJgdm = form.SrcJgdm
		dstJgdm = form.DstJgdm
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	ckLsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	t := time.Now().Format("2006-01-02 15:04:05")
	dm := &Devmod{
		Lsh:  ckLsh,
		Czrq: t,
		Czlx: czlx,
		Num:  len(ids),
		Czr:  czr,
		Jgdm: srcJgdm,
	}
	if err := tx.Table("devmod").Create(dm).Error; err != nil {
		tx.Rollback()
		return err
	}
	zt, sx := getState(czlx)
	for _, id := range ids {
		dev := &Devinfo{
			ID:   id,
			Czrq: t,
			Czr:  czr,
			Zt:   zt,
			Sx:   sx,
		}
		if err := tx.Table("devinfo").Where("id=?", dev.ID).Updates(dev).Error; err != nil {
			tx.Rollback()
			return err
		}
		d, err := GetDevinfoByID(id)
		if err != nil {
			tx.Rollback()
			return err
		}
		dmd := &Devmodetail{
			Lsh:   ckLsh,
			Czlx:  czlx,
			Czrq:  t,
			Lx:    d.Lx,
			DevID: d.ID,
			Zcbh:  d.Zcbh,
		}
		if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if czlx == "10" { //设备上交,创建待办任务并发送消息给设备机构管理员
		dto := &Devtodo{
			Czlx: czlx,
			Lsh:  ckLsh,
			Czr:  czr,
			Czrq: t,
			Jgdm: dstJgdm,
			Bz:   "等待管理员做上交入库操作",
		}
		if err := tx.Table("devtodo").Create(dto).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else { //设备下发,同时创建入库
		rkLsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
		dm2 := &Devmod{
			Lsh:    rkLsh,
			PreLsh: ckLsh,
			Czrq:   t,
			Czlx:   "1",
			Num:    len(ids),
			Czr:    czr,
			Jgdm:   dstJgdm,
		}
		if err := tx.Table("devmod").Create(dm2).Error; err != nil {
			tx.Rollback()
			return err
		}
		zt, sx = getState("1")
		for _, id := range ids {
			dev := &Devinfo{
				ID:   id,
				Czrq: t,
				Czr:  czr,
				Jgdm: dstJgdm,
				Zt:   zt,
				Sx:   sx,
			}
			if err := tx.Table("devinfo").Where("id=?", dev.ID).Updates(dev).Error; err != nil {
				tx.Rollback()
				return err
			}
			d, err := GetDevinfoByID(id)
			if err != nil {
				tx.Rollback()
				return err
			}
			dmd := &Devmodetail{
				Lsh:   rkLsh,
				Czlx:  "1",
				Czrq:  t,
				Lx:    d.Lx,
				DevID: d.ID,
				Zcbh:  d.Zcbh,
			}
			if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

//设备机构变更申请,创建待办任务并发送消息给设备机构管理员
func ChangeJgks(form OpForm, czr string) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	//创建待办任务并发送消息给设备机构管理员
	for _, id := range form.Ids {
		dto := &Devtodo{
			Czlx:      form.Czlx,
			Czr:       czr,
			Czrq:      t,
			Jgdm:      form.SrcJgdm,
			DstJgdm:   GetSjjgdm(form.DstJgdm),
			SrcJgksdm: form.SrcJgksdm,
			DstJgksdm: form.DstJgksdm,
			SrcCfwz:   form.SrcCfwz,
			DstCfwz:   form.DstCfwz,
			DevID:     id,
			Bz:        "等待管理员设备机构变更申请",
		}
		if err := tx.Table("devtodo").Create(dto).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

//同意设备机构变更申请
func AgreeChangeJgks(form OpForm, czr string) error {
	var (
		id        = form.Ids[0]
		dstJgdm   = form.DstJgdm
		dstJgksdm = form.DstJgksdm
		dstCfwz   = form.DstCfwz
		czlx      = form.Czlx
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	var udpMap = map[string]interface{}{"done": 1}
	if form.Agree == "N" { //不同意直接返回
		udpMap["bz"] = "不同意"
		//修改待办标志为已办
		if err := tx.Table("devtodo").Where("done = 0 and devid = ? ", id).
			Updates(udpMap).Error; err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}
	udpMap["bz"] = "同意"
	//修改待办标志为已办
	if err := tx.Table("devtodo").Where("done = 0 and devid = ? ", id).
		Updates(udpMap).Error; err != nil {
		tx.Rollback()
		return err
	}
	//同意后，增加出库&入库流水
	t := time.Now().Format("2006-01-02 15:04:05")
	d, err := GetDevinfoByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	//增加使用人申请流水
	ckLsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	dm := &Devmod{
		Lsh:  ckLsh,
		Czrq: t,
		Czlx: czlx,
		Num:  1,
		Czr:  d.Syr,
		Jgdm: d.Jgdm,
	}
	if err := tx.Table("devmod").Create(dm).Error; err != nil {
		tx.Rollback()
		return err
	}
	//变更管理机构代码，所属机构代码，存放位置
	dev := map[string]string{
		"id":   id,
		"czrq": t,
		"czr":  czr,
	}
	if dstJgdm != "" {
		dev["jgdm"] = dstJgdm
	}
	if dstJgksdm != "" {
		dev["jgksdm"] = dstJgksdm
	}
	if dstCfwz != "" {
		dev["cfwz"] = dstCfwz
	}
	if err := tx.Table("devinfo").Where("id=?", dev["id"]).Updates(dev).Error; err != nil {
		tx.Rollback()
		return err
	}
	dmd := &Devmodetail{
		Lsh:   ckLsh,
		Czlx:  czlx,
		Czrq:  t,
		Lx:    d.Lx,
		DevID: d.ID,
		Zcbh:  d.Zcbh,
		Syr:   d.Syr,
		Bz:    "申请变更",
	}
	if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
		tx.Rollback()
		return err
	}
	//增加操作人同意申请流水
	lsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	dmGly := &Devmod{
		Lsh:    lsh,
		PreLsh: ckLsh,
		Czrq:   t,
		Czlx:   czlx,
		Num:    1,
		Czr:    czr,
		Jgdm:   dstJgdm,
	}
	if err := tx.Table("devmod").Create(dmGly).Error; err != nil {
		tx.Rollback()
		return err
	}
	dmdGly := &Devmodetail{
		Lsh:   lsh,
		Czlx:  czlx,
		Czrq:  t,
		Lx:    d.Lx,
		DevID: d.ID,
		Zcbh:  d.Zcbh,
		Syr:   d.Syr,
		Bz:    "同意变更",
	}
	if err := tx.Table("devmodetail").Create(dmdGly).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//设备交回申请,创建待办任务并发送消息给设备机构管理员
func DevReback(form OpForm, syr, czr string) error {
	var (
		ids  = form.Ids
		dms  = form.Dms
		jgdm = form.DstJgdm
		cfwz = form.Cfwz
		czlx = form.Czlx
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	for i, id := range ids {
		dev := map[string]string{
			"id":   id,
			"czrq": t,
			"czr":  czr,
			"syr":  syr,
			"cfwz": cfwz,
		}
		if jgdm != "" {
			dev["jgdm"] = jgdm
		} else {
			dev["jgdm"] = dms[i]
		}
		dto := &Devtodo{
			Czlx:    czlx,
			Czr:     czr,
			Czrq:    t,
			Lsh:     util.RandomString(4) + strconv.Itoa(int(time.Now().Unix())),
			Jgdm:    dev["jgdm"],
			DstJgdm: dev["jgdm"],
			DevID:   id,
			Bz:      "等待管理员审批交回申请",
		}
		if err := tx.Table("devtodo").Create(dto).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

//同意设备交回申请
func AgreeDevReback(form OpForm, czr string) error {
	var (
		id   = form.Ids[0]
		czlx = form.Czlx
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	var udpMap = map[string]interface{}{"done": 1}
	if form.Agree == "N" { //不同意直接返回
		udpMap["bz"] = "不同意"
		//修改待办标志为已办
		if err := tx.Table("devtodo").Where("done = 0 and devid = ? ", id).
			Updates(udpMap).Error; err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}
	udpMap["bz"] = "同意"
	//修改待办标志为已办
	if err := tx.Table("devtodo").Where("done = 0 and devid = ? ", id).
		Updates(udpMap).Error; err != nil {
		tx.Rollback()
		return err
	}
	//同意后，增加出库&入库流水
	ckLsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	t := time.Now().Format("2006-01-02 15:04:05")
	d, err := GetDevinfoByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	//增加交回出库流水
	dm := &Devmod{
		Lsh:  ckLsh,
		Czrq: t,
		Czlx: czlx,
		Num:  1,
		Czr:  d.Syr,
		Jgdm: d.Jgdm,
	}
	if err := tx.Table("devmod").Create(dm).Error; err != nil {
		tx.Rollback()
		return err
	}
	//置空所属机构代码，存放位置，使用人
	dev := map[string]string{
		"id":     id,
		"czrq":   t,
		"czr":    czr,
		"syr":    "",
		"cfwz":   "",
		"jgksdm": "",
	}
	if err := tx.Table("devinfo").Where("id=?", dev["id"]).Updates(dev).Error; err != nil {
		tx.Rollback()
		return err
	}
	dmd := &Devmodetail{
		Lsh:   ckLsh,
		Czlx:  czlx,
		Czrq:  t,
		Lx:    d.Lx,
		DevID: d.ID,
		Zcbh:  d.Zcbh,
	}
	if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
		tx.Rollback()
		return err
	}
	//增加操作人同意申请流水
	lsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	//增加使用人申请入库流水
	dmGly := &Devmod{
		Lsh:    lsh,
		PreLsh: ckLsh,
		Czrq:   t,
		Czlx:   "1",
		Num:    1,
		Czr:    czr,
		Jgdm:   d.Jgdm,
	}
	if err := tx.Table("devmod").Create(dmGly).Error; err != nil {
		tx.Rollback()
		return err
	}
	dmdGly := &Devmodetail{
		Lsh:   lsh,
		Czlx:  "1",
		Czrq:  t,
		Lx:    d.Lx,
		DevID: d.ID,
		Zcbh:  d.Zcbh,
		Syr:   d.Syr,
	}
	if err := tx.Table("devmodetail").Create(dmdGly).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//设备分配&借出&收回&上交入库
func DevAllocate(form OpForm, syr, czr string) error {
	var (
		ids     = form.Ids
		dms     = form.Dms
		jgdm    = form.DstJgdm
		jgksdm  = form.DstJgksdm
		cfwz    = form.Cfwz
		czlx    = form.Czlx
		todoLsh = form.Lsh
	)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	lsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	t := time.Now().Format("2006-01-02 15:04:05")
	dm := &Devmod{
		Lsh:  lsh,
		Czrq: t,
		Czlx: czlx,
		Num:  len(ids),
		Czr:  czr,
		Jgdm: jgdm,
	}
	if err := tx.Table("devmod").Create(dm).Error; err != nil {
		tx.Rollback()
		return err
	}
	if syr == " " {
		syr = ""
	}
	if cfwz == " " {
		cfwz = ""
	}
	if len(todoLsh) > 0 { // todoLsh:上交入库,根据todoLsh修改待办标志为已办
		if err := tx.Table("devtodo").Where("done = 0 and lsh = ? ", todoLsh).
			Updates(map[string]interface{}{"done": 1, "bz": "已完成上交入库"}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	zt, sx := getState(czlx)
	for i, id := range ids {
		dev := map[string]string{
			"id":   id,
			"czrq": t,
			"czr":  czr,
			"syr":  syr,
			"cfwz": cfwz,
			"zt":   zt,
			"sx":   sx,
		}
		if jgdm != "" {
			dev["jgdm"] = jgdm
		} else {
			dev["jgdm"] = dms[i]
		}
		if jgksdm != "" {
			dev["jgksdm"] = jgksdm
		}
		if err := tx.Table("devinfo").Where("id=?", dev["id"]).Updates(dev).Error; err != nil {
			tx.Rollback()
			return err
		}
		d, err := GetDevinfoByID(id)
		if err != nil {
			tx.Rollback()
			return err
		}
		dmd := &Devmodetail{
			Lsh:   lsh,
			Czlx:  czlx,
			Czrq:  t,
			Lx:    d.Lx,
			DevID: d.ID,
			Zcbh:  d.Zcbh,
			Syr:   syr,
		}
		if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
			tx.Rollback()
			return err
		}
		/*if czlx == "1" && len(todoLsh) == 0 { //交回设备入库,修改待办标志为已办
			if err := tx.Table("devtodo").Where("done = 0 and devid = ? ", id).
				Update("done", 1).Error; err != nil {
				tx.Rollback()
				return err
			}
		}*/
	}
	if czlx == "7" { //设备收回时,同步进行入库操作
		czlx = "1"
		lsh = util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
		t = time.Now().Format("2006-01-02 15:04:05")
		dm := &Devmod{
			Lsh:  lsh,
			Czrq: t,
			Czlx: czlx,
			Num:  len(ids),
			Czr:  czr,
			Jgdm: jgdm,
		}
		if err := tx.Table("devmod").Create(dm).Error; err != nil {
			tx.Rollback()
			return err
		}
		zt, sx := getState(czlx)
		if syr == " " {
			syr = ""
		}
		if cfwz == " " {
			cfwz = ""
		}
		for i, id := range ids {
			dev := map[string]string{
				"id":   id,
				"czrq": t,
				"czr":  czr,
				"syr":  syr,
				"cfwz": cfwz,
				"zt":   zt,
				"sx":   sx,
			}
			if jgdm != "" {
				dev["jgdm"] = jgdm
			} else {
				dev["jgdm"] = dms[i]
			}
			if jgksdm != "" {
				dev["jgksdm"] = jgksdm
			}
			if err := tx.Table("devinfo").Where("id=?", dev["id"]).Updates(dev).Error; err != nil {
				tx.Rollback()
				return err
			}
			d, err := GetDevinfoByID(id)
			if err != nil {
				tx.Rollback()
				return err
			}
			dmd := &Devmodetail{
				Lsh:   lsh,
				Czlx:  czlx,
				Czrq:  t,
				Lx:    d.Lx,
				DevID: d.ID,
				Zcbh:  d.Zcbh,
				Syr:   syr,
			}
			if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func getState(czlx string) (zt, sx string) {
	switch czlx {
	case "1":
		zt, sx = "1", "1"
	case "2", "10":
		zt, sx = "4", "1"
	case "3":
		zt, sx = "2", "3"
	case "4":
		zt, sx = "2", "4"
	case "5":
		zt, sx = "3", "3"
	case "6":
		zt, sx = "3", "4"
	case "7":
		zt, sx = "4", "2"
	case "8":
		zt, sx = "2", "4"
	case "9":
		zt, sx = "5", "5"
	}
	return zt, sx
}

func EditDevinfo(dev *Devinfo) error {
	if err := db.Table("devinfo").Where("id=?", dev.ID).Updates(dev).Error; err != nil {
		return err
	}
	return nil
}

func EditDevinfoCfwz(updMap map[string]interface{}) error {
	if err := db.Table("devinfo").
		Where("id=?", updMap["id"]).Update("cfwz", updMap["cfwz"]).Error; err != nil {
		return err
	}
	return nil
}

func EditDevinfoPnum(id uint) error {
	s := fmt.Sprintf("update devinfo set pnum=pnum+1 where sbbh = %d", id)
	if err := db.Exec(s).Error; err != nil {
		return err
	}
	return nil
}

func DelDevinfo(id string) error {
	if err := db.Where("id=?", id).Delete(Devinfo{}).Error; err != nil {
		return err
	}
	return nil
}

//批量导入
type DevinfoErr struct {
	*Devinfo
	Msg string `json:"msg"`
}

func ImpDevinfos(fileName io.Reader, czr string) ([]*DevinfoErr, int, int, error) {
	devs, err := ReadDevinfoXmlToStructs(fileName, czr)
	if err != nil {
		log.Println(err)
		return nil, 0, 0, err
	}
	errDev, success, failed := InsertDevinfoXml(devs, czr)
	return errDev, success, failed, nil
}

func ReadDevinfoXmlToStructs(fileName io.Reader, czr string) ([]*Devinfo, error) {
	devs := make([]*Devinfo, 0)
	xlsx, err := excelize.OpenReader(fileName)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	//sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows("设备基本信息表")
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	//logging.Info(fmt.Sprintf("sheet name: %s", sheetName))
	//遍历行读取
	for k, row := range rows {
		// 跳过标题行，遍历每行的列读取
		if k == 0 {
			continue
		}
		d := Devinfo{}
		d.Czr = czr
		for i, cell := range row {
			if cell == "" {
				switch {
				case i == 1, i == 2, i == 3, i == 4, i == 5, i == 7, i == 8, i == 9:
					return nil, fmt.Errorf("%s", "文件校验错误，存在未录入项！")
				}
			}
			switch {
			case i == 0:
				d.Zcbh = cell
			case i == 1:
				d.Xlh = cell
			case i == 2:
				d.Lx = cell
			case i == 3:
				d.Mc = cell
			case i == 4:
				d.Grrq = cell
			case i == 5:
				d.Jg = cell
			case i == 6:
				d.Ly = cell
			case i == 7:
				d.Scrq = cell
			case i == 8:
				d.Scs = cell
			case i == 9:
				d.Xh = cell
			case i == 10:
				d.Gys = cell
			}
		}
		//logging.Debug(fmt.Sprintf("*: %+v", d))
		devs = append(devs, &d)
	}
	return devs, nil
}

func InsertDevinfoXml(devs []*Devinfo, czr string) ([]*DevinfoErr, int, int) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var (
		errDev   = make([]*DevinfoErr, 0)
		devsChan = make(chan *Devinfo)
		cntChan  = make(chan int)
		wg       = &sync.WaitGroup{}
		devsNum  = len(devs)
		wtNum    = 50
		seg      int
		cnt      int
	)
	//logging.Debug(fmt.Sprintf("------------------%d------", len(devs)))
	if devsNum > 0 {
		if devsNum%wtNum == 0 {
			seg = devsNum / wtNum
		} else {
			seg = (devsNum / wtNum) + 1
		}
		for j := 0; j < wtNum; j++ {
			beg := j * seg
			if beg > devsNum-1 {
				break
			}
			var end int
			if (j+1)*seg < devsNum {
				end = (j + 1) * seg
			} else {
				end = devsNum
			}
			//log.Println(beg, end)
			segDevs := devs[beg:end]
			go func() {
				for i, segDev := range segDevs {
					if segDev != nil {
						devsChan <- segDev
						cntChan <- i
					}
				}
			}()
		}
		go func() {
			for range cntChan {
				cnt++
				if cnt == devsNum {
					close(devsChan)
				}
			}
		}()
		lsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
		for k := 0; k < wtNum; k++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for dev := range devsChan {
					//if dev == nil {
					//	break
					//}
					t := time.Now().Format("2006-01-02 15:04:05")
					d := &Devinfo{
						Zcbh: dev.Zcbh,
						Mc:   dev.Mc,
						Xh:   dev.Xh,
						Xlh:  dev.Xlh,
						Ly:   dev.Ly,
						Scs:  dev.Scs,
						Scrq: dev.Scrq,
						Grrq: dev.Grrq,
						Bfnx: dev.Bfnx,
						Jg:   dev.Jg,
						Gys:  dev.Gys,
						Rkrq: t,
						Czrq: t,
						Czr:  dev.Czr,
						Zt:   "1",
						Jgdm: "00",
						Sx:   "1",
					}
					LxDm, err := GetDevtypeByMc(dev.Lx)
					if err != nil {
						errDev = append(errDev,
							&DevinfoErr{
								Devinfo: dev,
								Msg:     "获取设备类型代码失败,设备类型名称错误！",
							})
					} else {
						if IsDevXlhExist(dev.Xlh) {
							//logging.Info(fmt.Sprintf("%s:序列号已存在!", dev.Xlh))
							errDev = append(errDev,
								&DevinfoErr{
									Devinfo: dev,
									Msg:     "序列号已存在！",
								})
						} else {
							d.Lx = LxDm.Dm
							d.ID = GenerateSbbh(d.Lx, d.Xlh)
							//生成二维码
							info := d.ID + "$序列号[" + d.Xlh + "]$生产商[" + d.Scs + "]$设备型号[" + d.Xh + "]$生产日期[" + d.Scrq + "]$"
							name, _, err := qrcode.GenerateQrWithLogo(info, qrcode.GetQrCodeFullPath())
							if err != nil {
								log.Println(err)
							}
							d.QrUrl = qrcode.GetQrCodeFullUrl(name)
							if err := tx.Table("devinfo").Create(d).Error; err != nil {
								tx.Rollback()
								return
							}
							dmd := &Devmodetail{
								Lsh:   lsh,
								Czlx:  "1",
								Lx:    d.Lx,
								DevID: d.ID,
								Zcbh:  d.Zcbh,
								Czrq:  time.Now().Format("2006-01-02 15:04:05"),
							}
							if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
								tx.Rollback()
								return
							}
						}
					}

				}
			}()
		}
		wg.Wait()
		dm := &Devmod{
			Lsh:  lsh,
			Czrq: time.Now().Format("2006-01-02 15:04:05"),
			Czlx: "1",
			Num:  devsNum - len(errDev),
			Czr:  czr,
			Jgdm: "00",
		}
		if err := tx.Table("devmod").Create(dm).Error; err != nil {
			tx.Rollback()
			return nil, 0, 0
		}
		if devsNum == len(errDev) {
			tx.Rollback()
			return errDev, devsNum - len(errDev), len(errDev)
		}
		tx.Commit()
	}
	if len(errDev) > 0 {
		return errDev, devsNum - len(errDev), len(errDev)
	}
	return nil, devsNum, 0
}

type DevinfoResp struct {
	*Devinfo
	Jgmc      string `json:"jgmc"`
	Ksmc      string `json:"ksmc"`
	SyrName   string `json:"syr_name"`
	SyrMobile string `json:"syr_mobile"`
	Idstr     string `json:"idstr"` //6位短编号
}

func GetDevinfos(con map[string]string, pageNo, pageSize int, bz string) ([]*DevinfoResp, error) {
	query := `select devinfo.sbbh,devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,devinfo.rkrq,devinfo.pnum,
			devinfo.czrq,c.name as czr,devinfo.qrurl,devstate.mc as zt,a.jgdm as jgdm,a.jgmc as jgmc,
			b.jgdm as jgksdm,b.jgmc as ksmc,devinfo.cfwz,devproperty.mc as sx,devinfo.syr,
			(case when (d.name ='' OR d.name is null) then devinfo.syr else d.name end) as syr_name,d.mobile as syr_mobile,
       		concat(repeat('0',6-length(devinfo.sbbh)),devinfo.sbbh) as idstr
			from devinfo 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devproperty on devproperty.dm=devinfo.sx 
			left join devdept a on a.jgdm=devinfo.jgdm 
			left join devdept b on b.jgdm=devinfo.jgksdm 
			left join userdemo c on c.userid=devinfo.czr 
			left join userdemo d on d.userid=devinfo.syr 
			where 1=1`
	if con["jgdm"] != "" {
		query += fmt.Sprintf(" and devinfo.jgdm = '%s'", con["jgdm"])
	}
	if len(bz) > 0 {
		if bz == "0" {
			query += " and devinfo.zt = '1'"
		} else if bz == "3" {
			query += " and devinfo.zt = '2' and devinfo.sx = '3'"
		} else if bz == "4" {
			query += " and devinfo.zt = '2' and devinfo.sx = '4'"
		} else if bz == "6" {
			query += " and devinfo.zt = '3' and devinfo.sx = '4'"
		} else if bz == "10" {
			query += " and ((devinfo.zt = '2' and devinfo.sx = '3') or(devinfo.zt = '2' and devinfo.sx = '4')or(devinfo.zt = '3' and devinfo.sx = '4'))"
		}
	}
	if con["mc"] != "" {
		query += fmt.Sprintf(" and devinfo.mc like '%s'", "%"+con["mc"]+"%")
	}
	if con["xlh"] != "" {
		query += fmt.Sprintf(" and devinfo.xlh like '%s'", "%"+con["xlh"]+"%")
	}
	if con["zcbh"] != "" {
		query += fmt.Sprintf(" and devinfo.zcbh like '%s'", "%"+con["zcbh"]+"%")
	}
	if con["syr"] != "" {
		query += fmt.Sprintf(" and devinfo.syr = '%s'", con["syr"])
	}
	if con["sbbh"] != "" {
		query += fmt.Sprintf(" and devinfo.sbbh = '%s'", con["sbbh"])
	}
	if con["rkrqq"] != "" && con["rkrqz"] != "" {
		query += fmt.Sprintf(" and devinfo.rkrq >= '%s' and devinfo.rkrq <= '%s'", con["rkrqq"], con["rkrqz"])
	}
	if con["scrqq"] != "" && con["scrqz"] != "" {
		query += fmt.Sprintf(" and devinfo.scrq >= '%s' and devinfo.scrq <= '%s'", con["scrqq"], con["scrqz"])
	}
	query += ` order by devinfo.sbbh`
	if pageNo > 0 && pageSize > 0 {
		offset := (pageNo - 1) * pageSize
		query += fmt.Sprintf(" LIMIT %d,%d", offset, pageSize)
	}
	var devs []*DevinfoResp
	if err := db.Raw(query).Scan(&devs).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return devs, nil
}

func GetDevinfosGly(con map[string]string) ([]*DevinfoResp, error) {
	var devs []*DevinfoResp
	squery := `select devinfo.sbbh,devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,devinfo.rkrq,devinfo.pnum,
			devinfo.czrq,c.name as czr,devinfo.qrurl,devstate.mc as zt,a.jgdm as jgdm,a.jgmc as jgmc,
			b.jgdm as jgksdm,b.jgmc as ksmc,devinfo.cfwz,devproperty.mc as sx,devinfo.syr,
			(case when (d.name ='' OR d.name is null) then devinfo.syr else d.name end) as syr_name,d.mobile as syr_mobile,
       		concat(repeat('0',6-length(devinfo.sbbh)),devinfo.sbbh) as idstr
			from devinfo 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devproperty on devproperty.dm=devinfo.sx 
			left join devdept a on a.jgdm=devinfo.jgdm 
			left join devdept b on b.jgdm=devinfo.jgksdm 
			left join userdemo c on c.userid=devinfo.czr 
			left join userdemo d on d.userid=devinfo.syr 
			where devinfo.jgdm = '` + con["jgdm"] + `' `
	if con["sbbh"] != "" {
		squery += ` and devinfo.sbbh = '` + con["sbbh"] + `' `
	}
	if con["property"] != "" {
		squery += ` and devinfo.sx = '` + con["property"] + `' `
	}
	if con["state"] != "" {
		squery += ` and devinfo.zt = '` + con["state"] + `' `
	}
	if con["type"] != "" {
		squery += ` and devinfo.lx like '` + con["type"] + `%' `
	}
	if con["zcbh"] != "" {
		squery += ` and devinfo.zcbh like '` + con["zcbh"] + `%' `
	}
	if con["scrq"] != "" {
		squery += ` and devinfo.scrq like '` + con["scrq"] + `%' `
	}
	if con["rkrq"] != "" {
		squery += ` and devinfo.rkrq like '` + con["rkrq"] + `%' `
	}
	if con["xlh"] != "" {
		squery += ` and devinfo.xlh = '` + con["xlh"] + `' `
	}
	if con["rkrqq"] != "" && con["rkrqz"] != "" {
		squery += fmt.Sprintf(" and devinfo.rkrq >= '%s' and devinfo.rkrq <= '%s'", con["rkrqq"], con["rkrqz"])
	}
	if con["scrqq"] != "" && con["scrqz"] != "" {
		squery += fmt.Sprintf(" and devinfo.scrq >= '%s' and devinfo.scrq <= '%s'", con["scrqq"], con["scrqz"])
	}
	squery += ` order by devinfo.sbbh`
	//log.Println(squery)
	if err := db.Raw(squery).Scan(&devs).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return devs, nil
}

func GetDevinfoRespByID(id string) (*DevinfoResp, error) {
	var dev DevinfoResp
	query := `select devinfo.sbbh,devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,devinfo.rkrq,devinfo.pnum,
			devinfo.czrq,c.name as czr,devinfo.qrurl,devstate.mc as zt,a.jgdm as jgdm,a.jgmc as jgmc,
			b.jgdm as jgksdm ,b.jgmc as ksmc,devinfo.cfwz,devproperty.mc as sx,devinfo.syr,
			(case when (d.name ='' OR d.name is null) then devinfo.syr else d.name end) as syr_name,d.mobile as syr_mobile,
       		concat(repeat('0',6-length(devinfo.sbbh)),devinfo.sbbh) as idstr
			from devinfo 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devproperty on devproperty.dm=devinfo.sx 
			left join devdept a on a.jgdm=devinfo.jgdm 
			left join devdept b on b.jgdm=devinfo.jgksdm 
			left join userdemo c on c.userid=devinfo.czr 
			left join userdemo d on d.userid=devinfo.syr 
			where devinfo.id = '%s'`
	squery := fmt.Sprintf(query, id)
	if err := db.Raw(squery).Scan(&dev).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
	}
	return nil, nil
}

func GetDevinfoByID(id string) (*Devinfo, error) {
	var dev Devinfo
	if err := db.Table("devinfo").Where("id=?", id).First(&dev).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
	}
	return nil, nil
}

func GetDevinfosToBeStored() ([]*Devinfo, error) {
	var devs []*Devinfo
	if err := db.Table("devinfo").
		Where("zt=4").Scan(&devs).Error; err != nil {
		return nil, err
	}
	return devs, nil
}

func GetDevinfoBySbbh(sbbh uint) *Devinfo {
	var dev Devinfo
	err := db.Table("devinfo").Where("sbbh=?", sbbh).First(&dev).Error
	if err != nil {
		return nil
	}
	if len(dev.ID) > 0 {
		return &dev
	}
	return nil
}

func GetDevinfosByJgdm(jgdm string) ([]*Devinfo, error) {
	var devs []*Devinfo
	err := db.Table("devinfo").Where("jgdm=?", jgdm).Find(&devs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if len(devs) > 0 {
		return devs, nil
	}
	return nil, nil
}

func ConvSbbhToIdstr(sbbh uint) (idstr string) {
	switch {
	case sbbh < 10:
		idstr = "00000" + strconv.Itoa(int(sbbh))
	case sbbh >= 10 && sbbh < 100:
		idstr = "0000" + strconv.Itoa(int(sbbh))
	case sbbh >= 100 && sbbh < 1000:
		idstr = "000" + strconv.Itoa(int(sbbh))
	case sbbh >= 1000 && sbbh < 10000:
		idstr = "00" + strconv.Itoa(int(sbbh))
	case sbbh >= 10000 && sbbh < 100000:
		idstr = "0" + strconv.Itoa(int(sbbh))
	case sbbh >= 100000:
		idstr = strconv.Itoa(int(sbbh))
	}
	return idstr
}
