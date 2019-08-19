package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
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
			leaf := models.IsParentDepartment(department.ID)
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
