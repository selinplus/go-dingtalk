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

type StudyActForm struct {
	ID         uint     `json:"id"`
	TopicImage string   `json:"topic_image"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	ImageUrls  []string `json:"image_urls"`
	Share      string   `json:"share"`
	Status     string   `json:"status"`
	Deadline   string   `json:"deadline"`
}

//党建活动发布
func PostStudyAct(c *gin.Context) {
	var (
		appG       = app.Gin{C: c}
		form       StudyActForm
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
	act := &models.StudyAct{
		TopicImage: topicImage,
		Title:      form.Title,
		Content:    form.Content,
		ImageUrl:   imageUrl,
		Share:      form.Share,
		Status:     "0",
		Fbrq:       time.Now().Format("2006-01-02 15:04:05"),
		Deadline:   form.Deadline,
	}
	if err := models.AddStudyAct(act); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func UpdStudyAct(c *gin.Context) {
	var (
		appG       = app.Gin{C: c}
		form       StudyActForm
		topicImage string
		imageUrl   string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	act := &models.StudyAct{
		ID:   form.ID,
		Xgrq: time.Now().Format("2006-01-02 15:04:05"),
	}
	url := c.Request.URL.Path
	if strings.Contains(url, "edit") { //编辑
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
		act.TopicImage = topicImage
		act.Title = form.Title
		act.Content = form.Content
		act.ImageUrl = imageUrl
		act.Deadline = form.Deadline
	}
	if strings.Contains(url, "approve") { //审核发布
		act.Status = "1"
	}
	if strings.Contains(url, "share") { //分享
		act.Share = form.Share
	}
	if strings.Contains(url, "cancel") { //撤销
		act.ID = form.ID
		act.Status = "2"
	}
	if err := models.UpdStudyAct(act); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type ActResp struct {
	*models.StudyAct
	StudyHlts []*HltResp
}

//获取活动详情
func GetStudyAct(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		id     = c.Param("id")
		url    = c.Request.URL.Path
		userid = ""
	)
	act, err := models.GetStudyAct(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if act.ID > 0 {
		if strings.Contains(url, "v1") {
			if act.TopicImage != "" {
				act.TopicImage = fsdjsrv.GetFsdjImageFullUrl(act.TopicImage)
			}
			urls := make([]string, 0)
			if act.ImageUrl != "" {
				for _, imageUrl := range strings.Split(act.ImageUrl, ";") {
					urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(imageUrl))
				}
				act.ImageUrls = urls
			}
			if strings.Contains(act.Content, "fsdj/dj_image") {
				if strings.Contains(act.Content, "api/fsdj/dj_image") {
					act.Content = strings.ReplaceAll(
						act.Content, "api/fsdj/dj_image", "fsdj/dj_image")
				}
			}
		} else {
			if act.TopicImage != "" {
				act.TopicImage = fsdjsrv.GetFsdjEappImageFullUrl(act.TopicImage)
			}
			urls := make([]string, 0)
			if act.ImageUrl != "" {
				for _, imageUrl := range strings.Split(act.ImageUrl, ";") {
					urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(imageUrl))
				}
				act.ImageUrls = urls
			}
			if strings.Contains(act.Content, "fsdj/dj_image") {
				if !strings.Contains(act.Content, "api/fsdj/dj_image") {
					act.Content = strings.ReplaceAll(
						act.Content, "fsdj/dj_image", "api/fsdj/dj_image")
				}
			}

			token := c.GetHeader("Authorization")
			auth := c.Query("token")
			if len(auth) > 0 {
				token = auth
			}
			ts := strings.Split(token, ".")
			userid = ts[3]
		}

		act.JoinNum = len(act.StudyActdetails)
		joined := false //参与标志
		if len(userid) > 0 {
			if models.IsJoinStrudyAct(act.ID, userid) == "Y" {
				joined = true
			}
		}

		act.Joined = joined
		hltResps := make([]*HltResp, 0)
		hlts := act.StudyHlts
		if len(hlts) > 0 {
			for _, hltp := range hlts {
				hlt := hltp
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
				user, err := models.GetUserByUserid(hlt.UserID)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
					return
				}
				hltResps = append(hltResps, &HltResp{
					StudyHlt: &hlt,
					Name:     user.Name,
					Mobile:   user.Mobile,
				})
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS, &ActResp{
			StudyAct:  act,
			StudyHlts: hltResps,
		})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取活动列表
func GetStudyActs(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		share    = c.Query("share")
		status   = c.Query("status")   //0:未审核 1:审核通过(发布) 2:撤销发布
		dFlag    = c.Query("deadline") //Y:到期
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
	deadline := fmt.Sprintf("deadline>='2020-01-01'")
	if dFlag == "Y" {
		deadline = fmt.Sprintf(
			"deadline<='%s'", time.Now().Format("2006-01-02"))
	}
	if dFlag == "N" {
		deadline = fmt.Sprintf(
			"deadline>'%s'", time.Now().Format("2006-01-02"))
	}
	acts, err := models.GetStudyActs(share, status, deadline, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	var (
		data   []*ActResp
		userid = ""
	)
	if !strings.Contains(url, "v1") {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	if len(acts) > 0 {
		for _, act := range acts {
			if strings.Contains(url, "v1") {
				if act.TopicImage != "" {
					act.TopicImage = fsdjsrv.GetFsdjImageFullUrl(act.TopicImage)
				}
				urls := make([]string, 0)
				if act.ImageUrl != "" {
					for _, imageUrl := range strings.Split(act.ImageUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(imageUrl))
					}
					act.ImageUrls = urls
				}
				if strings.Contains(act.Content, "fsdj/dj_image") {
					if strings.Contains(act.Content, "api/fsdj/dj_image") {
						act.Content = strings.ReplaceAll(
							act.Content, "api/fsdj/dj_image", "fsdj/dj_image")
					}
				}
			} else {
				if act.TopicImage != "" {
					act.TopicImage = fsdjsrv.GetFsdjEappImageFullUrl(act.TopicImage)
				}
				urls := make([]string, 0)
				if act.ImageUrl != "" {
					for _, imageUrl := range strings.Split(act.ImageUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(imageUrl))
					}
					act.ImageUrls = urls
				}
				if strings.Contains(act.Content, "fsdj/dj_image") {
					if !strings.Contains(act.Content, "api/fsdj/dj_image") {
						act.Content = strings.ReplaceAll(
							act.Content, "fsdj/dj_image", "api/fsdj/dj_image")
					}
				}
			}
			act.JoinNum = len(act.StudyActdetails)
			joined := false //参与标志
			if len(userid) > 0 {
				if models.IsJoinStrudyAct(act.ID, userid) == "Y" {
					joined = true
				}
			}
			act.Joined = joined

			hltResps := make([]*HltResp, 0)
			hlts := act.StudyHlts
			if len(hlts) > 0 {
				for _, hltp := range hlts {
					hlt := hltp
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
					user, err := models.GetUserByUserid(hlt.UserID)
					if err != nil {
						appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
						return
					}
					hltResps = append(hltResps, &HltResp{
						StudyHlt: &hlt,
						Name:     user.Name,
						Mobile:   user.Mobile,
					})
				}
			}

			data = append(data, &ActResp{
				StudyAct:  act,
				StudyHlts: hltResps,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": data,
				"cnt":  models.GetStudyActsCnt(share, status, deadline),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取我的党建活动列表
func GetStudyActsMine(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		joinFlag = c.Query("joins")    //Y&N
		dFlag    = c.Query("deadline") //Y:到期
		url      = c.Request.URL.Path
		userid   string
	)
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
	deadline := fmt.Sprintf("deadline>='2020-01-01'")
	if dFlag == "Y" {
		deadline = fmt.Sprintf(
			"deadline<='%s'", time.Now().Format("2006-01-02"))
	}
	if dFlag == "N" {
		deadline = fmt.Sprintf(
			"deadline>'%s'", time.Now().Format("2006-01-02"))
	}
	acts, err := models.GetStudyActs(
		"", "1", deadline, 1, 10000)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(acts) > 0 {
		data := make([]*ActResp, 0)
		studyActs := make([]*models.StudyAct, 0)
		for _, act := range acts {
			if strings.Contains(url, "v1") {
				if act.TopicImage != "" {
					act.TopicImage = fsdjsrv.GetFsdjImageFullUrl(act.TopicImage)
				}
				urls := make([]string, 0)
				if act.ImageUrl != "" {
					for _, imageUrl := range strings.Split(act.ImageUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjImageFullUrl(imageUrl))
					}
					act.ImageUrls = urls
				}
				if strings.Contains(act.Content, "fsdj/dj_image") {
					if strings.Contains(act.Content, "api/fsdj/dj_image") {
						act.Content = strings.ReplaceAll(
							act.Content, "api/fsdj/dj_image", "fsdj/dj_image")
					}
				}
			} else {
				if act.TopicImage != "" {
					act.TopicImage = fsdjsrv.GetFsdjEappImageFullUrl(act.TopicImage)
				}
				urls := make([]string, 0)
				if act.ImageUrl != "" {
					for _, imageUrl := range strings.Split(act.ImageUrl, ";") {
						urls = append(urls, fsdjsrv.GetFsdjEappImageFullUrl(imageUrl))
					}
					act.ImageUrls = urls
				}
				if strings.Contains(act.Content, "fsdj/dj_image") {
					if !strings.Contains(act.Content, "api/fsdj/dj_image") {
						act.Content = strings.ReplaceAll(
							act.Content, "fsdj/dj_image", "api/fsdj/dj_image")
					}
				}
			}
			act.JoinNum = len(act.StudyActdetails)
			if joinFlag == models.IsJoinStrudyAct(act.ID, userid) {
				studyActs = append(studyActs, act)
			} else {
				continue
			}

			hltResps := make([]*HltResp, 0)
			hlts := act.StudyHlts
			if len(hlts) > 0 {
				for _, hltp := range hlts {
					hlt := hltp
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
					user, err := models.GetUserByUserid(hlt.UserID)
					if err != nil {
						appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
						return
					}
					hltResps = append(hltResps, &HltResp{
						StudyHlt: &hlt,
						Name:     user.Name,
						Mobile:   user.Mobile,
					})
				}
			}
			data = append(data, &ActResp{
				StudyAct:  act,
				StudyHlts: hltResps,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除党建活动
func DelStudyAct(c *gin.Context) {
	var appG = app.Gin{C: c}
	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS_VERIFY, "id can not be nil")
		return
	}
	if err := models.DelStudyAct(id); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type ActdetailForm struct {
	ID     uint   `json:"id"`
	ActID  uint   `json:"act_id"`
	UserID string `json:"userid"`
	Status string `json:"status"` //0:未审核 1:审核通过
}

//参与党建活动
func JoinStudyAct(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		id     = c.Param("id")
		status = "0"
		userid string
		sprq   string
	)
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
	actId, _ := strconv.Atoi(id)
	t := time.Now().Format("2006-01-02 15:04:05")
	if c.Query("status") != "" {
		status = c.Query("status")
		sprq = t
	}
	actdetail := &models.StudyActdetail{
		StudyActID: uint(actId),
		UserID:     userid,
		Status:     status,
		Bmrq:       t,
		Sprq:       sprq,
	}
	if err := models.AddStudyActdetail(actdetail); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type ActDetailResp struct {
	*models.StudyActdetail
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//获取党建活动参与信息审核列表
func GetApproveStudyActs(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		actId = c.Query("id")
		cond  = `study_act_id like '%'`
	)
	if len(actId) > 0 {
		cond = fmt.Sprintf("study_act_id='%s'", actId)
	}
	acts, err := models.GetApproveStudyActs(cond)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(acts) > 0 {
		var data = make([]*ActDetailResp, 0)
		for _, act := range acts {
			user, err := models.GetUserByUserid(act.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
				return
			}
			data = append(data, &ActDetailResp{
				StudyActdetail: act,
				Name:           user.Name,
				Mobile:         user.Mobile,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//党建活动参与信息审核
func ApproveStudyAct(c *gin.Context) {
	appG := app.Gin{C: c}
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	actdetail := &models.StudyActdetail{
		ID:     uint(id),
		Status: "1",
		Sprq:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.UpdStudyActdetail(actdetail); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//查询党建活动参与人员列表
func GetStudyActUsers(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Query("act_id")
	status := c.Query("status")
	actId, _ := strconv.Atoi(id)
	ads, err := models.GetStudyActdetails(uint(actId), status)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(ads) > 0 {
		var data = make([]*ActDetailResp, 0)
		for _, ad := range ads {
			user, err := models.GetUserByUserid(ad.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
				return
			}
			data = append(data, &ActDetailResp{
				StudyActdetail: ad,
				Name:           user.Name,
				Mobile:         user.Mobile,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//党建活动参与人员统计
func CountStudyAct(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		actId = c.Param("id")
	)
	groups, err := models.GetStudyGroupBySjdm("00")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	if len(groups) > 0 {
		data := make([]map[string]interface{}, 0)
		for _, group := range groups {
			resp, err := models.CountStudyAct(actId, group.Dm)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, err)
				return
			}
			if len(resp) > 0 {
				data = append(data, map[string]interface{}{
					"act_title": resp[0].Title,
					"group_mc":  resp[0].Mc,
					"list":      resp,
					"join_num":  len(resp),
				})
			}
		}
		if len(data) > 0 {
			appG.Response(http.StatusOK, e.SUCCESS, data)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
