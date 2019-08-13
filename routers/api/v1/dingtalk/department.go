package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
)

//获取部门列表
func GetDepartmentByParentID(c *gin.Context) {
	var appG = app.Gin{C: c}
	ParentID, _ := strconv.Atoi(c.Query("id"))
	departments, err := models.GetDepartmentByParentID(ParentID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
		return
	}
	if len(departments) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, departments)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
	}
}
