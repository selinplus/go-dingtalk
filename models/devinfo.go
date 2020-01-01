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
	ID    string `gorm:"primary_key;COMMENT:'设备编号'"`
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
	Czr   string `json:"czr" gorm:"COMMENT:'操作人'"`
	Czrq  string `json:"czrq" gorm:"COMMENT:'操作日期'"`
	Zt    string `json:"zt" gorm:"COMMENT:'设备状态'"`
	Jgdm  string `json:"jgdm" gorm:"COMMENT:'设备管理机构代码'"`
	Syr   string `json:"syr" gorm:"COMMENT:'设备使用人代码'"`
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

func EditDevinfo(dev *Devinfo) error {
	if err := db.Table("devinfo").Where("id=?", dev.ID).Updates(dev).Error; err != nil {
		return err
	}
	return nil
}

//批量导入
func ImpDevinfos(fileName io.Reader, czr string) ([]*Devinfo, int, int) {
	devs := ReadDevinfoXmlToStructs(fileName, czr)
	errDev, success, failed := InsertDevinfoXml(devs, czr)
	return errDev, success, failed
}

func ReadDevinfoXmlToStructs(fileName io.Reader, czr string) []*Devinfo {
	devs := make([]*Devinfo, 0)
	xlsx, err := excelize.OpenReader(fileName)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	sheetName := xlsx.GetSheetName(1)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		logging.Info(err.Error())
		return nil
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
			switch {
			case i == 0:
				d.Zcbh = cell
			case i == 1:
				d.Lx = cell
			case i == 2:
				d.Mc = cell
			case i == 3:
				d.Xh = cell
			case i == 4:
				d.Xlh = cell
			case i == 5:
				d.Ly = cell
			case i == 6:
				d.Scs = cell
			case i == 7:
				d.Scrq = cell
			case i == 8:
				d.Grrq = cell
			case i == 9:
				d.Bfnx = cell
			case i == 10:
				d.Jg = cell
			case i == 11:
				d.Gys = cell
			}
		}
		d.ID = GenerateSbbh(d.Lx, d.Xlh)
		//logging.Debug(fmt.Sprintf("*: %+v", d))
		devs = append(devs, &d)
	}
	return devs
}

func InsertDevinfoXml(devs []*Devinfo, czr string) ([]*Devinfo, int, int) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var (
		errDev   = make([]*Devinfo, 0)
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
					d := &Devinfo{
						ID:   dev.ID,
						Zcbh: dev.Zcbh,
						Lx:   dev.Lx,
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
						Czrq: time.Now().Format("2006-01-02 15:04:05"),
						Czr:  dev.Czr,
						Zt:   "1",
						Jgdm: "00",
						Sx:   "1",
					}
					if IsDevXlhExist(dev.Xlh) {
						logging.Info(fmt.Sprintf("%s:序列号已存在!", dev.Xlh))
						errDev = append(errDev, dev)
					} else {
						//生成二维码
						name, _, err := qrcode.GenerateQrWithLogo(dev.ID, qrcode.GetQrCodeFullPath())
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
		tx.Commit()
	}
	if len(errDev) > 0 {
		return errDev, devsNum - len(errDev), len(errDev)
	}
	return nil, devsNum, 0
}

func GetDevinfos(con map[string]string, pageNo, pageSize int) ([]*Devinfo, error) {
	var devs []*Devinfo
	offset := (pageNo - 1) * pageSize
	query := `select devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,
			devinfo.czrq,user.name as czr,devinfo.qrurl,devstate.mc as zt,devdept.jgmc as jgdm,
			devinfo.syr,devproperty.mc as sx
			from devinfo 
			left join devdept on devdept.jgdm=devinfo.jgdm 
			left join user on user.userid=devinfo.czr 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devproperty on devproperty.dm=devinfo.sx 
			where devinfo.mc like '%%%s%%' and devinfo.czrq > '%s' and devinfo.czrq < '%s'
			and devinfo.id like '%%%s%%' and devinfo.xlh like '%%%s%%' and devinfo.syr like '%%%s%%'
			order by devinfo.czrq desc LIMIT %d,%d`
	squery := fmt.Sprintf(query,
		con["mc"], con["rkrqq"], con["rkrqz"], con["sbbh"], con["xlh"], con["syr"], offset, pageSize)
	if err := db.Raw(squery).Scan(&devs).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return devs, nil
}

func GetDevinfoByID(id string) (*Devinfo, error) {
	var dev Devinfo
	query := `select devinfo.id,devinfo.zcbh,devtype.mc as lx,devinfo.mc,devinfo.xh,devinfo.xlh,devinfo.ly,
			devinfo.scs,devinfo.scrq,devinfo.grrq,devinfo.bfnx,devinfo.jg,devinfo.gys,
			devinfo.czrq,user.name as czr,devinfo.qrurl,devstate.mc as zt,devdept.jgmc as jgdm,
			devinfo.syr,devproperty.mc as sx
			from devinfo 
			left join devdept on devdept.jgdm=devinfo.jgdm 
			left join user on user.userid=devinfo.czr 
			left join devtype on devtype.dm=devinfo.lx 
			left join devstate on devstate.dm=devinfo.zt 
			left join devproperty on devproperty.dm=devinfo.sx 
			where devinfo.id like '%%%s%%'`
	squery := fmt.Sprintf(query, id)
	if err := db.Raw(squery).Scan(&dev).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
	}
	return nil, nil
}
