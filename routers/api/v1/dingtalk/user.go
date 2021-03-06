package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"strings"
)

//获取部门用户列表
func GetUserByDepartmentID(c *gin.Context) {
	var appG = app.Gin{C: c}
	DeptID := c.Query("deptId")
	users, err := models.GetUserByDepartmentID(DeptID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	if len(users) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, users)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
	}
}

//获取用户部门详情（内网）
func GetDepartmentByUserMobile(c *gin.Context) {
	var appG = app.Gin{C: c}
	var dts []*models.Department
	mb := c.Query("mobile")
	user, errd := models.GetUserByMobile(mb)
	if errd != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	for _, deptId := range strings.Split(user.Department, ",") {
		deptId, err := strconv.Atoi(deptId)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
			return
		}
		dt, errd := models.GetDepartmentByID(deptId)
		if errd != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
			return
		}
		dts = append(dts, dt)
	}
	appG.Response(http.StatusOK, e.SUCCESS, dts)
}

//模糊查询用户
func GetUserByMc(c *gin.Context) {
	appG := app.Gin{C: c}
	mc := c.Query("mc")
	if len(mc) == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, "名称不能为空")
		return
	}
	users, err := models.GetUserByMc(mc)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
		return
	}
	if len(users) == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, users)
}
