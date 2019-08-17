package dingtalk

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/cron"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"log"
	"net/http"
	"time"
)

type LoginForm struct {
	AuthCode string `json:"auth_code"`
}

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	session := sessions.Default(c)
	var form LoginForm
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if form.AuthCode == "" {
		log.Println("no auth code")
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	id := dingtalk.GetUserId(form.AuthCode)
	if id != "" {
		userInfo := dingtalk.GetUserInfo(id)
		session.Set("userid", userInfo.UserID)
		if err := session.Save(); err != nil {
			log.Printf("session.Save() err:%v", err)
		}
		appG.Response(http.StatusOK, e.SUCCESS, userInfo)
		return
	}
	log.Println("user id is empty:in Login")
}
func JsApiConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	url := c.Query("url")
	if url == "" {
		log.Println("no url")
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	sign := dingtalk.GetJsApiConfig(url)
	if sign != "" {
		appG.Response(http.StatusOK, e.SUCCESS, sign)
		return
	}
	log.Println("url is empty:in JsApiConfig")
}

//同步一次部门用户信息
func DepartmentUserSync(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		wt      = 20 //发生网页劫持后，发送递归请求的次数
		syncNum = 30 //goroutine数量
	)
	go cron.DepartmentUserSync(wt, syncNum)
	appG.Response(http.StatusOK, e.SUCCESS, "同步请求发送成功")
}

//获取部门用户信息同步条数
func DepartmentUserSyncNum(c *gin.Context) {
	appG := app.Gin{C: c}
	t := time.Now().Format("2006-01-02") + " 00:00:00"
	depNum, deperr := models.CountDepartmentSyncNum(t)
	if deperr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_NUMBER_FAIL, nil)
		return
	}
	userNum, usererr := models.CountUserSyncNum(t)
	if usererr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_NUMBER_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["syncTime"] = time.Now().Format("2006-01-02")
	data["depNum"] = depNum
	data["userNum"] = userNum
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
