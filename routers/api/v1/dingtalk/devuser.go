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
				appG.Response(http.StatusOK, e.ERROR_ADD_USER_FAIL, err.Error())
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
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备使用人员列表
func GetDevuserList(c *gin.Context) {
	appG := app.Gin{C: c}
	dus, err := models.GetDevuser(c.Query("jgdm"))
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEVUSER_FAIL, err.Error())
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
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	if err := models.DeleteDevuser(uint(id)); err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_USER_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//eapp登录后获取人员信息
func DevLoginInfo(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		syrDepts = make([]map[string]string, 0)
		glyDepts = make([]map[string]string, 0)
		sfbz     = "0" //0:syr;1:gly;2:super
	)

	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]

	sDepts, err := models.GetSyrDepts(userid)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	for _, sDept := range sDepts {
		syrDept := map[string]string{
			"jgdm": sDept.Jgdm,
			"jgmc": sDept.Jgmc,
		}
		syrDepts = append(syrDepts, syrDept)
	}
	gDepts, err := models.GetGlyDepts(userid)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	for _, gDept := range gDepts {
		if gDept.Jgdm == "00" {
			sfbz = "2"
		}
		glyDept := map[string]string{
			"jgdm": gDept.Jgdm,
			"jgmc": gDept.Jgmc,
		}
		glyDepts = append(glyDepts, glyDept)
	}
	if len(gDepts) > 0 && sfbz != "2" {
		sfbz = "1"
	}
	data := map[string]interface{}{
		"sfbz":      sfbz,
		"syr_depts": syrDepts,
		"gly_depts": glyDepts,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
