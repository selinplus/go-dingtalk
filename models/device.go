package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/qrcode"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strconv"
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

func AddDevice(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func EditDevice(data interface{}) error {
	if err := db.Model(&Device{}).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func ImpDevices(fileName string) []*Device {
	devs := ReadXmlToStructs(fileName)
	errDev := InsertDeviceXml(devs)
	return errDev
}

func ReadXmlToStructs(fileName string) []*Device {
	dev := make([]*Device, 0)
	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	inFile := setting.AppSetting.RuntimeRootPath + setting.AppSetting.ImageSavePath + fileName
	xlFile, err := xlsx.OpenFile(inFile)
	defer os.Remove(inFile)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	for _, sheet := range xlFile.Sheets {
		logging.Info(fmt.Sprintf("sheet name: %s", sheet.Name))
		//遍历行读取
		for k, row := range sheet.Rows {
			// 跳过标题行，遍历每行的列读取
			if k == 0 {
				continue
			}
			d := Device{}
			for i, cell := range row.Cells {
				text := cell.String()
				switch {
				case i == 0:
					d.Zcbh = text
				case i == 1:
					d.ID = text + "_" + timeStamp + "_" + d.Zcbh
					d.Lx = text
				case i == 2:
					d.Mc = text
				case i == 3:
					d.Xh = text
				case i == 4:
					d.Xlh = text
				case i == 5:
					d.Ly = text
				case i == 6:
					d.Scs = text
				case i == 7:
					d.Scrq = text
				case i == 8:
					d.Grrq = text
				case i == 9:
					d.Bfnx = text
				case i == 10:
					d.Jg = text
				case i == 11:
					d.Gys = text
				case i == 12:
					d.Rkrq = text
				case i == 13:
					d.Czr = text
				case i == 14:
					d.Zt = text
				}
			}
			logging.Debug(fmt.Sprintf("*: %+v", d))
			dev = append(dev, &d)
		}
	}
	return dev
}

func InsertDeviceXml(devs []*Device) []*Device {
	errDev := make([]*Device, 0)
	logging.Debug(fmt.Sprintf("------------------%d------", len(devs)))
	if len(devs) > 0 {
		for _, dev := range devs {
			//生成二维码
			//qrc := qrcode.NewQrCode(dev.ID, 300, 300, qr.M, qr.Auto)
			//name, _, err := qrc.Encode(qrcode.GetQrCodeFullPath())
			name, _, err := qrcode.GenerateQrWithLogo(dev.ID, qrcode.GetQrCodeFullPath())
			if err != nil {
				log.Println(err)
			}
			d := Device{
				ID:    dev.ID,
				Zcbh:  dev.Zcbh,
				Lx:    dev.Lx,
				Mc:    dev.Mc,
				Xh:    dev.Xh,
				Xlh:   dev.Xlh,
				Ly:    dev.Ly,
				Scs:   dev.Scs,
				Scrq:  dev.Scrq,
				Grrq:  dev.Grrq,
				Bfnx:  dev.Bfnx,
				Jg:    dev.Jg,
				Gys:   dev.Gys,
				Rkrq:  dev.Rkrq,
				QrUrl: qrcode.GetQrCodeFullUrl(name),
				Czr:   dev.Czr,
				Zt:    dev.Zt,
			}
			errd := db.Model(Device{}).Create(&d).Error
			if errd != nil {
				errDev = append(errDev, dev)
			}
		}
	}
	if len(errDev) > 0 {
		for _, e := range errDev {
			logging.Info(fmt.Sprintf("%+v", e))
		}
		return errDev
	}
	return nil
}

func GetDevices(mc string, pageNo, pageSize int) ([]*Device, error) {
	var devs []*Device
	offset := (pageNo - 1) * pageSize
	err := db.Raw("select * from device where mc like ? LIMIT ?,?", "%"+mc+"%", offset, pageSize).
		Scan(&devs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return devs, nil
}

func GetDevicesCount(mc string) (int, error) {
	var total int
	if err := db.Table("device").Where("mc like ?", "%"+mc+"%").Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
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

type DevResponse struct {
	ID   string
	Zcbh string `json:"zcbh"`
	Lx   string `json:"lx"`
	Mc   string `json:"mc"`
	Xh   string `json:"xh"`
	Xlh  string `json:"xlh"`
	Zt   string `json:"zt"`
	Sydw string `json:"sydw"`
	Syks string `json:"syks"`
	Syr  string `json:"syr"`
}

func GetDeviceModByDevID(devid string) (*DevResponse, error) {
	var dev DevResponse
	if err := db.Raw("select device.id,device.zcbh,device.lx,device.mc,device.xh,device.xlh,device.zt,devmodify.sydw,devmodify.syks,devmodify.syr from device left join devmodify on device.id=devmodify.devid left join devtype on devtype.dm=device.lx where device.id = ? and (devmodify.zzrq ='' or devmodify.zzrq is null) ", devid).
		Scan(&dev).Error; err != nil {
		return nil, err
	}
	if len(dev.ID) > 0 {
		return &dev, nil
	}
	return nil, nil
}
