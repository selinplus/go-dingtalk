package ydks

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/ydksrv"
	"net/http"
	"strconv"
	"time"
)

type DataForm struct {
	Lb   string `json:"lb"`
	Data string `json:"data"`
}

func Workrecord(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DataForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	data := &models.Ydksworkrecord{
		Lb:   form.Lb,
		Req:  form.Data,
		Crrq: time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddWorkrecord(&data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetWorkrecordSend(c *gin.Context) {
	appG := app.Gin{C: c}
	rq := c.Query("rq")
	flag, err := strconv.Atoi(c.Query("flag"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	records, err := models.GetYtstworkrecords(rq, flag)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(records) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, records)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func Recv(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DataForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	data := &models.Ydksdata{
		Lb:   form.Lb,
		Data: form.Data,
		Rq:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddYdksdata(&data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GenDataFile(c *gin.Context) {
	appG := app.Gin{C: c}
	ydksrv.WriteIntoFile()
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
