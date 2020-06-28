package fsdj

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strings"
	"time"
)

type GroupForm struct {
	Dm     string `json:"dm"`
	Mc     string `json:"mc"`
	Sjdm   string `json:"sjdm"`
	Gly    string `json:"gly"`
	Mobile string `json:"mobile"`
}

type GlyResp struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//增加学习小组
func AddGroup(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   GroupForm
		userid string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if len(form.Mobile) > 0 {
		user, err := models.GetUserByMobile(form.Mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, err)
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

	dm, err := models.GenGroupDmBySjjgdm(form.Sjdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	group := models.StudyGroup{
		Dm:   dm,
		Mc:   form.Mc,
		Sjdm: form.Sjdm,
		Gly:  form.Gly,
		Lrr:  userid,
		Lrrq: t,
		Xgr:  userid,
		Xgrq: t,
	}
	if err := models.AddStudyGroup(&group); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func UpdGroup(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   GroupForm
		userid string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if len(form.Mobile) > 0 {
		user, err := models.GetUserByMobile(form.Mobile)
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

	t := time.Now().Format("2006-01-02 15:04:05")
	group := models.StudyGroup{
		Dm:   form.Dm,
		Xgr:  userid,
		Xgrq: t,
	}
	url := c.Request.URL.Path
	if strings.Contains(url, "group/upd") { //修改学习小组
		group.Mc = form.Mc
		group.Sjdm = form.Sjdm
	}
	if strings.Contains(url, "gly/add") { //设置学习小组管理员
		group.Gly = form.Gly
	}
	if err := models.UpdStudyGroup(&group); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取学习小组信息
func GetGroup(c *gin.Context) {
	appG := app.Gin{C: c}
	d, err := models.GetStudyGroup(c.Query("dm"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	if d != nil {
		if d.Gly != "" {
			u, _ := models.GetUserByUserid(d.Gly)
			d.Gly = u.Name
		}
		if d.Lrr != "" {
			u, _ := models.GetUserByUserid(d.Lrr)
			d.Lrr = u.Name
		}
		if d.Xgr != "" {
			gly, _ := models.GetUserByUserid(d.Xgr)
			d.Xgr = gly.Name
		}
		appG.Response(http.StatusOK, e.SUCCESS, d)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取学习小组列表(树结构)
func GetGroupTree(c *gin.Context) {
	appG := app.Gin{C: c}
	dm := c.Query("dm")
	jgdms := strings.Split(dm, ",")
	dms := make([]string, 0)
	for _, dmSrc := range jgdms {
		flag := true
		for _, dm := range jgdms {
			if len(dmSrc) > len(dm) {
				if strings.Contains(dmSrc, dm) {
					flag = false
					break
				}
			}
		}
		if flag {
			dms = append(dms, dmSrc)
		}
	}
	data := make([]models.StudyGroupTree, 0)
	for _, dm := range dms {
		tree, err := models.GetStudyGroupTree(dm)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
			return
		}
		if len(tree) > 0 {
			data = append(data, tree...)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//删除学习小组
func DelGroup(c *gin.Context) {
	appG := app.Gin{C: c}
	dm := c.Query("dm")
	if models.IsStudySjjg(dm) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DETP_IS_PARENT, nil)
		return
	}
	if models.IsNullStudyGroup(dm) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DETP_NOT_NULL, nil)
		return
	}
	if err := models.DelStudyGroup(dm); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除当前学习小组管理员
func DelGroupGly(c *gin.Context) {
	appG := app.Gin{C: c}
	dm := c.Query("dm")
	group := map[string]interface{}{
		"dm":   dm,
		"gly":  "",
		"xgrq": time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.DelStudyGroupGly(group); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取当前学习小组管理员信息
func GetGroupGly(c *gin.Context) {
	appG := app.Gin{C: c}
	dm := c.Query("dm")
	ddept, err := models.GetStudyGroup(dm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
		return
	}
	resps := make([]*GlyResp, 0)
	if ddept.Gly == "" {
		appG.Response(http.StatusOK, e.SUCCESS, resps)
		return
	}
	user, err := models.GetUserByUserid(ddept.Gly)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
		return
	}
	resp := &GlyResp{
		UserID: user.UserID,
		Name:   user.Name,
		Mobile: user.Mobile,
	}
	resps = append(resps, resp)
	appG.Response(http.StatusOK, e.SUCCESS, resps)
}
