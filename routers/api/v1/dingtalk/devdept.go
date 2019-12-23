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
	ID     string
	Jgmc   string `json:"jgmc"`
	Sjjgdm string `json:"sjjgdm"`
	Gly    string `json:"gly"`
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
	jgdm, err := models.GenDevdeptDmBySjdm(form.Sjjgdm)
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
	if err := models.AddDevdept(devdept); err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改设备管理机构
func GetDevdeptByID(c *gin.Context) {
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
		Jgmc:   form.Jgmc,
		Sjjgdm: form.Sjjgdm,
		Gly:    form.Gly,
		Xgr:    userid,
		Xgrq:   t,
	}
	if err := models.UpdateDevdept(&devdept); err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备管理机构列表
func GetDevdeptTree(c *gin.Context) {
	appG := app.Gin{C: c}
	tree, err := models.GetDevdeptTree()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(tree) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, tree)
	} else {
		appG.Response(http.StatusOK, e.ERROR, nil)
	}
}
