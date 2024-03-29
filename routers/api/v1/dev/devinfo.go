package dev

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/export"
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
	Sbdl int    `json:"sbdl"`
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
	czr, err := models.GetUserdemoByMobile(form.Czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
			fmt.Sprintf("根据手机号[%s]获取操作人信息失败：%v", form.Czr, err))
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
		Sbdl: form.Sbdl,
	}
	if models.IsDevXlhExist(form.Xlh) {
		appG.Response(http.StatusInternalServerError, e.ERROR_XLHEXIST_FAIL, nil)
		return
	}
	//生成二维码
	info := sbbh + "$序列号[" + dev.Xlh + "]$生产商[" + dev.Scs + "]$设备型号[" + dev.Xh + "]$生产日期[" + dev.Scrq + "]$"
	name, _, err := qrcode.GenerateQrWithLogo(info, qrcode.GetQrCodeFullPath())
	if err != nil {
		log.Println(err)
	}
	dev.QrUrl = qrcode.GetQrCodeFullUrl(name)
	err = models.AddDevinfo(dev)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//批量导入
func ImpDevinfos(c *gin.Context) {
	appG := app.Gin{C: c}
	czr := c.Query("czr")
	user, err := models.GetUserdemoByMobile(czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
			fmt.Sprintf("根据手机号[%s]获取操作人信息失败：%v", czr, err))
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	errDev, success, failed, err := models.ImpDevinfos(file, user.UserID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, err)
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

type CfwzForm struct {
	Devid string `json:"devid"`
	Cfwz  string `json:"cfwz"`
}

//更新设备存放位置
func UpdateDevinfoCfwz(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form CfwzForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	updMap := map[string]interface{}{
		"id":   form.Devid,
		"cfwz": form.Cfwz,
	}
	if err := models.EditDevinfoCfwz(updMap); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DEV_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type UpdByAdminForm struct {
	Devid     string `json:"devid"`
	Czrmobile string `json:"mobile"` //操作人手机号
	Jgdm      string `json:"jgdm"`   //设备管理机构代码
	Jgksdm    string `json:"jgksdm"` //设备所属机构代码
	Syr       string `json:"syr"`    //设备使用人代码
	Cfwz      string `json:"cfwz"`   //存放位置
}

//更新设备管理机构、使用人、所属机构、所属位置
func UpdateDevinfoByAdmin(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form UpdByAdminForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	czr, _ := models.GetUserdemoByMobile(form.Czrmobile)
	updMap := map[string]interface{}{
		"id":     form.Devid,
		"jgdm":   form.Jgdm,
		"jgksdm": form.Jgksdm,
		"syr":    form.Syr,
		"cfwz":   form.Cfwz,
		"czr":    czr.UserID,
		"czrq":   time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.EditDevinfoByAdmin(updMap); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DEV_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
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
	czr, err := models.GetUserdemoByMobile(form.Czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
			fmt.Sprintf("根据手机号[%s]获取操作人信息失败：%v", form.Czr, err))
		return
	}
	if form.Syr != "" {
		suser, err := models.GetUserdemoByMobile(form.Syr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号[%s]获取使用人信息失败：%v", form.Syr, err))
			return
		}
		syr = suser.UserID
	}
	var lxdm string
	if strings.Contains(form.Lx, "1") {
		lxdm = form.Lx
	} else {
		lx, _ := models.GetDevtypeByMc(form.Lx)
		lxdm = lx.Dm
	}
	dev := &models.Devinfo{
		ID:   form.ID,
		Zcbh: form.Zcbh,
		Lx:   lxdm,
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
	//生成二维码
	info := form.ID + "$序列号[" + dev.Xlh + "]$生产商[" + dev.Scs + "]$设备型号[" + dev.Xh + "]$生产日期[" + dev.Scrq + "]$"
	name, _, err := qrcode.GenerateQrWithLogo(info, qrcode.GetQrCodeFullPath())
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DEV_FAIL, err)
		return
	}
	dev.QrUrl = qrcode.GetQrCodeFullUrl(name)
	if err = models.EditDevinfo(dev); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DEV_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "修改成功，请重新打印二维码！")
}

