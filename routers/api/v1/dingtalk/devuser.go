package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DevuserForm struct {
	ID   uint
	Jgdm string `json:"jgdm"`
	Syr  string `json:"syr"`
}

type DevuserResp struct {
	*models.Devuser
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//增加设备使用人员
func AddDevuser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DevuserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	for _, syr := range strings.Split(form.Syr, ",") {
		if !models.IsDevuserExist(form.Jgdm, syr) {
			devuser := models.Devuser{
				Jgdm: form.Jgdm,
				Syr:  syr,
			}
			if err := models.AddDevuser(&devuser); err != nil {
				appG.Response(http.StatusOK, e.ERROR_ADD_USER_FAIL, err)
				return
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改设备使用人员
func UpdateDevuser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DevuserForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	devuser := models.Devuser{
		ID:   form.ID,
		Jgdm: form.Jgdm,
		Syr:  form.Syr,
	}
	if err := models.UpdateDevuser(&devuser); err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备使用人员列表
func GetDevuserList(c *gin.Context) {
	appG := app.Gin{C: c}
	dus, err := models.GetDevuser(c.Query("jgdm"))
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_USER_FAIL, err)
		return
	}
	if len(dus) > 0 {
		resp := make([]*DevuserResp, 0)
		for _, du := range dus {
			user, err := models.GetUserByUserid(du.Syr)
			if err != nil {
				log.Println(err)
			}
			u := &DevuserResp{
				Devuser: du,
				Name:    user.Name,
				Mobile:  user.Mobile,
			}
			resp = append(resp, u)
		}
		appG.Response(http.StatusOK, e.SUCCESS, resp)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

//删除设备使用人员
func DeleteDevuser(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	if err := models.DeleteDevuser(uint(id)); err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_USER_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
