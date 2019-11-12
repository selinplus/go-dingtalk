package models

import "github.com/jinzhu/gorm"

type Note struct {
	ID      uint   `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	Title   string `json:"title" gorm:"COMMENT:'标题'"`
	Content string `json:"content" gorm:"COMMENT:'内容';size:1000"`
	UserID  string `json:"userid" gorm:"column:userid;COMMENT:'用户标识'"`
	Xgrq    string `json:"xgrq" gorm:"COMMENT:'修改日期'"`
}

func IsSameTitle(userid, title string) bool {
	var note Note
	if err := db.Where("userid =? and title=?", userid, title).First(&note).Error; err != nil {
		return false
	}
	return true
}

func AddNote(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteNote(id uint) error {
	if err := db.Where("id=?", id).Delete(Note{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateNote(note *Note) error {
	if err := db.Table("note").Where("id=?", note.ID).Updates(note).Error; err != nil {
		return err
	}
	return nil
}

func GetNoteList(userid string, pageNum, pageSize int) ([]*Note, error) {
	var notes []*Note
	sql := `SELECT note.id,note.title,note.content,user.name,note.xgrq
			FROM note LEFT JOIN user ON note.userid=user.userid
			WHERE note.userid = ? 
			ORDER BY note.xgrq DESC LIMIT ?,?`
	err := db.Raw(sql, userid, pageNum, pageSize).Scan(&notes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return notes, nil
}

func GetNoteDetail(id uint) (*Note, error) {
	var note Note
	sql := `SELECT note.id,note.title,note.content,user.name,note.xgrq
			FROM note LEFT JOIN user ON note.userid=user.userid
			WHERE note.id = ?`
	err := db.Raw(sql, id).Scan(&note).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &note, nil
}
