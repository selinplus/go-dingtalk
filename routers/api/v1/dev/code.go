package dev

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

//查询设备类型代码(树结构)
func GetDevtypeTree(c *gin.Context) {
	var appG = app.Gin{C: c}
	sjdm := c.Query("dm")
	parentDt, err := models.GetDevtypeByDm(sjdm)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	var dts []interface{}
	data := map[string]interface{}{
		"key":      parentDt.Dm,
		"value":    parentDt.Dm,
		"title":    parentDt.Mc,
		"children": dts,
	}
	devTypes, err := models.GetDevtypeBySjdm(sjdm)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	if len(devTypes) > 0 {
		for _, devType := range devTypes {
			leaf := models.IsLeafDevtype(devType.Dm)
			dt := map[string]interface{}{
				"key":    devType.Dm,
				"value":  devType.Dm,
				"title":  devType.Mc,
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

//查询操作属性代码
func GetDevProp(c *gin.Context) {
	appG := app.Gin{C: c}
	d, err := models.GetDevproperty()
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
		appG = app.Gin{C: c}
		data []map[string]string
		dm   = c.Query("dm")
		node = c.Query("node")
	)
	nodes, err := models.GetNextNode(dm, node)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if strings.Contains(nodes.Rname, "-1") {
		appG.Response(http.StatusOK, e.SUCCESS, "处理完成")
		return
	}
	cs := strings.Split(nodes.Rname, ",")
	for _, mobile := range cs {
		czrmp := map[string]string{}
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

//获取手工提报人员列表
func GetProcCustomizeList(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		data []map[string]string
	)
	xxzxUsers, err := models.GetUserByDepartmentID("70280083")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	for _, user := range xxzxUsers {
		czrmp := map[string]string{}
		czrmp["name"] = user.Name
		czrmp["mobile"] = user.Mobile
		data = append(data, czrmp)
	}
	users, err := models.GetUserByDepartmentID("29464263")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	for _, user := range users {
		czrmp := map[string]string{}
		czrmp["name"] = user.Name
		czrmp["mobile"] = user.Mobile
		data = append(data, czrmp)
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//查询提报类型代码
func GetProcType(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		pt   []*models.Proctype
		err  error
	)
	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]
	user, err := models.GetUserByUserid(userid)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
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
