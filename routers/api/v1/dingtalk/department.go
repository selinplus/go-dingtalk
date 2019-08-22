package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/cron"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"time"
)

//获取部门详情
func GetDepartmentByID(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, errc := strconv.Atoi(c.Query("id"))
	if errc != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	department, errd := models.GetDepartmentByID(id)
	if errd != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, department)
}

//获取部门列表
func GetDepartmentByParentID(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, errc := strconv.Atoi(c.Query("id"))
	if errc != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	parentDt, errd := models.GetDepartmentByID(id)
	if errd != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	var dts []interface{}
	data := map[string]interface{}{
		"key":      parentDt.ID,
		"value":    parentDt.ID,
		"title":    parentDt.Name,
		"children": dts,
	}
	departments, err := models.GetDepartmentByParentID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	if len(departments) > 0 {
		for _, department := range departments {
			leaf := models.IsLeafDepartment(department.ID)
			dt := map[string]interface{}{
				"key":    department.ID,
				"value":  department.ID,
				"title":  department.Name,
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

//同步一次部门用户信息
func DepartmentUserSync(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		wt      = 20 //发生网页劫持后，发送递归请求的次数
		syncNum = 30 //goroutine数量
	)
	go cron.DepartmentUserSync(wt, syncNum)
	appG.Response(http.StatusOK, e.SUCCESS, "同步请求发送成功")
}

//获取部门用户信息同步条数
func DepartmentUserSyncNum(c *gin.Context) {
	appG := app.Gin{C: c}
	t := time.Now().Format("2006-01-02") + " 00:00:00"
	depNum, deperr := models.CountDepartmentSyncNum(t)
	if deperr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_NUMBER_FAIL, nil)
		return
	}
	userNum, usererr := models.CountUserSyncNum(t)
	if usererr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_NUMBER_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["syncTime"] = time.Now().Format("2006-01-02")
	data["depNum"] = depNum
	data["userNum"] = userNum
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
