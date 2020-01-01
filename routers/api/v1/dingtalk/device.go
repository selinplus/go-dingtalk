package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/qrcode"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DeviceForm struct {
	ID   string
	Zcbh string `json:"zcbh"`
	Lx   string `json:"lx"`
	Mc   string `json:"mc" `
	Xh   string `json:"xh"`
	Xlh  string `json:"xlh"`
	Ly   string `json:"ly"`
	Scs  string `json:"scs"`
	Scrq string `json:"scrq"`
	Grrq string `json:"grrq"`
	Bfnx string `json:"bfnx"`
	Jg   string `json:"jg"`
	Zp   string `json:"zp"`
	Gys  string `json:"gys"`
	Rkrq string `json:"rkrq"`
	Czr  string `json:"czr"`
	Zt   string `json:"zt"`
}

//单项录入
func AddDevice(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DeviceForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if _, err := models.GetUserByMobile(form.Czr); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
		return
	}
	sbbh := models.GenerateSbbh(form.Lx, form.Xlh)
	dev := models.Device{
		ID:   sbbh,
		Zcbh: form.Zcbh,
		Lx:   form.Lx,
		Mc:   form.Mc,
		Xh:   form.Xh,
		Xlh:  form.Xlh,
		Ly:   form.Ly,
		Scs:  form.Scs,
		Scrq: form.Scrq,
		Grrq: form.Grrq,
		Bfnx: form.Bfnx,
		Jg:   form.Jg,
		Zp:   form.Zp,
		Gys:  form.Gys,
		Rkrq: time.Now().Format("2006-01-02 15:04:05"),
		Czr:  form.Czr,
		Zt:   form.Zt,
	}
	if models.IsXlhExist(form.Xlh) {
		appG.Response(http.StatusInternalServerError, e.ERROR_XLHEXIST_FAIL, nil)
		return
	}
	//生成二维码
	name, _, err := qrcode.GenerateQrWithLogo(sbbh, qrcode.GetQrCodeFullPath())
	if err != nil {
		log.Println(err)
	}
	dev.QrUrl = qrcode.GetQrCodeFullUrl(name)
	err = models.AddDevice(&dev)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
		return
	}
	if len(dev.ID) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
	}
}

//批量导入
func ImpDevices(c *gin.Context) {
	appG := app.Gin{C: c}
	czr := c.Query("czr")
	if _, err := models.GetUserByMobile(czr); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	errDev, success, failed := models.ImpDevices(file, czr)
	data := map[string]interface{}{
		"suNum":  success,
		"faNum":  failed,
		"errDev": errDev,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//更新设备信息
func UpdateDevice(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DeviceForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	dev := models.Device{
		ID:   form.ID,
		Zcbh: form.Zcbh,
		Lx:   form.Lx,
		Mc:   form.Mc,
		Xh:   form.Xh,
		Xlh:  form.Xlh,
		Ly:   form.Ly,
		Scs:  form.Scs,
		Scrq: form.Scrq,
		Grrq: form.Grrq,
		Bfnx: form.Bfnx,
		Jg:   form.Jg,
		Zp:   form.Zp,
		Gys:  form.Gys,
		Czr:  form.Czr,
		Zt:   form.Zt,
	}
	if form.Czr != "" {
		if _, err := models.GetUserByMobile(form.Czr); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
	}
	err := models.EditDevice(&dev)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DEV_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备列表
func GetDevices(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		rkrqq = c.Query("rkrqq")
		rkrqz = c.Query("rkrqz")
		sbbh  = c.Query("sbbh")
		xlh   = c.Query("xlh")
		syr   = c.Query("syr")
		mc    = c.Query("mc")
	)
	if rkrqq == "" {
		rkrqq = "2000-01-01 00:00:00"
	}
	if rkrqz == "" {
		rkrqz = "2099-01-01 00:00:00"
	}
	con := map[string]string{
		"rkrqq": rkrqq,
		"rkrqz": rkrqz,
		"sbbh":  sbbh,
		"xlh":   xlh,
		"syr":   syr,
		"mc":    mc,
	}
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	devs, err := models.GetDevices(con, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	for _, dev := range devs {
		if dev.Czr != "" {
			uczr, _ := models.GetUserByMobile(dev.Czr)
			dev.Czr = uczr.Name
		}
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = len(devs)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取设备详情
func GetDeviceByID(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Query("id")
	dev, err := models.GetDeviceByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEV_FAIL, nil)
		return
	}
	if len(dev.ID) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, dev)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEV_FAIL, nil)
	}
}

//查询设备信息及当前使用状态详情
func GetDeviceModByDevID(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Query("id")
	dev, err := models.GetDeviceModByDevID(id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEV_FAIL, nil)
		return
	}
	if len(dev.ID) > 0 {
		user, _ := models.GetUserByMobile(dev.Czr)
		dev.Czr = user.Name
		appG.Response(http.StatusOK, e.SUCCESS, dev)
	} else {
		appG.Response(http.StatusOK, e.ERROR_GET_DEV_FAIL, nil)
	}
}

//获取当前用户设备列表
func GetDevicesByUser(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		rkrqq    = c.Query("rkrqq")
		rkrqz    = c.Query("rkrqz")
		sbbh     = c.Query("sbbh")
		xlh      = c.Query("xlh")
		mc       = c.Query("mc")
		pageNo   int
		pageSize int
	)
	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]

	if rkrqq == "" {
		rkrqq = "2000-01-01 00:00:00"
	}
	if rkrqz == "" {
		rkrqz = "2099-01-01 00:00:00"
	}
	syr, uerr := models.GetUserByUserid(userid)
	if uerr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
		return
	}
	con := map[string]string{
		"rkrqq": rkrqq,
		"rkrqz": rkrqz,
		"sbbh":  sbbh,
		"xlh":   xlh,
		"syr":   syr.Mobile,
		"mc":    mc,
	}
	if c.Query("pageNo") == "" {
		pageNo = 1
	} else {
		pageNo, _ = strconv.Atoi(c.Query("pageNo"))
	}
	if c.Query("pageSize") == "" {
		pageSize = 10000
	} else {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	devs, err := models.GetDevices(con, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	for _, dev := range devs {
		if dev.Czr != "" {
			uczr, _ := models.GetUserByMobile(dev.Czr)
			dev.Czr = uczr.Name
		}
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = len(devs)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
