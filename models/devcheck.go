package models

type Devcheck struct {
	ID     uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Beg    string `json:"beg" gorm:"COMMENT:'时间起'"`
	End    string `json:"end" gorm:"COMMENT:'时间止'"`
	Fqr    string `json:"fqr" gorm:"COMMENT:'发起人代码'"`
	Ckself string `json:"ckself" gorm:"COMMENT:'是否自我盘点（Y,N）'"`
}
