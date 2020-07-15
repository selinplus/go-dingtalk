package fsdj

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/fsdjsrv"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StudyTopicForm struct {
	ID         uint     `json:"id"`
	TopicImage string   `json:"topic_image"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	ImageUrls  []string `json:"image_urls"`
	Rflag      string   `json:"rflag"`
	Share      string   `json:"share"`
	Status     string   `json:"status"`
}

//党建图文发布
func PostStudyTopic(c *gin.Context) {
	var (
		appG       = app.Gin{C: c}
		form       StudyTopicForm
		topicImage string
		imageUrl   string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if len(form.TopicImage) > 0 {
		i := strings.LastIndex(form.TopicImage, "/")
		topicImage = form.TopicImage[i+1:]
	}
	if len(form.ImageUrls) > 0 {
		for _, url := range form.ImageUrls {
			j := strings.LastIndex(url, "/")
			imageUrl += url[j+1:] + ";"
		}
		imageUrl = strings.TrimRight(imageUrl, ";")
	}
	topic := &models.StudyTopic{
		TopicImage: topicImage,
		Rflag:      form.Rflag,
		Title:      form.Title,
		Content:    form.Content,
		ImageUrl:   imageUrl,
		Share:      form.Share,
		Status:     "0",
		Fbrq:       time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddStudyTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func UpdStudyTopic(c *gin.Context) {
	var (
		appG       = app.Gin{C: c}
		form       StudyTopicForm
		topicImage string
		imageUrl   string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	topic := &models.StudyTopic{
		ID:   form.ID,
		Xgrq: time.Now().Format("2006-01-02 15:04:05"),
	}
	if len(form.TopicImage) > 0 {
		i := strings.LastIndex(form.TopicImage, "/")
		topicImage = form.TopicImage[i+1:]
	}
	if len(form.ImageUrls) > 0 {
		for _, url := range form.ImageUrls {
			j := strings.LastIndex(url, "/")
			imageUrl += url[j+1:] + ";"
		}
		imageUrl = strings.TrimRight(imageUrl, ";")
	}
	url := c.Request.URL.Path
	if strings.Contains(url, "edit") { //编辑
		topic.TopicImage = topicImage
		topic.Title = form.Title
		topic.Content = form.Content
		topic.ImageUrl = imageUrl
		topic.Rflag = form.Rflag
	}
	if strings.Contains(url, "approve") { //审核发布
		topic.Status = "1"
	}
	if strings.Contains(url, "share") { //分享
		topic.Share = form.Share
	}
	if strings.Contains(url, "cancel") { //撤销
		topic.Status = "2"
	}
	if err := models.UpdStudyTopic(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取图文详情
func GetStudyTopic(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		url  = c.Request.URL.Path
		id   = c.Param("id")
	)

	topic, err := models.GetStudyTopic(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if topic.ID > 0 {
		if strings.Contains(url, "v1") {
			if topic.TopicImage != "" {
				topic.TopicImage = fsdjsrv.GetFsdjImageFullUrl(topic.TopicImage)
			}
			urls := make([]string, 0)
			if topic.ImageUrl != "" {
				for _, imageUrl := range strings.Split(topic.ImageUrl, ";") {
					urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(imageUrl))
				}
				topic.ImageUrls = urls
			}
			if strings.Contains(topic.Content, "fsdj/dj_image") {
				if strings.Contains(topic.Content, "api/fsdj/dj_image") {
					topic.Content = strings.ReplaceAll(
						topic.Content, "api/fsdj/dj_image", "fsdj/dj_image")
				}
			}
		} else {
			if topic.TopicImage != "" {
				topic.TopicImage = fsdjsrv.GetFsdjEappImageFullUrl(topic.TopicImage)
			}
			urls := make([]string, 0)
			if topic.ImageUrl != "" {
				for _, imageUrl := range strings.Split(topic.ImageUrl, ";") {
					urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(imageUrl))
				}
				topic.ImageUrls = urls
			}
			if strings.Contains(topic.Content, "fsdj/dj_image") {
				if !strings.Contains(topic.Content, "api/fsdj/dj_image") {
					topic.Content = strings.ReplaceAll(
						topic.Content, "fsdj/dj_image", "api/fsdj/dj_image")
				}
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS, topic)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取图文列表
func GetStudyTopics(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		err      error
		rflag    = c.Query("rflag")
		share    = c.Query("share")
		status   = c.Query("status") //0:未审核 1:审核通过(发布) 2:撤销发布
		url      = c.Request.URL.Path
		pageSize int
		pageNo   int
	)
	if c.Query("pageNo") == "" {
		pageNo = 1
	} else {
		pageNo, _ = strconv.Atoi(c.Query("pageNo"))
	}
	if c.Query("pageSize") == "" {
		pageSize = 10
	} else {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	topics, err := models.GetStudyTopics(rflag, share, status, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(topics) > 0 {
		for _, topic := range topics {
			if strings.Contains(url, "v1") {
				if topic.TopicImage != "" {
					topic.TopicImage = fsdjsrv.GetFsdjImageFullUrl(topic.TopicImage)
				}
				urls := make([]string, 0)
				if topic.ImageUrl != "" {
					for _, imageUrl := range strings.Split(topic.ImageUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(imageUrl))
					}
					topic.ImageUrls = urls
				}
				if strings.Contains(topic.Content, "fsdj/dj_image") {
					if strings.Contains(topic.Content, "api/fsdj/dj_image") {
						topic.Content = strings.ReplaceAll(
							topic.Content, "api/fsdj/dj_image", "fsdj/dj_image")
					}
				}
			} else {
				if topic.TopicImage != "" {
					topic.TopicImage = fsdjsrv.GetFsdjEappImageFullUrl(topic.TopicImage)
				}
				urls := make([]string, 0)
				if topic.ImageUrl != "" {
					for _, imageUrl := range strings.Split(topic.ImageUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(imageUrl))
					}
					topic.ImageUrls = urls
				}
				if strings.Contains(topic.Content, "fsdj/dj_image") {
					if !strings.Contains(topic.Content, "api/fsdj/dj_image") {
						topic.Content = strings.ReplaceAll(
							topic.Content, "fsdj/dj_image", "api/fsdj/dj_image")
					}
				}
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": topics,
				"cnt":  models.GetStudyTopicsCnt(rflag, share, status),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除党建图文
func DelStudyTopic(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	if err := models.DelStudyTopic(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
