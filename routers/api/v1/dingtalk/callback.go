package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"net/http"
	"strings"
)

//注册事件回调
func RegisterCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	callbacks := []string{"user_add_org", "user_modify_org", "user_leave_org", "org_dept_create", "org_dept_modify", "org_dept_remove"}
	callbackURL := setting.DingtalkSetting.CallBackHost + "/api/v1/callback/detail"
	request := map[string]interface{}{
		"call_back_tag": callbacks,
		"token":         dingtalk.RandomString(8),
		"aes_key":       dingtalk.RandomString(43),
		"url":           callbackURL,
	}
	response, err := dingtalk.RegisterCallback(request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
	}
}

// 查询事件回调
func QueryCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	response, err := dingtalk.QueryCallback()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
	}
}

// 更新事件回调
func UpdateCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	call_back_tag := c.PostForm("callbacks")
	callbacks := strings.Split(call_back_tag, ",")
	//callbacks := []string{"user_add_org", "user_modify_org", "user_leave_org", "org_dept_create", "org_dept_modify", "org_dept_remove"}
	callbackURL := setting.DingtalkSetting.CallBackHost + "/callback/detail"
	request := map[string]interface{}{
		"call_back_tag": callbacks,
		"token":         dingtalk.RandomString(8),
		"aes_key":       dingtalk.RandomString(43),
		"url":           callbackURL,
	}
	response, err := dingtalk.UpdateCallback(request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
	}
}

// 删除事件回调
func DeleteCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	response, err := dingtalk.DeleteCallback()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
	}
}

// 获取回调失败的结果
func GetFailedCallbacks(c *gin.Context) {
	appG := app.Gin{C: c}
	response, err := dingtalk.GetFailedCallbacks()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
	}
}
