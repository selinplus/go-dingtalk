package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strings"
	"time"
)

type DevdeptForm struct {
	Jgdm   string `json:"jgdm"`
	Jgmc   string `json:"jgmc"`
	Sjjgdm string `json:"sjjgdm"`
	Gly    string `json:"gly"`
	Mobile string `json:"mobile"`
}

type GlyResp struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//增加设备管理机构
func AddDevdept(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   DevdeptForm
		userid string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if len(form.Mobile) > 0 {
		u, err := models.GetUserByMobile(form.Mobile)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_GET_USERBYMOBILE_FAIL, err)
			return
		}
		userid = u.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}
	jgdm, err := models.GenDevdeptDmBySjjgdm(form.Sjjgdm)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	devdept := models.Devdept{
		Jgdm:   jgdm,
		Jgmc:   form.Jgmc,
		Sjjgdm: form.Sjjgdm,
		Gly:    form.Gly,
		Lrr:    userid,
		Lrrq:   t,
		Xgr:    userid,
		Xgrq:   t,
	}
	if err := models.AddDevdept(&devdept); err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改设备管理机构
func UpdateDevdept(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   DevdeptForm
		userid string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if len(form.Mobile) > 0 {
		u, err := models.GetUserByMobile(form.Mobile)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_GET_USER_FAIL, err)
			return
		}
		userid = u.UserID
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
	devdept := models.Devdept{
		Jgdm:   form.Jgdm,
		Jgmc:   form.Jgmc,
		Sjjgdm: form.Sjjgdm,
		Gly:    form.Gly,
		Xgr:    userid,
		Xgrq:   t,
	}
	if err := models.UpdateDevdept(&devdept); err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备管理机构列表(树结构)
func GetDevdeptTree(c *gin.Context) {
	appG := app.Gin{C: c}
	tree, err := models.GetDevdeptTree()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	if len(tree) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, tree)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

//获取设备管理机构列表(循环遍历)
func GetDevdeptBySjjgdm(c *gin.Context) {
	var appG = app.Gin{C: c}
	jgdm := c.Query("jgdm")
	parentDt, err := models.GetDevdept(jgdm)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	var dts []interface{}
	data := map[string]interface{}{
		"key":      parentDt.Jgdm,
		"value":    parentDt.Jgdm,
		"title":    parentDt.Jgmc,
		"children": dts,
	}
	departments, err := models.GetDevdeptBySjjgdm(jgdm)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	if len(departments) > 0 {
		for _, department := range departments {
			leaf := models.IsLeafDevdept(department.Jgdm)
			dt := map[string]interface{}{
				"key":    department.Jgdm,
				"value":  department.Jgdm,
				"title":  department.Jgmc,
				"isLeaf": leaf,
			}
			dts = append(dts, dt)
		}
		data["children"] = dts
		appG.Response(http.StatusOK, e.SUCCESS, data)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, data)
	}
}

//删除设备管理机构
func DeleteDevdept(c *gin.Context) {
	appG := app.Gin{C: c}
	jgdm := c.Query("jgdm")
	if models.IsSjjg(jgdm) {
		appG.Response(http.StatusOK, e.ERROR_DELETE_DEVDETP_IS_PARENT, nil)
		return
	}
	if models.IsDevdeptUserExist(jgdm) {
		appG.Response(http.StatusOK, e.ERROR_DELETE_DEVDETP_NOT_NULL, nil)
		return
	}
	if err := models.DeleteDevdept(jgdm); err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取当前机构管理员信息
func GetDevdeptGly(c *gin.Context) {
	appG := app.Gin{C: c}
	jgdm := c.Query("jgdm")
	ddept, err := models.GetDevdept(jgdm)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_USER_FAIL, err)
		return
	}
	resps := make([]*GlyResp, 0)
	if ddept.Gly == "" {
		appG.Response(http.StatusOK, e.SUCCESS, resps)
		return
	}
	user, err := models.GetUserByUserid(ddept.Gly)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_USER_FAIL, err)
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