type DevinfoPnumForm struct {
	Ids []uint `json:"ids"` //id数组
}

//更新设备二维码打印次数
func UpdateDevinfoPnum(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   DevinfoPnumForm
		errIds []uint
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	for _, id := range form.Ids {
		if err := models.EditDevinfoPnum(id); err != nil {
			errIds = append(errIds, id)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, errIds)
}

//删除设备信息
func DelDevinfo(c *gin.Context) {
	appG := app.Gin{C: c}
	if err := models.DelDevinfo(c.Query("id")); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备列表(inner多条件查询设备)
func GetDevinfos(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		rkrqq    = c.Query("rkrqq")
		rkrqz    = c.Query("rkrqz")
		scrqq    = c.Query("scrqq")
		scrqz    = c.Query("scrqz")
		sbbh     = c.Query("sbbh")
		xlh      = c.Query("xlh")
		syr      = c.Query("syr")
		mc       = c.Query("mc")
		jgdm     = c.Query("jgdm")
		bz       = c.Query("bz")
		zcbh     = c.Query("zcbh")
		sbdl     = c.Query("sbdl")
		pageNo   int
		pageSize int
	)
	if syr != "" {
		user, err := models.GetUserdemoByMobile(syr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号[%s]获取使用人信息失败：%v", syr, err))
			return
		}
		syr = user.UserID
	}
	con := map[string]string{
		"rkrqq": rkrqq,
		"rkrqz": rkrqz,
		"scrqq": scrqq,
		"scrqz": scrqz,
		"sbbh":  sbbh,
		"xlh":   xlh,
		"syr":   syr,
		"mc":    mc,
		"jgdm":  jgdm,
		"zcbh":  zcbh,
		"sbdl":  sbdl,
	}
	if c.Query("pageNo") == "" {
		pageNo = 0
	} else {
		pageNo, _ = strconv.Atoi(c.Query("pageNo"))
	}
	if c.Query("pageSize") == "" {
		pageSize = 0
	} else {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	devs, err := models.GetDevinfos(con, pageNo, pageSize, bz)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = len(devs)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//导出设备清册
func ExportDevInfosGly(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		mobile   = c.Query("mobile")
		sbbh     = c.Query("sbbh")
		property = c.Query("property")
		state    = c.Query("state")
		devtype  = c.Query("type")
		xlh      = c.Query("xlh")
		jgdm     = c.Query("jgdm")
		zcbh     = c.Query("zcbh")
		scrq     = c.Query("scrq")
		rkrq     = c.Query("rkrq")
		rkrqq    = c.Query("rkrqq")
		rkrqz    = c.Query("rkrqz")
		scrqq    = c.Query("scrqq")
		scrqz    = c.Query("scrqz")
		sbdl     = c.Query("sbdl")
		depts    = make([]*models.Devdept, 0)
		err      error
	)
	glydm := make([]string, 0)
	if jgdm == "" {
		if len(mobile) > 0 {
			gly, err := models.GetUserdemoByMobile(mobile)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
					fmt.Sprintf("根据手机号[%s]获取管理员信息失败：%v", mobile, err))
				return
			}
			depts, err = models.GetDevdeptsHasGlyByUserid(gly.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
				return
			}
		} else {
			depts, err = models.GetDevdeptsHasGly()
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
				return
			}
		}
		for _, dept := range depts {
			glydm = append(glydm, dept.Jgdm)
		}
	} else {
		glydm = strings.Split(jgdm, ",")
	}
	devs := make([]*models.DevinfoResp, 0)
	for _, dm := range glydm {
		con := map[string]string{
			"sbbh":     sbbh,
			"property": property,
			"state":    state,
			"type":     devtype,
			"xlh":      xlh,
			"jgdm":     dm,
			"zcbh":     zcbh,
			"scrq":     scrq,
			"rkrq":     rkrq,
			"rkrqq":    rkrqq,
			"rkrqz":    rkrqz,
			"scrqq":    scrqq,
			"scrqz":    scrqz,
			"sbdl":     sbdl,
		}
		ds, err := models.GetDevinfosGly(con)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
			return
		}
		devs = append(devs, ds...)
	}
	records := make([]map[string]string, 0)
	for _, resp := range devs {
		records = append(records, map[string]string{
			"设备编号":   resp.Idstr,
			"资产编号":   resp.Zcbh,
			"设备类型":   resp.Lx,
			"设备型号":   resp.Xh,
			"序列号":    resp.Xlh,
			"设备来源":   resp.Ly,
			"供应商":    resp.Gys,
			"价格":     resp.Jg,
			"生产商":    resp.Scs,
			"生产日期":   resp.Scrq,
			"购入日期":   resp.Grrq,
			"设备报废年限": resp.Bfnx,
			"入库日期":   resp.Rkrq,
			"操作人":    resp.Czr,
			"操作日期":   resp.Czrq,
			"设备状态":   resp.Zt,
			"设备属性":   resp.Sx,
			"使用人":    resp.SyrName,
			"使用人手机号": resp.SyrMobile,
			"设备管理机构": resp.Jgmc,
			"设备所属机构": resp.Ksmc,
			"存放位置":   resp.Cfwz,
		})
	}
	// sort map key
	sortedKeys := make([]string, 22)
	for field := range records[0] {
		switch field {
		case "设备编号":
			sortedKeys[0] = field
		case "资产编号":
			sortedKeys[1] = field
		case "设备类型":
			sortedKeys[2] = field
		case "设备型号":
			sortedKeys[3] = field
		case "序列号":
			sortedKeys[4] = field
		case "设备来源":
			sortedKeys[5] = field
		case "供应商":
			sortedKeys[6] = field
		case "价格":
			sortedKeys[7] = field
		case "生产商":
			sortedKeys[8] = field
		case "生产日期":
			sortedKeys[9] = field
		case "购入日期":
			sortedKeys[10] = field
		case "设备报废年限":
			sortedKeys[11] = field
		case "入库日期":
			sortedKeys[12] = field
		case "操作人":
			sortedKeys[13] = field
		case "操作日期":
			sortedKeys[14] = field
		case "设备状态":
			sortedKeys[15] = field
		case "设备属性":
			sortedKeys[16] = field
		case "使用人":
			sortedKeys[17] = field
		case "使用人手机号":
			sortedKeys[18] = field
		case "设备管理机构":
			sortedKeys[19] = field
		case "设备所属机构":
			sortedKeys[20] = field
		case "存放位置":
			sortedKeys[21] = field
		}
		//sorted_keys = append(sorted_keys, field)
	}
	fileName := "设备清册" + time.Now().Format("150405")
	url, err := export.WriteIntoExecel(fileName, sortedKeys, records)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, url)
}

