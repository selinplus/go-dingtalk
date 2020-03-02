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
	Sbbh  uint   `json:"sbbh" gorm:"primary_key;AUTO_INCREMENT"`
	ID    string `json:"ID" gorm:"COMMENT:'设备编号'"`
	Zcbh  string `json:"zcbh" gorm:"COMMENT:'资产编号'"`
	Lx    string `json:"lx" gorm:"COMMENT:'设备类型'"`
	Mc    string `json:"mc" gorm:"COMMENT:'设备名称'"`
	Xh    string `json:"xh" gorm:"COMMENT:'设备型号'"`
	Xlh   string `json:"xlh" gorm:"COMMENT:'序列号'"`
	Ly    string `json:"ly" gorm:"COMMENT:'设备来源'"`
	Gys   string `json:"gys" gorm:"COMMENT:'供应商'"`
	Jg    string `json:"jg" gorm:"COMMENT:'价格'"`
	Scs   string `json:"scs" gorm:"COMMENT:'生产商'"`
	Scrq  string `json:"scrq" gorm:"COMMENT:'生产日期'"`
	Grrq  string `json:"grrq" gorm:"COMMENT:'购入日期'"`
	Bfnx  string `json:"bfnx" gorm:"COMMENT:'设备报废年限'"`
	QrUrl string `json:"qrurl" gorm:"COMMENT:'二维码URL';column:qrurl"`
	Rkrq  string `json:"rkrq" gorm:"COMMENT:'入库日期'"`
	Czr   string `json:"czr" gorm:"COMMENT:'操作人'"`
	Czrq  string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Zt    string `json:"zt" gorm:"COMMENT:'设备状态'"`
	Jgdm  string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	Syr   string `json:"syr" gorm:"COMMENT:'设备使用人代码'"`
	Cfwz  string `json:"cfwz" gorm:"COMMENT:'存放位置'"`
	Sx    string `json:"sx" gorm:"COMMENT:'设备属性'"`
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

//设备入库
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

