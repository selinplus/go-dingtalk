package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
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
