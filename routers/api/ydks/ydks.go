package ydks

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/ydksrv"
	"net/http"
	"time"
)

type DataForm struct {
	Lb       string `json:"lb"`
	Data     string `json:"data"`
	Nsrsbh   string `json:"nsrsbh"`
	RecordID string `json:"record_id"`
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
		Lb:     form.Lb,
		Req:    form.Data,
		UserID: form.Nsrsbh,
		Crrq:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddWorkrecord(&data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data.ID)
}
func UpdWorkrecord(c *gin.Context) {
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
		Lb:       "updRecord",
		UserID:   form.Nsrsbh,
		RecordID: form.RecordID,
		Crrq:     time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddWorkrecord(&data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data.ID)
}
func GetWorkrecords(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		rq     = c.Query("rq")
		id     = c.Query("id")
		flag   = c.Query("flag") //""全部,"1"未推送,"2"已推送
		lbCond string            //发起待办||更新待办
		Cond   string            //id||rq=
	)
	if c.Query("lb") == "" {
		lbCond = "lb='updRecord'"
	} else {
		lbCond = "lb!='updRecord'"
	}
	if rq != "" {
		Cond = "crrq like '" + rq + "%'"
	}
	if id != "" {
		Cond = "id = " + id
	}
	records, err := models.GetYtstworkrecords(flag, Cond, lbCond)
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
	ydksrv.WriteIntoFile(c.Query("date"))
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