//设备下发
func DevIssued(ids []string, srcJgdm, dstJgdm, czr string) error {
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
		Czlx: "2",
		Num:  len(ids),
		Czr:  czr,
		Jgdm: srcJgdm,
	}
	if err := tx.Table("devmod").Create(dm).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, id := range ids {
		d, err := getDevinfoByID(id)
		if err != nil {
			tx.Rollback()
			return err
		}
		dmd := &Devmodetail{
			Lsh:   ckLsh,
			Czlx:  "2",
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
	rkLsh := util.RandomString(4) + strconv.Itoa(int(time.Now().Unix()))
	dm2 := &Devmod{
		Lsh:  rkLsh,
		Czrq: t,
		Czlx: "1",
		Num:  len(ids),
		Czr:  czr,
		Jgdm: dstJgdm,
	}
	if err := tx.Table("devmod").Create(dm2).Error; err != nil {
		tx.Rollback()
		return err
	}
	zt, sx := getState("1")
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
		d, err := getDevinfoByID(id)
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
	return tx.Commit().Error
}

//设备分配&借出
func DevAllocate(ids, dms []string, jgdm, syr, cfwz, czr, czlx string) error {
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
		if err := tx.Table("devinfo").Where("id=?", dev["id"]).Updates(dev).Error; err != nil {
			tx.Rollback()
			return err
		}
		d, err := getDevinfoByID(id)
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
		}
		if err := tx.Table("devmodetail").Create(dmd).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func getState(czlx string) (zt, sx string) {
	switch czlx {
	case "1":
		zt, sx = "1", "1"
	case "2":
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
		zt, sx = "4", "2"
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

//批量导入
type DevinfoErr struct {
	*Devinfo
	Msg string `json:"msg"`
}

func ImpDevinfos(fileName io.Reader, czr string) ([]*DevinfoErr, int, int, error) {
	devs, err := ReadDevinfoXmlToStructs(fileName, czr)
	if err != nil {
		return nil, 0, 0, err
	}
	errDev, success, failed := InsertDevinfoXml(devs, czr)
	return errDev, success, failed, nil
}

func ReadDevinfoXmlToStructs(fileName io.Reader, czr string) ([]*Devinfo, error) {
	devs := make([]*Devinfo, 0)
	xlsx, err := excelize.OpenReader(fileName)
	if err != nil {
		logging.Info(err.Error())
		return nil, err
	}
	sheetName := xlsx.GetSheetName(1)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		logging.Info(err.Error())
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
							info := d.ID + "$序列号[" + d.Xlh + "]$生产商[" + d.Scs + "]$生产日期[" + d.Scrq + "]$"
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
								Lx:    dev.Lx,
								DevID: dev.ID,
								Zcbh:  dev.Zcbh,
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
	Jgmc string `json:"jgmc"`
}

func GetDevinfos(con map[string]string, pageNo, pageSize int, bz string) ([]*DevinfoResp, error) {
	var (
		devs    []*DevinfoResp
		jgdmCon string
		ztCon   string
	)
	offset := (pageNo - 1) * pageSize
	query := `select devinfo.sbbh,devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,devinfo.rkrq,
			devinfo.czrq,user.name as czr,devinfo.qrurl,devstate.mc as zt,devinfo.jgdm,devdept.jgmc,
			devinfo.syr,devinfo.cfwz,devproperty.mc as sx
			from devinfo 
			left join user on user.userid=devinfo.czr 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devdept on devdept.jgdm=devinfo.jgdm 
			left join devproperty on devproperty.dm=devinfo.sx 
			where devinfo.mc like '%%%s%%' and devinfo.rkrq >= '%s' and devinfo.czrq <= '%s'
			and devinfo.id like '%%%s%%' and devinfo.xlh like '%%%s%%' and devinfo.syr like '%%%s%%'
			and devinfo.jgdm %s %s
			order by devinfo.czrq desc LIMIT %d,%d`
	if con["jgdm"] == "" {
		jgdmCon = "like '00%'"
	} else {
		jgdmCon = "= '" + con["jgdm"] + "'"
	}
	if bz == "0" {
		ztCon = " and devinfo.zt = '1'"
	} else if bz == "3" {
		ztCon = " and devinfo.zt = '2' and devinfo.sx = '3'"
	} else if bz == "4" {
		ztCon = " and devinfo.zt = '2' and devinfo.sx = '4'"
	} else if bz == "6" {
		ztCon = " and devinfo.zt = '3' and devinfo.sx = '4'"
	} else if bz == "10" {
		ztCon = " and ((devinfo.zt = '2' and devinfo.sx = '3') or(devinfo.zt = '2' and devinfo.sx = '4')or(devinfo.zt = '3' and devinfo.sx = '4'))"
	}
	squery := fmt.Sprintf(query,
		con["mc"], con["rkrqq"], con["rkrqz"], con["sbbh"], con["xlh"], con["syr"], jgdmCon, ztCon, offset, pageSize)
	if err := db.Raw(squery).Scan(&devs).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return devs, nil
}

func GetDevinfosGly(con map[string]string) ([]*DevinfoResp, error) {
	var (
		devs []*DevinfoResp
	)
	squery := `select devinfo.sbbh,devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,devinfo.rkrq,
			devinfo.czrq,user.name as czr,devinfo.qrurl,devstate.mc as zt,devinfo.jgdm,devdept.jgmc,
			devinfo.syr,devinfo.cfwz,devproperty.mc as sx
			from devinfo 
			left join user on user.userid=devinfo.czr 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devdept on devdept.jgdm=devinfo.jgdm 
			left join devproperty on devproperty.dm=devinfo.sx 
			where devinfo.jgdm like '` + con["jgdm"] + `%' `
	if con["sbbh"] != "" {
		squery += `and devinfo.sbbh = '` + con["sbbh"] + `' `
	}
	if con["property"] != "" {
		squery += `and devinfo.sx = '` + con["property"] + `' `
	}
	if con["state"] != "" {
		squery += `and devinfo.zt = '` + con["state"] + `' `
	}
	if con["type"] != "" {
		squery += `and devinfo.lx = '` + con["type"] + `' `
	}
	if con["xlh"] != "" {
		squery += `and devinfo.xlh = '` + con["xlh"] + `' `
	}
	//log.Println(squery)
	if err := db.Raw(squery).Scan(&devs).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return devs, nil
}

func GetDevinfoByID(id string) (*DevinfoResp, error) {
	var dev DevinfoResp
	query := `select devinfo.sbbh,devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,devinfo.rkrq,
			devinfo.czrq,user.name as czr,devinfo.qrurl,devstate.mc as zt,devinfo.jgdm,devdept.jgmc,
			devinfo.syr,devinfo.cfwz,devproperty.mc as sx
			from devinfo 
			left join devdept on devdept.jgdm=devinfo.jgdm 
			left join user on user.userid=devinfo.czr 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devproperty on devproperty.dm=devinfo.sx 
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

func getDevinfoByID(id string) (*Devinfo, error) {
	var dev Devinfo
	if err := db.Table("devinfo").Where("id=?", id).First(&dev).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
	}
	return nil, nil
}

func GetDevinfoBySbbh(sbbh uint) (*Devinfo, error) {
	var dev Devinfo
	if err := db.Table("devinfo").Where("sbbh=?", sbbh).First(&dev).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
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
