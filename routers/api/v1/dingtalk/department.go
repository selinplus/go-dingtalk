package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/cron"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//获取多部门详情时，排除outer属性
func GetDepartmentByIDs(c *gin.Context) {
	var (
		appG       = app.Gin{C: c}
		ids        = c.Query("ids")
		err        error
		department *models.Department
	)
	ss := strings.Split(ids, ",")
	for _, s := range ss {
		id, _ := strconv.Atoi(s)
		department, err = models.GetDepartmentByID(id)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
			return
		}
		if !department.OuterDept {
			appG.Response(http.StatusOK, e.SUCCESS, department)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, department)
}

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

//获取部门列表(不含外部部门)
func GetDepartmentByParentIDWithNoOuter(c *gin.Context) {
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
			if !department.OuterDept {
				leaf := models.IsLeafDepartment(department.ID)
				dt := map[string]interface{}{
					"key":    department.ID,
					"value":  department.ID,
					"title":  department.Name,
					"isLeaf": leaf,
				}
				dts = append(dts, dt)
			}
		}
		data["children"] = dts
		appG.Response(http.StatusOK, e.SUCCESS, data)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, data)
	}
}

//同步一次部门用户信息
func DepartmentUserSync(c *gin.Context) {
	var appG = app.Gin{C: c}
	go cron.DepartmentUserSync(20, 30)
	appG.Response(http.StatusOK, e.SUCCESS, "同步请求发送成功")
}

//获取当日部门用户信息同步条数
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

//获取企业员工人数
func OrgUserCount(c *gin.Context) {
	appG := app.Gin{C: c}
	cnt, err := dingtalk.OrgUserCount(20)
	if err != nil {
		appG.Response(http.StatusOK, e.SUCCESS, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, cnt)
}
