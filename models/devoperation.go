package models

type Devoperation struct {
	ID uint   `gorm:"primary_key"`
	Dm string `json:"dm" gorm:"COMMENT:'操作类型代码'"`
	Mc string `json:"mc" gorm:"COMMENT:'操作类型'"`
}

/*
代码 名称
1	设备入库
2	设备出库
3	分配共用
4	分配专有
5	借出共用
6	借出专有
7	设备收回
8	设备交回 --- 人->库
9	设备报废
10	设备上交 --- 库->库
11	机构变更
12  设备初始化
13	管理员变更
*/
func GetDevOp() ([]*Devoperation, error) {
	var ds []*Devoperation
	if err := db.Table("devoperation").Find(&ds).Error; err != nil {
		return nil, err
	}
	return ds, nil
}
