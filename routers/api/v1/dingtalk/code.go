package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strings"
)

//查询设备状态代码
func GetDevstate(c *gin.Context) {
	appG := app.Gin{C: c}
	d, err := models.GetDevstate()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(d) == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
}

//查询设备类型代码
func GetDevtype(c *gin.Context) {
	appG := app.Gin{C: c}
	d, err := models.GetDevtype()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(d) == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
}

//查询操作类型代码
func GetDevOp(c *gin.Context) {
	appG := app.Gin{C: c}
	d, err := models.GetDevOp()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(d) == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
}

//获取下一节点操作人
func GetProcCzr(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		data  []map[string]string
		czrmp = map[string]string{}
		dm    = c.Query("dm")
		node  = c.Query("node")
	)
	nodes, err := models.GetNextNode(dm, node)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	cs := strings.Split(nodes.Rname, ",")
	for _, mobile := range cs {
		czr, uerr := models.GetUserByMobile(mobile)
		if uerr != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		czrmp["name"] = czr.Name
		czrmp["mobile"] = mobile
		data = append(data, czrmp)
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//查询提报类型代码
func GetProcType(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		mobile = c.Query("mobile")
		pt     []*models.Proctype
		err    error
	)
	user, err := models.GetUserByMobile(mobile)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if strings.Contains(user.Department, "70280083") {
		pt, err = models.GetProctypeAll()
	} else {
		pt, err = models.GetProctype()
	}
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(pt) == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, pt)
}
