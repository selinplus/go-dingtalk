package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
)

//新增盘点任务
func GetDevCkTask(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.Devcheck
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if err := models.AddDevCheckTask(&form); err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	if form.Ckself == "Y" {

	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取盘点任务列表
func GetDevCkTasks(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取盘点任务清册明细
func GetDevCkDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//显示其自我盘点过的历史记录
func GetDevCkHistory(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//设备盘点
func GetDevCheck(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
