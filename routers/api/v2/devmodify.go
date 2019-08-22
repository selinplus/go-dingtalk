package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"time"
)

type AddDevmodifyForm struct {
	DevID string `json:"devid"`
	Czlx  int    `json:"czlx"`
	Sydw  string `json:"sydw"`
	Syks  string `json:"syks"`
	Syr   string `json:"syr"`
	Cfwz  string `json:"cfwz"`
	Czrq  string `json:"czrq"`
	Czr   string `json:"czr"`
}

//流转
func AddDeviceMod(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddDevmodifyForm
		err  error
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	d := models.Devmodify{
		DevID: form.DevID,
		Czlx:  form.Czlx,
		Sydw:  form.Sydw,
		Syks:  form.Syks,
		Syr:   form.Syr,
		Cfwz:  form.Cfwz,
		Czrq:  form.Czrq,
		Czr:   form.Czr,
		Qsrq:  t,
	}
	err = models.ModifyZzrq(d.DevID, t)
	err = models.AddDevMod(d)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if d.ID > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, d.ID)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
	}
}
