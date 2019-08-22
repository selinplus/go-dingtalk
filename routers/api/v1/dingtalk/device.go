package dingtalk

import (
	"github.com/boombuler/barcode/qr"
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
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	sbbh := string(form.Lx) + "-" + timeStamp
	//生成二维码
	qrc := qrcode.NewQrCode(sbbh, 300, 300, qr.M, qr.Auto)
	name, _, err := qrc.Encode(qrcode.GetQrCodeFullPath())
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
