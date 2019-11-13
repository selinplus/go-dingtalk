package models

import "github.com/jinzhu/gorm"

type Netdisk struct {
	ID       int    `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	UserID   string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	TreeID   int    `json:"tree_id" gorm:"primary_key;COMMENT:'部门id'"`
	FileName string `json:"file_name" gorm:"COMMENT:'文件原始名'"`
	FileUrl  string `json:"file_url" gorm:"COMMENT:'文件真实路径'"`
	FileSize int    `json:"file_size" gorm:"COMMENT:'文件大小';size:20"`
	Xgrq     string `json:"xgrq" gorm:"COMMENT:'修改时间'"`
	Tag      int    `json:"tag" gorm:"COMMENT:'0：已删除 1: 网盘文件';type:varchar(255);type:int(11);default:'1'"`
}

func GetNetdiskFileDir(userid string, id int) (int, error) {
	var nt Netdisk
	if err := db.Where("userid =? and id=?", userid, id).First(&nt).Error; err != nil {
		return 0, err
	}
	return nt.TreeID, nil
}

func AddNetdiskFile(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteNetdiskFile(id uint) error {
	if err := db.Where("id=?", id).Delete(Netdisk{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateNetdiskFile(netdisk *Netdisk) error {
	if err := db.Table("netdisk").Where("id=?", netdisk.ID).Updates(netdisk).Error; err != nil {
		return err
	}
	return nil
}

func GetNetdiskFileList(userid string, tag, treeid, pageNum, pageSize int) ([]*Netdisk, error) {
	var netdisks []*Netdisk
	sql := `SELECT netdisk.id,netdisk.userid,netdisk.file_name,netdisk.file_url,
					netdisk.file_url,user.name,netdisk.xgrq,netdisk.tag
			FROM netdisk LEFT JOIN user ON netdisk.userid=user.userid
			WHERE netdisk.userid = ? and netdisk.tag=? and netdisk.tree_id=?
			ORDER BY netdisk.xgrq DESC LIMIT ?,?`
	err := db.Raw(sql, userid, tag, treeid, pageNum, pageSize).Scan(&netdisks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return netdisks, nil
}

type NetdiskResp struct {
	Netdisk
	Name string `json:"name"`
}

func GetNetdiskFileDetail(id int) (*Netdisk, error) {
	var netdisk Netdisk
	sql := `SELECT netdisk.id,netdisk.userid,netdisk.file_name,netdisk.file_url,
					netdisk.file_url,user.name,netdisk.xgrq,netdisk.tag
			FROM netdisk LEFT JOIN user ON netdisk.userid=user.userid
			WHERE netdisk.id = ?`
	err := db.Raw(sql, id).Scan(&netdisk).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &netdisk, nil
}
