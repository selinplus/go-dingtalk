package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"time"
)

type DevmodifyForm struct {
	DevID string `json:"devid"`
	Czlx  string `json:"czlx"`
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
		form DevmodifyForm
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
	flag, err := models.IsLastModifyZzrqExist(d.DevID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if flag {
		err = models.ModifyZzrq(d.DevID, t)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
	}
	err = models.AddDevMod(&d)
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

//设备流水记录查询
func GetDevModList(c *gin.Context) {
	appG := app.Gin{C: c}
	devid := c.Query("devid")
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	devs, err := models.GetDevMods(devid, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	total, er := models.GetDevModCount(devid)
	if er != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
