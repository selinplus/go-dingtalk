package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"net/http"
	"strconv"
	"time"
)

type AddDeviceForm struct {
	ID   string
	Zcbh string `json:"zcbh"`
	Lx   int    `json:"lx"`
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
	Zt   int    `json:"zt"`
}

//单项录入
func AddDevice(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddDeviceForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
	sbbh := string(form.Lx) + timeStamp
	dev := models.Device{
		ID:   sbbh,
		Rkrq: t,
	}
	err := models.AddDevice(&dev)
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
