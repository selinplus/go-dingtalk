package fsdj

import (
	"fmt"
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

type StudyHltForm struct {
	ID      uint     `json:"id"`
	UserID  string   `json:"userid"`
	ActID   uint     `json:"act_id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	HltUrls []string `json:"hlt_urls"`
	Flag    string   `json:"flag"` //0:图文 1:视频
	Status  string   `json:"status"`
}

//党员风采发布
func PostStudyHlt(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   StudyHltForm
		hltUrl string
		userid string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	if len(c.Query("mobile")) > 0 {
		user, err := models.GetUserByMobile(c.Query("mobile"))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
			return
		}
		userid = user.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	if len(form.HltUrls) > 0 {
		for _, url := range form.HltUrls {
			j := strings.LastIndex(url, "/")
			hltUrl += url[j+1:] + ";"
		}
		hltUrl = strings.TrimRight(hltUrl, ";")
	}
	topic := &models.StudyHlt{
		StudyActID: form.ActID,
		UserID:     userid,
		Title:      form.Title,
		Content:    form.Content,
		HltUrl:     hltUrl,
		Flag:       form.Flag,
		Status:     "0",
		Fbrq:       time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddStudyHlt(topic); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func UpdStudyHlt(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   StudyHltForm
		hltUrl string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	hlt := &models.StudyHlt{
		ID:   form.ID,
		Xgrq: time.Now().Format("2006-01-02 15:04:05"),
	}
	url := c.Request.URL.Path
	if strings.Contains(url, "edit") { //编辑
		if len(form.HltUrls) > 0 {
			for _, url := range form.HltUrls {
				j := strings.LastIndex(url, "/")
				hltUrl += url[j+1:] + ";"
			}
			hltUrl = strings.TrimRight(hltUrl, ";")
		}
		hlt.Title = form.Title
		hlt.Content = form.Content
		hlt.HltUrl = hltUrl
		hlt.Status = "0"
	}
	if strings.Contains(url, "approve") { //审核发布&驳回
		hlt.Status = form.Status // 1:通过(发布)  3:驳回
	}
	if strings.Contains(url, "cancel") { //撤销
		hlt.Status = "2"
	}
	if err := models.UpdStudyHlt(hlt); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type HltResp struct {
	*models.StudyHlt
	Name          string `json:"name"`
	Mobile        string `json:"mobile"`
	StudyHltStars []*HltStarResp
}
type HltStarResp struct {
	*models.StudyHltStar
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//获取党员风采详情
func GetStudyHlt(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		id   = c.Param("id")
		url  = c.Request.URL.Path
	)
	hlt, err := models.GetStudyHlt(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if hlt.ID > 0 {
		if strings.Contains(url, "v1") {
			urls := make([]string, 0)
			if hlt.HltUrl != "" {
				for _, hltUrl := range strings.Split(hlt.HltUrl, ";") {
					urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(hltUrl))
				}
				hlt.HltUrls = urls
			}
			if strings.Contains(hlt.Content, "fsdj/dj_image") {
				if strings.Contains(hlt.Content, "api/fsdj/dj_image") {
					hlt.Content = strings.ReplaceAll(
						hlt.Content, "api/fsdj/dj_image", "fsdj/dj_image")
				}
			}
		} else {
			urls := make([]string, 0)
			if hlt.HltUrl != "" {
				for _, hltUrl := range strings.Split(hlt.HltUrl, ";") {
					urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(hltUrl))
				}
				hlt.HltUrls = urls
			}
			if strings.Contains(hlt.Content, "fsdj/dj_image") {
				if !strings.Contains(hlt.Content, "api/fsdj/dj_image") {
					hlt.Content = strings.ReplaceAll(
						hlt.Content, "fsdj/dj_image", "api/fsdj/dj_image")
				}
			}

			token := c.GetHeader("Authorization")
			auth := c.Query("token")
			if len(auth) > 0 {
				token = auth
			}
			ts := strings.Split(token, ".")
			hlt.Star = models.IsStudyHltStar(hlt.ID, ts[3])
		}
		hlt.StarNum = len(hlt.StudyHltStars)
		stars := make([]*HltStarResp, 0)
		if len(hlt.StudyHltStars) > 0 {
			for _, star := range hlt.StudyHltStars {
				user, err := models.GetUserByUserid(star.UserID)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
					return
				}
				stars = append(stars, &HltStarResp{
					StudyHltStar: &star,
					Avatar:       user.Avatar,
					Name:         user.Name,
					Mobile:       user.Mobile,
				})
			}
		}
		user, err := models.GetUserByUserid(hlt.UserID)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, HltResp{
			StudyHlt:      hlt,
			Name:          user.Name,
			Mobile:        user.Mobile,
			StudyHltStars: stars,
		})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取活动id获取党员风采列表
func GetStudyHlts(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		actId    = c.Query("actId")
		flag     = c.Query("flag")
		dm       = c.Query("dm")
		status   = c.Query("status") //0:未审核 1:审核通过(发布) 2:审核驳回 3:撤销发布
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

	cond := fmt.Sprintf(
		"status like '%s' and flag like '%s'", status+"%", flag+"%")
	if actId != "" {
		cond += fmt.Sprintf(" and act_id='%s'", actId)
	} else {
		cond += " and act_id != '0'"
	}
	if dm != "" {
		cond += fmt.Sprintf(" and study_member.dm='%s'", dm)
	}
	hlts, err := models.GetStudyHlts(cond, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	var data []*HltResp
	if len(hlts) > 0 {
		for _, hlt := range hlts {
			if strings.Contains(url, "v1") {
				urls := make([]string, 0)
				if hlt.HltUrl != "" {
					for _, hltUrl := range strings.Split(hlt.HltUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(hltUrl))
					}
					hlt.HltUrls = urls
				}
				if strings.Contains(hlt.Content, "fsdj/dj_image") {
					if strings.Contains(hlt.Content, "api/fsdj/dj_image") {
						hlt.Content = strings.ReplaceAll(
							hlt.Content, "api/fsdj/dj_image", "fsdj/dj_image")
					}
				}
			} else {
				urls := make([]string, 0)
				if hlt.HltUrl != "" {
					for _, hltUrl := range strings.Split(hlt.HltUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(hltUrl))
					}
					hlt.HltUrls = urls
				}
				if strings.Contains(hlt.Content, "fsdj/dj_image") {
					if !strings.Contains(hlt.Content, "api/fsdj/dj_image") {
						hlt.Content = strings.ReplaceAll(
							hlt.Content, "fsdj/dj_image", "api/fsdj/dj_image")
					}
				}
				token := c.GetHeader("Authorization")
				auth := c.Query("token")
				if len(auth) > 0 {
					token = auth
				}
				ts := strings.Split(token, ".")
				hlt.Star = models.IsStudyHltStar(hlt.ID, ts[3])
			}
			hlt.StarNum = len(hlt.StudyHltStars)
			stars := make([]*HltStarResp, 0)
			if len(hlt.StudyHltStars) > 0 {
				for _, star := range hlt.StudyHltStars {
					user, err := models.GetUserByUserid(star.UserID)
					if err != nil {
						appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
						return
					}
					stars = append(stars, &HltStarResp{
						StudyHltStar: &star,
						Avatar:       user.Avatar,
						Name:         user.Name,
						Mobile:       user.Mobile,
					})
				}
			}
			user, err := models.GetUserByUserid(hlt.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
				return
			}
			data = append(data, &HltResp{
				StudyHlt:      hlt,
				Name:          user.Name,
				Mobile:        user.Mobile,
				StudyHltStars: stars,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": data,
				"cnt":  models.GetStudyHltsCnt(cond),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取userid获取党员风采列表
func GetStudyHltsByUserid(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		status   = c.Query("status") //0:未审核 1:审核通过(发布) 2:撤销发布 3:审核驳回
		flag     = c.Query("flag")   //0:图文 1:视频
		url      = c.Request.URL.Path
		pageSize int
		pageNo   int
		userid   string
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

	if len(c.Query("mobile")) > 0 {
		user, err := models.GetUserByMobile(c.Query("mobile"))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
			return
		}
		userid = user.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	hlts, err := models.GetStudyHltsByUserid(userid, status, flag, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	var data []*HltResp
	if len(hlts) > 0 {
		for _, hlt := range hlts {
			if strings.Contains(url, "v1") {
				urls := make([]string, 0)
				if hlt.HltUrl != "" {
					for _, hltUrl := range strings.Split(hlt.HltUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(hltUrl))
					}
					hlt.HltUrls = urls
				}
				if strings.Contains(hlt.Content, "fsdj/dj_image") {
					if strings.Contains(hlt.Content, "api/fsdj/dj_image") {
						hlt.Content = strings.ReplaceAll(
							hlt.Content, "api/fsdj/dj_image", "fsdj/dj_image")
					}
				}
			} else {
				urls := make([]string, 0)
				if hlt.HltUrl != "" {
					for _, hltUrl := range strings.Split(hlt.HltUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(hltUrl))
					}
					hlt.HltUrls = urls
				}
				if strings.Contains(hlt.Content, "fsdj/dj_image") {
					if !strings.Contains(hlt.Content, "api/fsdj/dj_image") {
						hlt.Content = strings.ReplaceAll(
							hlt.Content, "fsdj/dj_image", "api/fsdj/dj_image")
					}
				}
			}
			hlt.StarNum = len(hlt.StudyHltStars)
			hlt.Star = models.IsStudyHltStar(hlt.ID, userid)
			stars := make([]*HltStarResp, 0)
			if len(hlt.StudyHltStars) > 0 {
				for _, star := range hlt.StudyHltStars {
					user, err := models.GetUserByUserid(star.UserID)
					if err != nil {
						appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
						return
					}
					stars = append(stars, &HltStarResp{
						StudyHltStar: &star,
						Avatar:       user.Avatar,
						Name:         user.Name,
						Mobile:       user.Mobile,
					})
				}
			}
			user, err := models.GetUserByUserid(hlt.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
				return
			}
			data = append(data, &HltResp{
				StudyHlt:      hlt,
				Name:          user.Name,
				Mobile:        user.Mobile,
				StudyHltStars: stars,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": data,
				"cnt":  models.GetStudyHltsCntByUserid(userid, status, flag),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除党员风采
func DelStudyHlt(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		id   = c.Param("id")
	)
	if err := models.DelStudyHlt(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//点赞党员风采
func AddStudyHltStar(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		id     = c.Param("id")
		userid string
	)
	hltId, _ := strconv.Atoi(id)
	if len(c.Query("mobile")) > 0 {
		user, err := models.GetUserByMobile(c.Query("mobile"))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
			return
		}
		userid = user.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	star := &models.StudyHltStar{
		StudyHltID: uint(hltId),
		UserID:     userid,
	}
	if models.IsHltStar(star) {
		appG.Response(http.StatusOK, e.ERROR, "已点赞!")
		return
	}
	star.Stime = time.Now().Format("2006-01-02 15:04:05")
	if err := models.CreateStudyHltStar(star); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//取消点赞党员风采
func CancelStudyHltStar(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		id     = c.Param("id")
		userid string
	)
	hltId, _ := strconv.Atoi(id)
	if len(c.Query("mobile")) > 0 {
		user, err := models.GetUserByMobile(c.Query("mobile"))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
			return
		}
		userid = user.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	if err := models.CancelStudyHltStar(uint(hltId), userid); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
