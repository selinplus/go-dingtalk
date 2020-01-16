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
	if !db.HasTable("note") {
		db.CreateTable(Note{})
	} else {
		db.AutoMigrate(Note{})
	}
	if !db.HasTable("netdisk") {
		db.CreateTable(Netdisk{})
	} else {
		db.AutoMigrate(Netdisk{})
	}
	if !db.HasTable("netdisk_cap") {
		db.CreateTable(NetdiskCap{})
	} else {
		db.AutoMigrate(NetdiskCap{})
	}
	if !db.HasTable("netdisk_tree") {
		db.CreateTable(NetdiskTree{})
	} else {
		db.AutoMigrate(NetdiskTree{})
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
	if !db.HasTable("devinfo") {
		db.CreateTable(Devinfo{})
	} else {
		db.AutoMigrate(Devinfo{})
	}
	if !db.HasTable("devdept") {
		db.CreateTable(Devdept{})
	} else {
		db.AutoMigrate(Devdept{})
	}
	if !db.HasTable("devuser") {
		db.CreateTable(Devuser{})
	} else {
		db.AutoMigrate(Devuser{})
	}
	if !db.HasTable("devmod") {
		db.CreateTable(Devmod{})
	} else {
		db.AutoMigrate(Devmod{})
	}
	if !db.HasTable("devmodetail") {
		db.CreateTable(Devmodetail{})
	} else {
		db.AutoMigrate(Devmodetail{})
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
	if !db.HasTable("devproperty") {
		db.CreateTable(Devproperty{})
	} else {
		db.AutoMigrate(Devproperty{})
	}
	if !db.HasTable("process") {
		db.CreateTable(Process{})
	} else {
		db.AutoMigrate(Process{})
	}
	if !db.HasTable("process_tag") {
		db.CreateTable(ProcessTag{})
	} else {
		db.AutoMigrate(ProcessTag{})
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
		logging.Error(fmt.Sprintf("init Devtype error: %v", err))
		return
	}
	if cnt == 0 {
		AddType()
	}
	err = db.Select("id").Model(&Devstate{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init Devstate error: %v", err))
		return
	}
	if cnt == 0 {
		AddState()
	}
	err = db.Select("id").Model(&Devoperation{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init Devoperation error: %v", err))
		return
	}
	if cnt == 0 {
		AddOpera()
	}
	err = db.Select("id").Model(&Devproperty{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init Devproperty error: %v", err))
		return
	}
	if cnt == 0 {
		AddProperty()
	}
	err = db.Select("id").Model(&Procnode{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init Procnode error: %v", err))
		return
	}
	if cnt == 0 {
		AddProcNode()
	}
	err = db.Select("id").Model(&Proctype{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init Proctype error: %v", err))
		return
	}
	if cnt == 0 {
		AddProcType()
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

func AddProperty() {
	devProperty := readXmlToMapProperty()
	InsertProperty(devProperty)
}

func AddProcNode() {
	procNode := readXmlToStructProcNode()
	InsertNode(procNode)
}

func AddProcType() {
	procType := readXmlToMapProcType()
	InsertProcType(procType)
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
			for r, row := range sheet.Rows {
				if r == 0 {
					continue
				}
				m := make(map[string]string, 0)
				for i, cell := range row.Cells {
					text := cell.String()
					switch {
					case i == 0:
						m["dm"] = text
					case i == 1:
						m["sjdm"] = text
					case i == 2:
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

func readXmlToMapProperty() []map[string]string {
	res := make([]map[string]string, 0)
	inFile := setting.AppSetting.RuntimeRootPath + setting.AppSetting.ExportSavePath + "device.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	for k, sheet := range xlFile.Sheets {
		if k == 4 {
			logging.Info(fmt.Sprintf("sheet name: %s", sheet.Name))
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

func readXmlToStructProcNode() []*Procnode {
	pns := make([]*Procnode, 0)
	inFile := setting.AppSetting.RuntimeRootPath + "submit.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	for k, sheet := range xlFile.Sheets {
		if k == 0 {
			logging.Info(fmt.Sprintf("sheet name: %s", sheet.Name))
			for r, row := range sheet.Rows {
				if r == 0 {
					continue
				}
				pn := Procnode{}
				for i, cell := range row.Cells {
					text := cell.String()
					switch {
					case i == 0:
						pn.Dm = text
					case i == 1:
						pn.Node = text
					case i == 2:
						pn.Last = text
					case i == 3:
						pn.Next = text
					case i == 4:
						pn.Rname = text
					case i == 5:
						pn.Role = text
					case i == 6:
						pn.Flag = text
					}
				}
				logging.Debug(fmt.Sprintf("*: %+v", pn))
				pns = append(pns, &pn)
			}
		}
	}
	return pns
}

func readXmlToMapProcType() []map[string]string {
	res := make([]map[string]string, 0)
	inFile := setting.AppSetting.RuntimeRootPath + "submit.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		logging.Info(err.Error())
		return nil
	}
	for k, sheet := range xlFile.Sheets {
		if k == 1 {
			logging.Info(fmt.Sprintf("sheet name: %s", sheet.Name))
			for r, row := range sheet.Rows {
				if r == 0 {
					continue
				}
				m := make(map[string]string, 0)
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
				Dm:   d["dm"],
				Sjdm: d["sjdm"],
				Mc:   d["mc"],
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

func InsertProperty(devProperty []map[string]string) {
	er := make([]Devproperty, 0)
	logging.Debug(fmt.Sprintf("------------------%d------", len(devProperty)))
	if len(devProperty) > 0 {
		for _, d := range devProperty {
			dev := Devproperty{
				Dm: d["dm"],
				Mc: d["mc"],
			}
			err := db.Model(Devproperty{}).Create(&dev).Error
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

func InsertNode(pns []*Procnode) {
	er := make([]*Procnode, 0)
	logging.Debug(fmt.Sprintf("------------------%d------", len(pns)))
	if len(pns) > 0 {
		for _, pnx := range pns {
			pn := Procnode{
				Dm:    pnx.Dm,
				Node:  pnx.Node,
				Last:  pnx.Last,
				Next:  pnx.Next,
				Rname: pnx.Rname,
				Role:  pnx.Role,
				Flag:  pnx.Flag,
			}
			err := db.Model(Procnode{}).Create(&pn).Error
			if err != nil {
				er = append(er, pnx)
			}
		}
	}
	if len(er) > 0 {
		for _, e := range er {
			logging.Info(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertProcType(procType []map[string]string) {
	er := make([]Proctype, 0)
	logging.Debug(fmt.Sprintf("------------------%d------", len(procType)))
	if len(procType) > 0 {
		for _, d := range procType {
			pt := Proctype{
				Dm: d["dm"],
				Mc: d["mc"],
			}
			err := db.Model(Proctype{}).Create(&pt).Error
			if err != nil {
				er = append(er, pt)
			}
		}
	}
	if len(er) > 0 {
		for _, e := range er {
			logging.Info(fmt.Sprintf("%+v", e))
		}
	}
}
