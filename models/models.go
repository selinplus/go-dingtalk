package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/selinplus/go-dingtalk/pkg/setting"
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

	if setting.DatabaseSetting.LogMode == 1 {
		db.LogMode(true)
	}
	db.SingularTable(true)
	CheckTable()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func QueryData(squery string) ([]map[string]string, error) {
	rows, err := db.Raw(squery).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range cols {
		scans[i] = &vals[i]
	}
	var results []map[string]string
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return results, err
		}
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	return results, nil
}

func ExecuSql(sql string) error {
	err := db.Exec(sql).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return nil
}
