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

type DevinfoForm struct {
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
	Gys  string `json:"gys"`
	Czr  string `json:"czr"`
	Zt   string `json:"zt"`
	Jgdm string `json:"jgdm"`
	Syr  string `json:"syr"`
	Sx   string `json:"sx"`
}

type DevOpForm struct {
	Ids     []string `json:"ids"`
	SrcJgdm string   `json:"src_jgdm"`
	DstJgdm string   `json:"dst_jgdm"`
	Czr     string   `json:"czr"`     //inner传递操作人mobile
	Syr     string   `json:"syr"`     //inner传递使用人mobile
	CuserID string   `json:"cuserid"` //epp传递操作人userid
	SuserID string   `json:"suserid"` //epp传递使用人userid
	Cfwz    string   `json:"cfwz"`
	Czlx    string   `json:"czlx"`
}

//单项录入
func AddDevinfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DevinfoForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	czr, err := models.GetUserByMobile(form.Czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, err.Error())
		return
	}
	sbbh := models.GenerateSbbh(form.Lx, form.Xlh)
	t := time.Now().Format("2006-01-02 15:04:05")
	dev := &models.Devinfo{
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
		Gys:  form.Gys,
		Rkrq: t,
		Czrq: t,
		Czr:  czr.UserID,
		Zt:   "1",
		Jgdm: "00",
		Sx:   "1",
	}
	if models.IsDevXlhExist(form.Xlh) {
		appG.Response(http.StatusInternalServerError, e.ERROR_XLHEXIST_FAIL, nil)
		return
	}
	//生成二维码
	name, _, err := qrcode.GenerateQrWithLogo(sbbh, qrcode.GetQrCodeFullPath())
	if err != nil {
		log.Println(err)
	}
	dev.QrUrl = qrcode.GetQrCodeFullUrl(name)
	err = models.AddDevinfo(dev)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//批量导入
func ImpDevinfos(c *gin.Context) {
	appG := app.Gin{C: c}
	czr := c.Query("czr")
	user, err := models.GetUserByMobile(czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, err.Error())
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	errDev, success, failed, err := models.ImpDevinfos(file, user.UserID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, err.Error())
		return
	}
	if success == 0 && failed == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
		return
	}
	data := map[string]interface{}{
		"suNum":  success,
		"faNum":  failed,
		"errDev": errDev,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//更新设备信息
func UpdateDevinfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DevinfoForm
		syr  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	czr, err := models.GetUserByMobile(form.Czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, err.Error())
		return
	}
	if form.Syr != "" {
		suser, err := models.GetUserByMobile(form.Syr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, err.Error())
			return
		}
		syr = suser.UserID
	}
	dev := &models.Devinfo{
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
		Gys:  form.Gys,
		Czr:  czr.UserID,
		Czrq: time.Now().Format("2006-01-02 15:04:05"),
		Zt:   form.Zt,
		Jgdm: form.Jgdm,
		Syr:  syr,
		Sx:   form.Sx,
	}
	err = models.EditDevinfo(dev)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DEV_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type DevResponse struct {
	*models.Devinfo
	Jgmc string `json:"jgmc"`
}

//获取设备列表(inner多条件查询设备)
func GetDevinfos(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		rkrqq    = c.Query("rkrqq")
		rkrqz    = c.Query("rkrqz")
		sbbh     = c.Query("sbbh")
		xlh      = c.Query("xlh")
		syr      = c.Query("syr")
		mc       = c.Query("mc")
		jgdm     = c.Query("jgdm")
		pageNo   int
		pageSize int
	)
	if rkrqq == "" {
		rkrqq = "2000-01-01 00:00:00"
	}
	if rkrqz == "" {
		rkrqz = "2099-01-01 00:00:00"
	}
	if syr != "" {
		user, err := models.GetUserByMobile(syr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		syr = user.UserID
	}
	con := map[string]string{
		"rkrqq": rkrqq,
		"rkrqz": rkrqz,
		"sbbh":  sbbh,
		"xlh":   xlh,
		"syr":   syr,
		"mc":    mc,
		"jgdm":  jgdm,
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
	devs, err := models.GetDevinfos(con, pageNo, pageSize, "0")
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	resps := make([]*DevResponse, 0)
	for _, dev := range devs {
		if dev.Syr != "" {
			usyr, _ := models.GetUserByUserid(dev.Syr)
			dev.Syr = usyr.Name
		}
		dept, err := models.GetDevdept(dev.Jgdm)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
			return
		}
		resp := &DevResponse{
			Devinfo: dev,
			Jgmc:    dept.Jgmc,
		}
		resps = append(resps, resp)
	}
	data := make(map[string]interface{})
	data["lists"] = resps
	data["total"] = len(resps)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取设备详情
func GetDevinfoByID(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Query("id")
	dev, err := models.GetDevinfoByID(id)
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

// 获取设备列表(管理员查询||eapp使用人查询)
// jgdm!="",bz=0:管理人名下在库设备;jgdm!="",bz=3:管理人名下共用设备;
// jgdm!="",bz=4:管理人名下已分配设备;jgdm!="",bz=6:管理人名下已借出设备;
// eapp :jgdm!="",bz=10:管理人名下已分配&已借出设备;jgdm=="",bz=10:使用人;
func GetDevinfosByUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		jgdm = c.Query("jgdm")
		bz   = c.Query("bz")
	)
	//使用人查看名下设备
	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]

	con := map[string]string{
		"rkrqq": "2000-01-01 00:00:00",
		"rkrqz": "2099-01-01 00:00:00",
		"syr":   "",
		"jgdm":  "",
	}
	resps := make([]*models.Devinfo, 0)
	if jgdm != "" {
		for _, dm := range strings.Split(jgdm, ",") {
			con["jgdm"] = dm
			devs, err := models.GetDevinfos(con, 1, 10000, bz)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
				return
			}
			resps = append(resps, devs...)
		}
	} else {
		con["syr"] = userid
		devs, err := models.GetDevinfos(con, 1, 10000, bz)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
			return
		}
		resps = append(resps, devs...)
	}
	data := make(map[string]interface{})
	data["lists"] = resps
	data["total"] = len(resps)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//设备下发
func DevIssued(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DevOpForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	czr, err := models.GetUserByMobile(form.Czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
		return
	}
	if err := models.DevIssued(form.Ids, form.SrcJgdm, form.DstJgdm, czr.UserID); err != nil {
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//设备分配&借出&收回&交回
func DevAllocate(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DevOpForm
		czr  string
		syr  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if form.Czr != "" {
		cuser, err := models.GetUserByMobile(form.Czr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		czr = cuser.UserID
	}
	if form.CuserID != "" {
		czr = form.CuserID
	}
	if form.Syr != "" {
		if form.Syr != " " {
			suser, err := models.GetUserByMobile(form.Syr)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
				return
			}
			syr = suser.UserID
		} else {
			syr = form.Syr
		}
	}
	if form.SuserID != "" {
		syr = form.SuserID
	}
	if err := models.DevAllocate(form.Ids, form.DstJgdm, syr, form.Cfwz, czr, form.Czlx); err != nil {
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
