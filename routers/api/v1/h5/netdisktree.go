package h5

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
)

type NetdiskTreeForm struct {
	ID     int    `json:"id"`
	PId    int    `json:"pId"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"` //inner useful
}

//新建网盘文件夹
func AddNetdiskDir(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    NetdiskTreeForm
		mobile  string
		userID  string
		err     error
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	mobile = form.Mobile
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	nd := models.NetdiskTree{
		PId:    form.PId,
		Name:   form.Name,
		UserID: userID,
	}
	if err = models.AddNetdiskDir(&nd); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DIRL_FAIL, nil)
		return
	}
	if nd.ID > 0 {
		if nd.ID == 1 { //treeid can not be 1
			_ = models.DeleteNetdiskDir(nd.UserID, 1)
			nd2 := models.NetdiskTree{
				PId:    form.PId,
				Name:   form.Name,
				UserID: userID,
			}
			if err = models.AddNetdiskDir(&nd2); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DIRL_FAIL, nil)
				return
			}
			appG.Response(http.StatusOK, e.SUCCESS, nd2.ID)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nd.ID)
	} else {
		appG.Response(http.StatusOK, e.ERROR_ADD_DIRL_FAIL, nil)
	}
}

//修改文件夹
func UpdateNetdiskDir(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    NetdiskTreeForm
		userID  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	id, mobile := form.ID, form.Mobile
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	dir := models.NetdiskTree{
		ID:     id,
		PId:    form.PId,
		Name:   form.Name,
		UserID: userID,
	}
	err := models.UpdateNetdiskDir(&dir)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DIR_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除文件夹
func DeleteNetdiskDir(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	id, _ := strconv.Atoi(c.Query("id"))
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	if models.IsParentDir(userID, id) {
		appG.Response(http.StatusOK, e.ERROR_DELETE_DIR_IS_PARENT, nil)
		return
	}
	if models.IsDirContainFile(userID, id) {
		appG.Response(http.StatusOK, e.ERROR_DELETE_DIR_HAS_FILE, nil)
		return
	}
	err := models.DeleteNetdiskDir(userID, id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DIR_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取用户网盘文件夹列表
func GetNetdiskDirTree(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
		err     error
	)
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	tree, err := models.GetNetdiskTree(userID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DIR_LIST_FAIL, nil)
		return
	}
	if len(tree) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, tree)
	} else {
		appG.Response(http.StatusOK, e.ERROR, nil)
	}
}
