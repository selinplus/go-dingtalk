package models

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/qrcode"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

type Device struct {
	ID    string `gorm:"primary_key;COMMENT:'设备编号'"`
	Zcbh  string `json:"zcbh" gorm:"COMMENT:'资产编号'"`
	Lx    string `json:"lx" gorm:"COMMENT:'设备类型'"`
	Mc    string `json:"mc" gorm:"COMMENT:'设备名称'"`
	Xh    string `json:"xh" gorm:"COMMENT:'设备型号'"`
	Xlh   string `json:"xlh" gorm:"COMMENT:'序列号'"`
	Ly    string `json:"ly" gorm:"COMMENT:'设备来源'"`
	Scs   string `json:"scs" gorm:"COMMENT:'生产商'"`
	Scrq  string `json:"scrq" gorm:"COMMENT:'生产日期'"`
	Grrq  string `json:"grrq" gorm:"COMMENT:'购入日期'"`
	Bfnx  string `json:"bfnx" gorm:"COMMENT:'设备报废年限'"`
	Jg    string `json:"jg" gorm:"COMMENT:'价格'"`
	Zp    string `json:"zp" gorm:"COMMENT:'设备照片'"`
	Gys   string `json:"gys" gorm:"COMMENT:'供应商'"`
	Rkrq  string `json:"rkrq" gorm:"COMMENT:'入库日期'"`
	QrUrl string `json:"qrurl" gorm:"COMMENT:'二维码URL';column:qrurl"`
	Czr   string `json:"czr" gorm:"COMMENT:'操作人'"`
	Zt    string `json:"zt" gorm:"COMMENT:'设备状态'"`
}

//生成设备编号
func GenerateSbbh(lx, xlh string) string {
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	sbbh := lx + xlh + timeStamp[:13]
	return sbbh
}

//判断序列号是否存在
func IsXlhExist(xlh string) bool {
	var dev Device
	err := db.Table("device").Where("xlh=?", xlh).First(&dev).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func AddDevice(data interface{}) error {
	if err := db.Save(data).Error; err != nil {
		return err
	}
	return nil
}

func EditDevice(dev *Device) error {
	if err := db.Table("device").Where("id=?", dev.ID).Updates(dev).Error; err != nil {
		return err
	}
	return nil
}

//批量导入
func ImpDevices(fileName io.Reader, czr string) ([]*Device, int, int) {
	devs := ReadXmlToStructs(fileName, czr)
	errDev, success, failed := InsertDeviceXml(devs)
	return errDev, success, failed
}

func ReadXmlToStructs(fileName io.Reader, czr string) []*Device {
	devs := make([]*Device, 0)
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
		d := Device{}
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
			case i == 12:
				d.Zt = cell
			}
		}
		d.ID = GenerateSbbh(d.Lx, d.Xlh)
		//logging.Debug(fmt.Sprintf("*: %+v", d))
		devs = append(devs, &d)
	}
	return devs
}

func InsertDeviceXml(devs []*Device) ([]*Device, int, int) {
	var (
		errDev   = make([]*Device, 0)
		devsChan = make(chan *Device)
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
		for k := 0; k < wtNum; k++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for dev := range devsChan {
					//if dev == nil {
					//	break
					//}
					d := Device{
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
						Rkrq: time.Now().Format("2006-01-02 15:04:05"),
						Czr:  dev.Czr,
						Zt:   dev.Zt,
					}
					if IsXlhExist(dev.Xlh) {
						logging.Info(fmt.Sprintf("%s:序列号已存在!", dev.Xlh))
						errDev = append(errDev, dev)
					} else {
						//生成二维码
						name, _, err := qrcode.GenerateQrWithLogo(dev.ID, qrcode.GetQrCodeFullPath())
						if err != nil {
							log.Println(err)
						}
						d.QrUrl = qrcode.GetQrCodeFullUrl(name)
						if err = db.Save(&d).Error; err != nil {
							errDev = append(errDev, dev)
						}
					}
				}
			}()
		}
		wg.Wait()
	}
	if len(errDev) > 0 {
		return errDev, devsNum - len(errDev), len(errDev)
	}
	return nil, devsNum, 0
}

type DevResponse struct {
	ID      string
	Zcbh    string `json:"zcbh"`
	Lx      string `json:"lx"`
	Mc      string `json:"mc"`
	Xh      string `json:"xh"`
	Xlh     string `json:"xlh"`
	Ly      string `json:"ly"`
	Scs     string `json:"scs"`
	Scrq    string `json:"scrq"`
	Grrq    string `json:"grrq"`
	Bfnx    string `json:"bfnx"`
	Jg      string `json:"jg"`
	Zp      string `json:"zp"`
	Gys     string `json:"gys"`
	Rkrq    string `json:"rkrq"`
	Czr     string `json:"czr"`
	Qrurl   string `json:"qrurl"`
	Zt      string `json:"zt"`
	Sydw    string `json:"sydw"`
	Syks    string `json:"syks"`
	SyrName string `json:"syr_name"`
	Syr     string `json:"syr"`
	Cfwz    string `json:"cfwz"`
	Czrq    string `json:"czrq"`
}

func GetDevices(con map[string]string, pageNo, pageSize int) ([]*DevResponse, error) {
	var devs []*DevResponse
	offset := (pageNo - 1) * pageSize
	if err := db.Raw("select device.id,device.zcbh,devtype.mc as lx,device.mc,device.xh,device.xlh,device.ly,device.scs,device.scrq,device.grrq,device.bfnx,device.jg,device.zp,device.gys,device.rkrq,device.czr,device.qrurl,devstate.mc as zt,department.name as sydw,department.name as syks,user.name as syr,devmodify.czrq from device left join devmodify on device.id=devmodify.devid left join department on department.id=devmodify.sydw left join user on user.mobile=devmodify.syr left join devtype on devtype.dm=device.lx left join devstate on devstate.dm=device.zt where (devmodify.zzrq =''or devmodify.zzrq is null) and device.mc like ? and device.rkrq > ? and device.rkrq < ? and device.id like ? and device.xlh like ? and ifnull(devmodify.syr,'') like ? order by device.rkrq desc LIMIT ?,?", "%"+con["mc"]+"%", con["rkrqq"], con["rkrqz"], "%"+con["sbbh"]+"%", "%"+con["xlh"]+"%", "%"+con["syr"]+"%", offset, pageSize).
		Scan(&devs).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return devs, nil
}

func GetDeviceByID(id string) (*Device, error) {
	var dev Device
	if err := db.Find(&dev, "id=?", id).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
	}
	return nil, nil
}

func GetDeviceModByDevID(devid string) (*DevResponse, error) {
	var dev DevResponse
	if err := db.Raw("select device.id,device.zcbh,devtype.mc as lx,device.mc,device.xh,device.xlh,device.ly,device.scs,device.scrq,device.grrq,device.bfnx,device.jg,device.zp,device.gys,device.rkrq,device.czr,devstate.mc as zt,department.name as sydw,department.name as syks,user.name as syr_name,devmodify.syr,devmodify.cfwz,devmodify.czrq from device left join devmodify on device.id=devmodify.devid left join department on department.id=devmodify.sydw left join user on user.mobile=devmodify.syr left join devtype on devtype.dm=device.lx left join devstate on devstate.dm=device.zt where device.id = ? and (devmodify.zzrq ='' or devmodify.zzrq is null) ", devid).
		Scan(&dev).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
	}
	return nil, nil
}
