package dev

import (
	"fmt"
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
				appG.Response(http.StatusOK, e.ERROR_ADD_USER_FAIL,
					fmt.Sprintf("增加设备使用人员错误：%v", err))
				return
			} else {
				user, err := models.GetUserByUserid(syr)
				if err != nil {
					appG.Response(http.StatusOK, e.ERROR_ADD_USER_FAIL,
						fmt.Sprintf("增加设备使用人员错误,[%s]获取user失败：%v", syr, err))
					return
				}
				if err := models.SaveUserdemo(&user); err != nil {
					appG.Response(http.StatusOK, e.ERROR_ADD_USER_FAIL,
						fmt.Sprintf("增加设备使用人员userdemo错误：%v", err))
					return
				}
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
		appG.Response(http.StatusOK, e.ERROR_GET_DEPT_USER_FAIL, err)
		return
	}
	if len(dus) > 0 {
		resp := make([]*DevuserResp, 0)
		for _, du := range dus {
			var name, mobile string
			user, err := models.GetUserdemoByUserid(du.Syr)
			if err != nil {
				log.Println(err)
				name = du.Syr
			} else {
				name, mobile = user.Name, user.Mobile
			}
			u := &DevuserResp{
				Devuser: du,
				Name:    name,
				Mobile:  mobile,
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
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if models.IsUserDevExist(c.Query("id")) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USERDEV_FAIL, nil)
		return
	}
	if err := models.DeleteDevuser(uint(id)); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USER_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//登录后获取人员身份信息
func LoginInfo(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		syrDepts = make([]map[string]string, 0)
		glyDepts = make([]map[string]string, 0)
		sfbz     = "3" //0:syr;1:gly;2:super;3:undefined
		userid   string
	)
	if c.Query("mobile") != "" {
		user, err := models.GetUserByMobile(c.Query("mobile"))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL,
				fmt.Sprintf("人员获取失败：%s", c.Query("mobile")))
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

	sDepts, err := models.GetSyrDepts(userid)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	if len(sDepts) > 0 {
		sfbz = "0"
		for _, sDept := range sDepts {
			syrDept := map[string]string{
				"jgdm": sDept.Jgdm,
				"jgmc": sDept.Jgmc,
			}
			syrDepts = append(syrDepts, syrDept)
		}
	}
	gDepts, err := models.GetGlyDepts(userid)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	if len(gDepts) > 0 {
		sfbz = "1"
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
	}
	data := map[string]interface{}{
		"sfbz":      sfbz,
		"syr_depts": syrDepts,
		"gly_depts": glyDepts,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