//获取设备列表(管理员端,多条件查询设备)
func GetDevinfosGly(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		rkrqq    = c.Query("rkrqq")
		rkrqz    = c.Query("rkrqz")
		scrqq    = c.Query("scrqq")
		scrqz    = c.Query("scrqz")
		mobile   = c.Query("mobile")
		sbbh     = c.Query("sbbh")
		property = c.Query("property")
		state    = c.Query("state")
		devtype  = c.Query("type")
		xlh      = c.Query("xlh")
		jgdm     = c.Query("jgdm")
		zcbh     = c.Query("zcbh")
		scrq     = c.Query("scrq")
		rkrq     = c.Query("rkrq")
		sbdl     = c.Query("sbdl")
		depts    = make([]*models.Devdept, 0)
		err      error
	)
	glydm := make([]string, 0)
	if jgdm == "" {
		if len(mobile) > 0 {
			gly, err := models.GetUserdemoByMobile(mobile)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
					fmt.Sprintf("根据手机号[%s]获取管理员信息失败：%v", mobile, err))
				return
			}
			depts, err = models.GetDevdeptsHasGlyByUserid(gly.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
				return
			}
		} else {
			depts, err = models.GetDevdeptsHasGly()
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err)
				return
			}
		}
		for _, dept := range depts {
			glydm = append(glydm, dept.Jgdm)
		}
	} else {
		glydm = strings.Split(jgdm, ",")
	}
	devs := make([]*models.DevinfoResp, 0)
	for _, dm := range glydm {
		con := map[string]string{
			"sbbh":     sbbh,
			"property": property,
			"state":    state,
			"type":     devtype,
			"xlh":      xlh,
			"jgdm":     dm,
			"zcbh":     zcbh,
			"scrq":     scrq,
			"rkrq":     rkrq,
			"rkrqq":    rkrqq,
			"rkrqz":    rkrqz,
			"scrqq":    scrqq,
			"scrqz":    scrqz,
			"sbdl":     sbdl,
		}
		ds, err := models.GetDevinfosGly(con)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
			return
		}
		devs = append(devs, ds...)
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = len(devs)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type Resp struct {
	*models.DevinfoResp
}

