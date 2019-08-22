package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
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
