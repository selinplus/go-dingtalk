package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/tealeg/xlsx"
	"log"
)

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
	if !db.HasTable("msg_addressbook") {
		db.CreateTable(MsgAddressbook{})
	} else {
		db.AutoMigrate(MsgAddressbook{})
	}
	if !db.HasTable("msg_contacter") {
		db.CreateTable(MsgContacter{})
	} else {
		db.AutoMigrate(MsgContacter{})
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
	if !db.HasTable("department") {
		db.CreateTable(Department{})
	} else {
		db.AutoMigrate(Department{})
	}
	if !db.HasTable("userdemo") {
		db.CreateTable(Userdemo{})
	} else {
		db.AutoMigrate(Userdemo{})
	}
	if !db.HasTable("user") {
		db.CreateTable(User{})
	} else {
		db.AutoMigrate(User{})
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
	if !db.HasTable("devtodo") {
		db.CreateTable(Devtodo{})
	} else {
		db.AutoMigrate(Devtodo{})
	}
	if !db.HasTable("devcheck") {
		db.CreateTable(Devcheck{})
	} else {
		db.AutoMigrate(Devcheck{})
	}
	if !db.HasTable("devckdetail") {
		db.CreateTable(Devckdetail{})
	} else {
		db.AutoMigrate(Devckdetail{})
	}
	if !db.HasTable("devcktodd") {
		db.CreateTable(Devcktodd{})
	} else {
		db.AutoMigrate(Devcktodd{})
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
	if !db.HasTable("onduty") {
		db.CreateTable(Onduty{})
	} else {
		db.AutoMigrate(Onduty{})
	}
	if !db.HasTable("ydksworkrecord") {
		db.CreateTable(Ydksworkrecord{})
	} else {
		db.AutoMigrate(Ydksworkrecord{})
	}
	if !db.HasTable("ydksdata") {
		db.CreateTable(Ydksdata{})
	} else {
		db.AutoMigrate(Ydksdata{})
	}
	if !db.HasTable("study_group") {
		db.CreateTable(StudyGroup{})
	} else {
		db.AutoMigrate(StudyGroup{})
	}
	if !db.HasTable("study_member") {
		db.CreateTable(StudyMember{})
	} else {
		db.AutoMigrate(StudyMember{})
	}
	if !db.HasTable("study_act") {
		db.CreateTable(StudyAct{})
	} else {
		db.AutoMigrate(StudyAct{})
	}
	if !db.HasTable("study_actdetail") {
		db.CreateTable(StudyActdetail{})
	} else {
		db.AutoMigrate(StudyActdetail{})
	}
	if !db.HasTable("study_hlt") {
		db.CreateTable(StudyHlt{})
	} else {
		db.AutoMigrate(StudyHlt{})
	}
	if !db.HasTable("study_hlt_star") {
		db.CreateTable(StudyHltStar{})
	} else {
		db.AutoMigrate(StudyHltStar{})
	}
	if !db.HasTable("study_signin") {
		db.CreateTable(StudySignin{})
	} else {
		db.AutoMigrate(StudySignin{})
	}
	if !db.HasTable("study_topic") {
		db.CreateTable(StudyTopic{})
	} else {
		db.AutoMigrate(StudyTopic{})
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
	err = db.Select("id").Model(&StudyGroup{}).Count(&cnt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init StudyGroup error: %v", err))
		return
	}
	if cnt == 0 {
		initGrouproot()
	}
	var sql = `
create or replace view v_devgljg as (
    select t.jgdm                                                       jgksdm,
           t.jgmc                                                       jgks,
           if(t.gly is not null and t.gly <> '', t.jgdm,
              if(b.gly is not null and b.gly <> '', b.jgdm, c.jgdm)) as jgdm,
           if(t.gly is not null and t.gly <> '', t.jgmc,
              if(b.gly is not null and b.gly <> '', b.jgmc, c.jgmc)) as jgmc
    from devdept t
             left join devdept b on t.sjjgdm = b.jgdm
             left join devdept c on b.sjjgdm = c.jgdm);`
	err = db.Exec(sql).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Error(fmt.Sprintf("init v_devgljg error: %v", err))
		return
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
		logging.Error(err.Error())
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
		logging.Error(err.Error())
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
		logging.Error(err.Error())
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
		logging.Error(err.Error())
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
		logging.Error(err.Error())
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
		logging.Error(err.Error())
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
	log.Printf("Devtype------------------%d------", len(devType))
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
			logging.Error(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertState(devState []map[string]string) {
	er := make([]Devstate, 0)
	log.Printf("Devstate------------------%d------", len(devState))
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
			logging.Error(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertOpera(devOpera []map[string]string) {
	er := make([]Devoperation, 0)
	log.Printf("Devoperation------------------%d------", len(devOpera))
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
			logging.Error(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertProperty(devProperty []map[string]string) {
	er := make([]Devproperty, 0)
	log.Printf("Devproperty------------------%d------", len(devProperty))
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
			logging.Error(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertNode(pns []*Procnode) {
	er := make([]*Procnode, 0)
	log.Printf("Procnode------------------%d------", len(pns))
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
			logging.Error(fmt.Sprintf("%+v", e))
		}
	}
}

func InsertProcType(procType []map[string]string) {
	er := make([]Proctype, 0)
	log.Printf("Proctype------------------%d------", len(procType))
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
			logging.Error(fmt.Sprintf("%+v", e))
		}
	}
}