//获取设备详情
func GetDevinfoByID(c *gin.Context) {
	appG := app.Gin{C: c}
	id := strings.Split(c.Query("id"), "$")[0]
	dev, err := models.GetDevinfoRespByID(id)
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

//获取交回设备待入库列表
func GetDevinfosToBeStored(c *gin.Context) {
	appG := app.Gin{C: c}
	var userid string
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserdemoByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号[%s],获取人员失败：%v", mobile, err))
			return
		}
		userid = user.UserID
	} else {
		u := c.Request.URL.Path
		if strings.Index(u, "api/v3") != -1 {
			token := c.GetHeader("Authorization")
			auth := c.Query("token")
			if len(auth) > 0 {
				token = auth
			}
			ts := strings.Split(token, ".")
			userid = ts[3]
		}
	}

	devs, err := models.GetDevinfosToBeStored()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEV_FAIL, nil)
		return
	}
	if len(devs) > 0 {
		data := make([]*models.Devinfo, 0)
		for _, dev := range devs {
			ddept, err := models.GetDevdept(dev.Jgdm)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, err)
				continue
			}
			if ddept.Gly == userid {
				user, err := models.GetUserdemoByUserid(dev.Czr)
				if err != nil {
					log.Println(fmt.Sprintf(
						"根据userid[%s],操作人失败：%v", dev.Czr, err))
				} else {
					dev.Czr = user.Name
				}
				data = append(data, dev)
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
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
		appG   = app.Gin{C: c}
		jgdm   = c.Query("jgdm")
		bz     = c.Query("bz")
		mc     = c.Query("mc")
		sbbh   = c.Query("sbbh")
		zcbh   = c.Query("zcbh")
		xlh    = c.Query("xlh")
		sbdl   = c.Query("sbdl")
		userid string
	)
	//使用人查看名下设备
	u := c.Request.URL.Path
	if strings.Index(u, "api/v3") != -1 {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	con := map[string]string{
		"rkrqq": "",
		"rkrqz": "",
		"scrqq": "",
		"scrqz": "",
		"syr":   "",
		"jgdm":  "",
		"mc":    mc,
		"sbbh":  sbbh,
		"zcbh":  zcbh,
		"xlh":   xlh,
		"sbdl":  sbdl,
	}
	resps := make([]*models.DevinfoResp, 0)
	if jgdm != "" {
		for _, dm := range strings.Split(jgdm, ",") {
			con["jgdm"] = dm
			devs, err := models.GetDevinfos(con, 0, 0, bz)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
				return
			}
			resps = append(resps, devs...)
		}
	} else {
		con["syr"] = userid
		devs, err := models.GetDevinfos(con, 0, 0, bz)
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
func Issued(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.OpForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	czr, err := models.GetUserdemoByMobile(form.Czr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
			fmt.Sprintf("根据手机号[%s]获取操作人信息失败：%v", form.Czr, err))
		return
	}
	if err := models.DevIssued(form, czr.UserID, "2"); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//设备机构变更申请
func ChangeJgks(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.OpForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	czr := models.GetCommonGly(form.SrcJgksdm, form.DstJgksdm)
	if err := models.ChangeJgks(form, czr); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//管理员处理设备机构变更申请&交回申请
func GlyChangeJgks(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.OpForm
		czr  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if form.Czr != "" {
		cuser, err := models.GetUserdemoByMobile(form.Czr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号[%s]获取操作人信息失败：%v", form.Czr, err))
			return
		}
		czr = cuser.UserID
	}
	if form.CuserID != "" {
		czr = form.CuserID
	}
	if form.Czlx == "8" { //交回入库
		if err := models.AgreeDevReback(form, czr); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	} else if form.Czlx == "11" { //机构变更
		if err := models.AgreeChangeJgks(form, czr); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//设备分配(管理员入库)&借出&收回&交回申请&上交
func Allocate(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.OpForm
		czr  string
		syr  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if form.Czr != "" {
		cuser, err := models.GetUserdemoByMobile(form.Czr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号[%s]获取操作人信息失败：%v", form.Czr, err))
			return
		}
		czr = cuser.UserID
	}
	if form.CuserID != "" {
		czr = form.CuserID
	}
	if form.Syr != "" {
		if form.Syr != " " {
			suser, err := models.GetUserdemoByMobile(form.Syr)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
					fmt.Sprintf("根据手机号[%s]获取使用人信息失败：%v", form.Syr, err))
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
	if form.Czlx == "8" { //交回申请
		if err := models.DevReback(form, syr, czr); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	} else if form.Czlx == "10" { //上交
		if err := models.DevIssued(form, czr, form.Czlx); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	} else { //设备分配(上交后管理员入库)&借出&收回
		if err := models.DevAllocate(form, syr, czr); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//根据id获取待办&已办详情
func GetDevTodosOrDonesByTodoid(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, _ := strconv.Atoi(c.Query("id"))
	done, _ := strconv.Atoi(c.Query("done"))
	data, err := models.GetDevTodosOrDonesByToid(uint(id), done)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
			fmt.Sprintf("获取待办详情失败：%v", err))
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取待办&已办列表(交回设备)
func GetDevTodosOrDones(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		url    = c.Request.URL.Path
		mobile = c.Query("mobile")
		userid string
		done   int
	)
	if len(mobile) > 0 {
		user, err := models.GetUserdemoByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号[%s]获取人员信息失败：%v", mobile, err))
			return
		}
		userid = user.UserID
	} else {
		if strings.Index(url, "api/v3") != -1 {
			token := c.GetHeader("Authorization")
			auth := c.Query("token")
			if len(auth) > 0 {
				token = auth
			}
			ts := strings.Split(token, ".")
			userid = ts[3]
		}
	}

	if strings.Contains(url, "dev/todolist") {
		done = 0
	}
	if strings.Contains(url, "dev/donelist") {
		done = 1
	}
	donelist, err := models.GetDevTodosOrDones(done)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	data := make([]interface{}, 0)
	if len(donelist) > 0 {
		for _, p := range donelist {
			if p.Czlx == "11" { //机构变更时，应判断操作人是否为登录人
				p.Gly = p.Czrid
			}
			if p.Gly == userid && len(p.DevID) > 0 {
				data = append(data, p)
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取待办列表(上交设备)
func GetUpDevTodosOrDones(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		url    = c.Request.URL.Path
		mobile = c.Query("mobile")
		done   int
	)
	user, err := models.GetUserdemoByMobile(mobile)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
			fmt.Sprintf("根据手机号[%s]获取人员信息失败：%v", mobile, err))
		return
	}
	if strings.Contains(url, "dev/uptodolist") {
		done = 0
	}
	if strings.Contains(url, "dev/updonelist") {
		done = 1
	}
	todolist, err := models.GetUpDevTodosOrDones(done, user.UserID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, todolist)
}
