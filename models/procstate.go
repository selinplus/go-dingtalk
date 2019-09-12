package models

type Procstate struct {
	ID uint   `gorm:"primary_key"`
	Dm string `json:"dm" gorm:"COMMENT:'流程状态代码'"`
	Mc string `json:"mc" gorm:"COMMENT:'流程状态'"`
}

func GetProcstate() ([]*Procstate, error) {
	var ps []*Procstate
	if err := db.Table("procstate").Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps, nil
}
