package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"net/http"
	"strings"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserForm struct {
	UserAccount string `json:"user_account"`
	Username    string `json:"username"`
	DepartID    string `json:"depart_id"`
	DepartName  string `json:"depart_name"`
	UserRole    string `json:"user_role"`
}

// 用户登录
func Login(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LoginForm
	)
	resp := map[string]interface{}{}
	result := map[string]interface{}{}
	userInfo := map[string]interface{}{}
	depart := map[string]interface{}{}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	//todo: inner deployment
	// 引入"crypto/tls":解决golang https请求提示x509: certificate signed by unknown authority
	/*ts := &tls.Config{InsecureSkipVerify: true}
	pMap := map[string]string{
		"username": form.Username,
		"password": form.Password,
	}
	_, body, errs := gorequest.New().TLSClientConfig(ts).
		Post(setting.AppSetting.LoginUrl + "/jeecg-boot/sys/login").
		Type(gorequest.TypeJSON).SendMap(pMap).End()
	if len(errs) > 0 {
		data := fmt.Sprintf("login err:%v", errs[0])
		appG.Response(http.StatusOK, e.ERROR, data)
		return
	} else {
		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			data := fmt.Sprintf("unmarshall body error:%v", err)
			appG.Response(http.StatusOK, e.ERROR, data)
			return
		}
		if resp["result"] != nil {
			result = resp["result"].(map[string]interface{})
			userInfo = result["userInfo"].(map[string]interface{})
			departs := result["departs"].([]interface{})
			depart = departs[0].(map[string]interface{})
		}
	}*/

	//todo: internet test
	token, _ := util.GenerateToken(form.Username, form.Password)
	if form.Username == "test" {
		resp["success"] = "True"
		resp["message"] = "登录成功"
		result["token"] = token
		userInfo["id"] = "13706002531"
		userInfo["username"] = "张三"
		userInfo["userAccount"] = "test"
		depart["id"] = "13706130900"
		depart["departName"] = "XX市XX区信息中心"
		depart["parentId"] = "13706130000"
	}
	if form.Username == "test1" {
		resp["success"] = "True"
		resp["message"] = "登录成功"
		result["token"] = token
		userInfo["id"] = "test1"
		userInfo["username"] = "王五"
		userInfo["userAccount"] = "test1"
		depart["id"] = "13706001800"
		depart["departName"] = "XX市信息中心"
		depart["parentId"] = "13706000000"
	}
	if form.Username == "admin" {
		resp["success"] = "True"
		resp["message"] = "登录成功"
		result["token"] = token
		userInfo["id"] = "fsdj_admin"
		userInfo["username"] = "超级管理员"
		userInfo["userAccount"] = "admin"
		depart["id"] = "13706001800"
		depart["departName"] = "福山品牌党建"
		depart["parentId"] = "13706000000"
	}

	data := map[string]interface{}{
		"success":     resp["success"],
		"message":     resp["message"],
		"token":       result["token"],
		"userid":      userInfo["id"],
		"username":    userInfo["username"],
		"userAccount": userInfo["userAccount"],
		"departID":    depart["id"],
		"departName":  depart["departName"],
		"parentId":    depart["parentId"],
	}

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 获取用户信息,存入session
func UserInfo(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		session = sessions.Default(c)
		form    UserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	session.Set("userAccount", form.UserAccount)
	session.Set("username", form.Username)
	session.Set("departID", form.DepartID)
	session.Set("departName", form.DepartName)
	session.Set("userRole", form.UserRole)
	if err := session.Save(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 代理转发智税平台用户登录
func Rlogin(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form LoginForm
	)
	resp := map[string]interface{}{}
	result := map[string]interface{}{}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// 引入"crypto/tls":解决golang https请求提示x509: certificate signed by unknown authority
	ts := &tls.Config{InsecureSkipVerify: true}
	pMap := map[string]string{
		"username": form.Username,
		"password": form.Password,
	}
	_, body, errs := gorequest.New().TLSClientConfig(ts).
		Post(setting.AppSetting.LoginUrl + "/jeecg-boot/sys/login").
		Type(gorequest.TypeJSON).SendMap(pMap).End()
	if len(errs) > 0 {
		data := fmt.Sprintf("login err:%v", errs[0])
		appG.Response(http.StatusOK, e.ERROR, data)
		return
	} else {
		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			data := fmt.Sprintf("unmarshall body error:%v", err)
			appG.Response(http.StatusOK, e.ERROR, data)
			return
		}
		if resp["result"] != nil {
			result = resp["result"].(map[string]interface{})
			appG.Response(http.StatusOK, e.SUCCESS, result)
			return
		}
		appG.Response(http.StatusOK, e.ERROR, nil)
	}
}

//获取数据中台路由表
func GetRoutes(c *gin.Context) {
	var appG = app.Gin{C: c}
	var IeHeader = `Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Win64; x64; Trident/4.0; .NET CLR 2.0.limit727; SLCC2; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET4.0C; .NET4.0E)`

	token := c.Query("token")
	// 引入"crypto/tls":解决golang https请求提示x509: certificate signed by unknown authority
	ts := &tls.Config{InsecureSkipVerify: true}
	url := setting.AppSetting.LoginUrl + "/jeecg-boot/sys/permission/getUserPermissionByToken" +
		"?token=" + token
	_, body, errs := gorequest.New().TLSClientConfig(ts).Get(url).
		Set("Content-Type", "application/json").
		Set("X-Access-Token", token).Set("Accept-Language", "zh-CN").
		Set("Accept", "*/*").Set("User-Agent", IeHeader).End()
	if len(errs) > 0 {
		data := fmt.Sprintf("Get routes err:%v", errs[0])
		appG.Response(http.StatusOK, e.ERROR, data)
		return
	} else {
		if !strings.Contains(body, "查询成功") {
			appG.Response(http.StatusOK, e.ERROR, body)
			return
		}
		resp := make(map[string]interface{})
		result := map[string]interface{}{}
		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			data := fmt.Sprintf("unmarshall body error:%v", err)
			appG.Response(http.StatusOK, e.ERROR, data)
			return
		}
		if resp["result"] != nil {
			result = resp["result"].(map[string]interface{})
			menu := make(map[string]interface{})
			for _, m := range result["menu"].([]interface{}) {
				mMap := m.(map[string]interface{})
				if mMap["path"].(string) == "/shuiwuyewu" {
					menu = map[string]interface{}{
						"redirect":  mMap["redirect"],
						"path":      mMap["path"],
						"component": mMap["component"],
						"route":     mMap["route"],
					}
					for _, c := range mMap["children"].([]interface{}) {
						cMap := c.(map[string]interface{})
						if cMap["path"].(string) == "/zxfxsm" {
							menu["children"] = []map[string]interface{}{
								{
									"path":      cMap["path"],
									"component": cMap["component"],
									"route":     cMap["route"],
									"children":  cMap["children"],
								},
							}
						}
					}
				}
			}
			appG.Response(http.StatusOK, e.SUCCESS,
				map[string]interface{}{
					"allAuth": result["allAuth"],
					"auth":    result["auth"],
					"menu": []map[string]interface{}{
						menu,
					},
				})
			return
		}
		appG.Response(http.StatusOK, e.ERROR, nil)
	}
}
