package models

import (
	"fmt"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/tealeg/xlsx"
	"os"
	"strconv"
	"time"
)

type Device struct {
	ID    string `gorm:"primary_key;COMMENT:'设备编号'"`
	Zcbh  string `json:"zcbh" gorm:"COMMENT:'资产编号'"`
	Lx    int    `json:"lx" gorm:"COMMENT:'设备类型'"`
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
	Zt    int    `json:"zt" gorm:"COMMENT:'设备状态'"`
}

func AddDevice(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
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
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	inFile := setting.AppSetting.ImageSavePath + fileName
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
					d.ID = text + timeStamp
					lx, _ := strconv.Atoi(text)
					d.Lx = lx
				case i == 2:

				case i == 3:
				case i == 4:
				case i == 5:
				case i == 6:
				case i == 7:
				case i == 8:
				case i == 9:
				case i == 10:
				case i == 11:
				case i == 12:
				case i == 13:
				case i == 14:
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
			d := Device{
				Mc: dev.Mc,
			}
			err := db.Model(Device{}).Create(&d).Error
			if err != nil {
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

func EditDevice(data interface{}) error {
	if err := db.Model(&Device{}).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func QrCode(id string) {
	//qr := qrcode.NewQrCode(id, 300, 300, qr.M, qr.Auto)
	//posterName := qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	//poster_url := qrcode.GetQrCodeFullUrl(posterName)
	//poster_save_url := filePath + posterName,
}
