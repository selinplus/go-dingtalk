package models

import "github.com/jinzhu/gorm"

type Netdisk struct {
	ID       int    `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	UserID   string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	TreeID   int    `json:"tree_id" gorm:"COMMENT:'文件夹id，回收站0，网盘>0'"`
	OrID     int    `json:"orid" gorm:"column:orid;default:'9999';COMMENT:'源文件夹id，存储回收站内文件源位置id'"`
	FileName string `json:"file_name" gorm:"COMMENT:'文件原始名'"`
	FileUrl  string `json:"file_url" gorm:"COMMENT:'文件真实文件名'"`
	FileSize int    `json:"file_size" gorm:"COMMENT:'文件大小'"`
	Xgrq     string `json:"xgrq" gorm:"COMMENT:'修改时间'"`
}

func IsDirContainFile(userid string, id int) bool {
	var nt Netdisk
	if err := db.Where("userid =? and id=?", userid, id).First(&nt).Error; err != nil {
		return false
	}
	if nt.ID > 0 {
		return true
	}
	return false
}

func AddNetdiskFile(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteNetdiskFile(id int) error {
	if err := db.Where("id=?", id).Delete(&Netdisk{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateNetdiskFile(netdisk *Netdisk) error {
	if err := db.Save(netdisk).Error; err != nil {
		return err
	}
	return nil
}

type NetdiskResp struct {
	Netdisk
	Name string `json:"name"`
}

func GetNetdiskFileList(userid string, treeid, pageNum, pageSize int) ([]*NetdiskResp, error) {
	var netdisks []*NetdiskResp
	sql := `SELECT netdisk.*,user.name
			FROM netdisk LEFT JOIN user ON netdisk.userid=user.userid
			WHERE netdisk.userid = ? and netdisk.tree_id=?
			ORDER BY netdisk.xgrq DESC LIMIT ?,?`
	err := db.Raw(sql, userid, treeid, pageNum, pageSize).Scan(&netdisks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return netdisks, nil
}

func GetNetdiskFileDetail(id int) (*Netdisk, error) {
	var netdisk Netdisk
	sql := `SELECT netdisk.*,user.name
			FROM netdisk LEFT JOIN user ON netdisk.userid=user.userid
			WHERE netdisk.id = ?`
	err := db.Raw(sql, id).Scan(&netdisk).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &netdisk, nil
}

func GetTrashFiles() ([]*Netdisk, error) {
	var netdisks []*Netdisk
	err := db.Where("tree_id=0 and DATEDIFF(NOW(),xgrq)>30").Find(&netdisks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return netdisks, nil
}
