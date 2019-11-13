package models

type NetdiskTree struct {
	ID       int    `gorm:"primary_key"`
	UserID   string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	Name     string `json:"name" gorm:"COMMENT:'文件夹名称'"`
	Nodeid   int    `json:"nodeid" gorm:"COMMENT:'节点id，表示文件夹层级'"`
	Parentid int    `json:"parentid" gorm:"COMMENT:'父部门id，根节点为1'"`
}

func AddNetdiskTree(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func IsParentDir(userid string, nodeid int) bool {
	var nt NetdiskTree
	if err := db.Where("userid =? and nodeid=?", userid, nodeid).First(&nt).Error; err != nil {
		return false
	}
	return true
}

func DeleteNetdiskTree(id uint) error {
	if err := db.Where("id=?", id).Delete(Netdisk{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateNetdiskTree(netdisk *NetdiskTree) error {
	if err := db.Table("netdisk").Where("id=?", netdisk.ID).Updates(netdisk).Error; err != nil {
		return err
	}
	return nil
}
