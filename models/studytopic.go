package models

//党建图文
type StudyTopic struct {
	ID         uint     `gorm:"primary_key;size:11;AUTO_INCREMENT"`
	TopicImage string   `json:"topic_image" gorm:"COMMENT:'主题图片'"`
	Title      string   `json:"title" gorm:"COMMENT:'标题'"`
	Content    string   `json:"content" gorm:"COMMENT:'内容';size:65535"`
	ImageUrl   string   `json:"image_url" gorm:"COMMENT:'图片真实路径';size:65535"`
	ImageUrls  []string `json:"image_urls" gorm:"-"` //返回前台[]Url
	Fbrq       string   `json:"fbrq" gorm:"COMMENT:'发布日期'"`
	Xgrq       string   `json:"xgrq" gorm:"COMMENT:'修改日期'"`
	Rflag      string   `json:"rflag" gorm:"COMMENT:'轮播标志,0:否 1:轮播';default:'0'"`
	Share      string   `json:"share" gorm:"COMMENT:'分享标志,0:否 1:分享';default:'0'"`
	Status     string   `json:"status" gorm:"COMMENT:'状态,0:未审核 1:审核通过(发布) 2:撤销发布';default:'0'"`
}

func AddStudyTopic(data interface{}) error {
	if err := db.Create(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdStudyTopic(topic *StudyTopic) error {
	if err := db.Table("study_topic").
		Where("id=?", topic.ID).Updates(topic).Error; err != nil {
		return err
	}
	return nil
}

func DelStudyTopic(id string) error {
	if err := db.Where("id=?", id).Delete(StudyTopic{}).Error; err != nil {
		return err
	}
	return nil
}

func GetStudyTopic(id string) (*StudyTopic, error) {
	var topic StudyTopic
	if err := db.Where("id=?", id).First(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func GetStudyTopics(rflag, share, status string, pageNo, pageSize int) ([]*StudyTopic, error) {
	var topics []*StudyTopic
	err := db.
		Where("rflag like ? and share like ? and status like ?",
			rflag+"%", share+"%", status+"%").
		Limit(pageSize).Offset(pageSize * (pageNo - 1)).Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func GetStudyTopicsCnt(rflag, share, status string) (cnt int) {
	err := db.
		Where("rflag like ? and share like ? and status like ?",
			rflag+"%", share+"%", status+"%").Count(&cnt).Error
	if err != nil {
		cnt = 0
	}
	return cnt
}
