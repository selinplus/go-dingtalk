package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/tealeg/xlsx"
	"log"
)

var db *gorm.DB

// Setup initializes the database instance
func Setup() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	CheckTable()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CheckTable() {
	if !db.HasTable("attachment") {
		db.CreateTable(Attachment{})
	} else {
		db.AutoMigrate(Attachment{})
	}
	if !db.HasTable("msg") {
		db.CreateTable(Msg{})
	} else {
		db.AutoMigrate(Msg{})
	}
	if !db.HasTable("msg_tag") {
		db.CreateTable(MsgTag{})
	} else {
		db.AutoMigrate(MsgTag{})
	}
	if !db.HasTable("department") {
		db.CreateTable(Department{})
	} else {
		db.AutoMigrate(Department{})
	}
	if !db.HasTable("user") {
		db.CreateTable(User{})
	} else {
		db.AutoMigrate(User{})
	}
	if !db.HasTable("device") {
		db.CreateTable(Device{})
	} else {
		db.AutoMigrate(Device{})
	}
	if !db.HasTable("devmodify") {
		db.CreateTable(Devmodify{})
	} else {
		db.AutoMigrate(Devmodify{})
	}
	if !db.HasTable("devoperation") {
		db.CreateTable(Devoperation{})
	} else {
		db.AutoMigrate(Devoperation{})
	}
	if !db.HasTable("devstate") {
		db.CreateTable(Devstate{})
	} else {
		db.AutoMigrate(Devstate{})
	}
	if !db.HasTable("devtype") {
		db.CreateTable(Devtype{})
	} else {
		db.AutoMigrate(Devtype{})
	}
	if !db.HasTable("process") {
		db.CreateTable(Process{})
	} else {
		db.AutoMigrate(Process{})
	}
	if !db.HasTable("procmodify") {
		db.CreateTable(Procmodify{})
	} else {
		db.AutoMigrate(Procmodify{})
	}
	if !db.HasTable("procnode") {
		db.CreateTable(Procnode{})
	} else {
		db.AutoMigrate(Procnode{})
	}
	if !db.HasTable("procstate") {
		db.CreateTable(Procstate{})
	} else {
		db.AutoMigrate(Procstate{})
	}
	if !db.HasTable("proctype") {
		db.CreateTable(Proctype{})
	} else {
		db.AutoMigrate(Proctype{})
	}
}

func InitDb() {
	CheckTable()
	var cnt int
	err := db.Select("id").Model(&Devtype{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init db error: %v", err))
		return
	}
	if cnt == 0 {
		AddType()
	}
	err = db.Select("id").Model(&Devstate{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init db error: %v", err))
		return
	}
	if cnt == 0 {
		AddState()
	}
	err = db.Select("id").Model(&Devoperation{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init db error: %v", err))
		return
	}
	if cnt == 0 {
		AddOpera()
	}
}

func AddType() {
	devType := readXmlToMapType()
	InsertType(devType)
}

func AddState() {
	devState := readXmlToMapState()
	InsertState(devState)
}

func AddOpera() {
	devOpera := readXmlToMapOpera()
	InsertOpera(devOpera)
}

func readXmlToMapType() []map[string]string {
	res := make([]map[string]string, 0)
	inFile := setting.AppSetting.RuntimeRootPath + setting.AppSetting.ExportSavePath + "device.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	for k, sheet := range xlFile.Sheets {
		if k == 1 {
			logging.Info(fmt.Sprintf("sheet name: %s", sheet.Name))
			//遍历行读取
			for r, row := range sheet.Rows {
				if r == 0 {
					continue
				}
				m := make(map[string]string, 0)
				// 遍历每行的列读取
				for i, cell := range row.Cells {
					text := cell.String()
					switch {
					case i == 0:
						m["dm"] = text
					case i == 1:
						m["mc"] = text
					}
				}
				res = append(res, m)
			}
		}
	}
	return res
}

func readXmlToMapState() []map[string]string {
	res := make([]map[string]string, 0)
	inFile := setting.AppSetting.RuntimeRootPath + setting.AppSetting.ExportSavePath + "device.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	for k, sheet := range xlFile.Sheets {
		if k == 2 {
			logging.Info(fmt.Sprintf("sheet name: %s", sheet.Name))
			//遍历行读取
			for r, row := range sheet.Rows {
				if r == 0 {
					continue
				}
				m := make(map[string]string, 0)
				// 遍历每行的列读取
				for i, cell := range row.Cells {
					text := cell.String()
					switch {
					case i == 0:
						m["dm"] = text
					case i == 1:
						m["mc"] = text
					}
				}
				res = append(res, m)
			}
		}
	}
	return res
}

func readXmlToMapOpera() []map[string]string {
	res := make([]map[string]string, 0)
	inFile := setting.AppSetting.RuntimeRootPath + setting.AppSetting.ExportSavePath + "device.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	for k, sheet := range xlFile.Sheets {
		if k == 3 {
			logging.Info(fmt.Sprintf("sheet name: %s", sheet.Name))
			//遍历行读取
			for r, row := range sheet.Rows {
				if r == 0 {
					continue
				}
				m := make(map[string]string, 0)
				// 遍历每行的列读取
				for i, cell := range row.Cells {
					text := cell.String()
					switch {
					case i == 0:
						m["dm"] = text
					case i == 1:
						m["mc"] = text
					}
				}
				res = append(res, m)
			}
		}
	}
	return res
}

func InsertType(devType []map[string]string) {
	er := make([]Devtype, 0)
	logging.Debug(fmt.Sprintf("------------------%d------", len(devType)))
	if len(devType) > 0 {
		for _, d := range devType {
			dev := Devtype{
				Dm: d["dm"],
				Mc: d["mc"],
			}
			err := db.Model(Devtype{}).Create(&dev).Error
			if err != nil {
				er = append(er, dev)
			}
		}
	}
	if len(er) > 0 {
		for _, e := range er {
			logging.Info(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertState(devState []map[string]string) {
	er := make([]Devstate, 0)
	logging.Debug(fmt.Sprintf("------------------%d------", len(devState)))
	if len(devState) > 0 {
		for _, d := range devState {
			dev := Devstate{
				Dm: d["dm"],
				Mc: d["mc"],
			}
			err := db.Model(Devstate{}).Create(&dev).Error
			if err != nil {
				er = append(er, dev)
			}
		}
	}
	if len(er) > 0 {
		for _, e := range er {
			logging.Info(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertOpera(devOpera []map[string]string) {
	er := make([]Devoperation, 0)
	logging.Debug(fmt.Sprintf("------------------%d------", len(devOpera)))
	if len(devOpera) > 0 {
		for _, d := range devOpera {
			dev := Devoperation{
				Dm: d["dm"],
				Mc: d["mc"],
			}
			err := db.Model(Devoperation{}).Create(&dev).Error
			if err != nil {
				er = append(er, dev)
			}
		}
	}
	if len(er) > 0 {
		for _, e := range er {
			logging.Info(fmt.Sprintf("%+v", e))
		}
	}
}
