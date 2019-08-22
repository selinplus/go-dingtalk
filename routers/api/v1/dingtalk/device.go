package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/qrcode"
	"github.com/selinplus/go-dingtalk/pkg/upload"
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
	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	sbbh := string(form.Lx) + "_" + timeStamp + "_" + form.Zcbh
	//生成二维码
	//qrc := qrcode.NewQrCode(sbbh, 300, 300, qr.M, qr.Auto)
	//name, _, err := qrc.Encode(qrcode.GetQrCodeFullPath())
	name, _, err := qrcode.GenerateQrWithLogo(sbbh, qrcode.GetQrCodeFullPath())
	if err != nil {
		log.Println(err)
	}
	dev := models.Device{
		ID:    sbbh,
		Zcbh:  form.Zcbh,
		Lx:    form.Lx,
		Mc:    form.Mc,
		Xh:    form.Xh,
		Xlh:   form.Xlh,
		Ly:    form.Ly,
		Scs:   form.Scs,
		Scrq:  form.Scrq,
		Grrq:  form.Grrq,
		Bfnx:  form.Bfnx,
		Jg:    form.Jg,
		Zp:    form.Zp,
		Gys:   form.Gys,
		Rkrq:  form.Rkrq,
		QrUrl: qrcode.GetQrCodeFullUrl(name),
		Czr:   form.Czr,
		Zt:    form.Zt,
	}
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
	file, image, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if image == nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	src := fullPath + imageName
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_FILE_FORMAT, nil)
		return
	}
	if err = upload.CheckImage(fullPath); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_FILE_FAIL, nil)
		return
	}
	if err = c.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_FILE_FAIL, nil)
		return
	}
	if errDev := models.ImpDevices(imageName); errDev != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, errDev)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备列表
func GetDevices(c *gin.Context) {
	appG := app.Gin{C: c}
	mc := c.Query("mc")
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	devs, err := models.GetDevices(mc, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	total, er := models.GetDevicesCount(mc)
	if er != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = total
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
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEV_FAIL, nil)
		return
	}
	if len(dev.ID) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, dev)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEV_FAIL, nil)
	}
}

//生成二维码
func QrCode(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		qrs  []interface{}
	)
	ids := c.Query("id")
	for _, id := range strings.Split(ids, ",") {
		name, _, err := qrcode.GenerateQrWithLogo(id, qrcode.GetQrCodeFullPath())
		if err != nil {
			log.Println(err)
		}
		data := map[string]string{
			"qrName": name,
			"qrUrl":  qrcode.GetQrCodeFullUrl(name),
		}
		qrs = append(qrs, &data)
	}
	appG.Response(http.StatusOK, e.SUCCESS, qrs)
}
