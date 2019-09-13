package models

type Procnode struct {
	ID    uint   `gorm:"primary_key"`
	Dm    string `json:"dm" gorm:"COMMENT:'提报类型代码'"`
	Node  string `json:"node" gorm:"COMMENT:'当前节点代码'"`
	Last  string `json:"last" gorm:"COMMENT:'上一节点（0表示起始）'"`
	Next  string `json:"next" gorm:"COMMENT:'下一节点（-1表示终结）'"`
	Rname string `json:"rname" gorm:"COMMENT:'节点人员代码集'"`
	Role  string `json:"role" gorm:"COMMENT:'节点操作人员角色'"`
	Flag  string `json:"flag" gorm:"COMMENT:'是否自动跳过（0不跳过，1跳过）';default:'0'"`
}

//获取节点信息
func GetNode(dm, node string) (*Procnode, error) {
	var pn Procnode
	if err := db.Where("dm=? and node=?", dm, node).First(&pn).Error; err != nil {
		return nil, err
	}
	return &pn, nil
}

//获取上一节点信息
func GetLastNode(dm, next string) (*Procnode, error) {
	var pn Procnode
	if err := db.Where("dm=? and next=?", dm, next).First(&pn).Error; err != nil {
		return nil, err
	}
	return &pn, nil
}

//获取下一节点信息
func GetNextNode(dm, last string) (*Procnode, error) {
	var pn Procnode
	if err := db.Where("dm=? and last=?", dm, last).First(&pn).Error; err != nil {
		return nil, err
	}
	return &pn, nil
}
