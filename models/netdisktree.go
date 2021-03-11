package models

type NetdiskTree struct {
	ID     int    `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	PId    int    `json:"pId" gorm:"column:pId;COMMENT:'父部门id，根节点为1,回收站为0'"`
	Name   string `json:"name" gorm:"COMMENT:'文件夹名称'"`
	UserID string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
}

func AddNetdiskDir(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func IsDirExist(userid string, id int) bool {
	var nt NetdiskTree
	if err := db.Where("userid =? and id=?", userid, id).
		First(&nt).Error; err != nil {
		return false
	}
	return true
}

func IsParentDir(userid string, id int) bool {
	var nt NetdiskTree
	if err := db.Where("userid =? and pid=?", userid, id).
		First(&nt).Error; err != nil {
		return false
	}
	return true
}

func DeleteNetdiskDir(userid string, id int) error {
	if err := db.Where("userid=? and id=?", userid, id).
		Delete(&NetdiskTree{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateNetdiskDir(netdisk *NetdiskTree) error {
	if err := db.Table("netdisk_tree").
		Where("id=?", netdisk.ID).Updates(netdisk).Error; err != nil {
		return err
	}
	return nil
}

type Tree struct {
	ID       int     `json:"id"`
	Label    string  `json:"label"`
	UserID   string  `json:"userid"`
	Children []*Tree `json:"children"`
}

//获取文件树列表
func GetNetdiskTree(userID string) ([]Tree, error) {
	perms := make([]Tree, 0)
	child := Tree{
		ID:       1,
		Label:    "我的网盘",
		UserID:   userID,
		Children: []*Tree{},
	}
	err := getTreeNode(1, userID, &child)
	if err != nil {
		return nil, err
	}
	perms = append(perms, child)
	return perms, nil
}

//递归获取子节点
func getTreeNode(pId int, userID string, tree *Tree) error {
	var perms []*NetdiskTree
	err := db.Where("pId=? and userid=?", pId, userID).
		Find(&perms).Error //根据父结点Id查询数据表，获取相应的子结点信息
	if err != nil {
		return err
	}
	for i := 0; i < len(perms); i++ {
		child := Tree{
			ID:       perms[i].ID,
			Label:    perms[i].Name,
			UserID:   perms[i].UserID,
			Children: []*Tree{},
		}
		tree.Children = append(tree.Children, &child)
		err = getTreeNode(perms[i].ID, perms[i].UserID, &child)
	}
	return err
}
