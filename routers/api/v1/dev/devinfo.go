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
}

type OpForm struct {
	Ids     []string `json:"ids"`
	Dms     []string `json:"dms"` //交回&批量收回
	SrcJgdm string   `json:"src_jgdm"`
	DstJgdm string   `json:"dst_jgdm"` //分配,下发,上交
	Lsh     string   `json:"lsh"`      //上交时,用于修改devtodo表done
	Czr     string   `json:"czr"`      //inner传递操作人mobile
	Syr     string   `json:"syr"`      //inner传递使用人mobile
	CuserID string   `json:"cuserid"`  //epp传递操作人userid
	SuserID string   `json:"suserid"`  //epp传递使用人userid
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
		sbbh     = c.Query("sbbh")
		xlh      = c.Query("xlh")
		syr      = c.Query("syr")
		mc       = c.Query("mc")
		jgdm     = c.Query("jgdm")
		bz       = c.Query("bz")
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
		"sbbh":  sbbh,
		"xlh":   xlh,
		"syr":   syr,
		"mc":    mc,
		"jgdm":  jgdm,
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
	resps := make([]*Resp, 0)
	for _, dev := range devs {
		var syrName, syrMobile string
		if dev.Syr != "" {
			suser, err := models.GetUserdemoByUserid(dev.Syr)
			if err != nil {
				log.Println(fmt.Sprintf("根据userid[%s],获取设备使用人失败：%v", dev.Syr, err))
				syrName, syrMobile = dev.Syr, dev.Syr
			} else {
				syrName, syrMobile = suser.Name, suser.Mobile
			}
		}
		d := &Resp{dev, models.ConvSbbhToIdstr(dev.Sbbh), syrName, syrMobile}
		resps = append(resps, d)
	}
	data := make(map[string]interface{})
	data["lists"] = resps
	data["total"] = len(resps)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//导出设备清册
func ExportDevInfos(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		rkrqq    = c.Query("rkrqq")
		rkrqz    = c.Query("rkrqz")
		sbbh     = c.Query("sbbh")
		xlh      = c.Query("xlh")
		syr      = c.Query("syr")
		mc       = c.Query("mc")
		jgdm     = c.Query("jgdm")
		bz       = c.Query("bz")
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
		"sbbh":  sbbh,
		"xlh":   xlh,
		"syr":   syr,
		"mc":    mc,
		"jgdm":  jgdm,
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
	if len(devs) == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	resps := make([]*Resp, 0)
	for _, dev := range devs {
		var syrName, syrMobile string
		if dev.Syr != "" {
			suser, err := models.GetUserdemoByUserid(dev.Syr)
			if err != nil {
				log.Println(fmt.Sprintf("根据userid[%s],获取设备使用人失败：%v", dev.Syr, err))
				syrName, syrMobile = dev.Syr, dev.Syr
			} else {
				syrName, syrMobile = suser.Name, suser.Mobile
			}
		}
		d := &Resp{dev, models.ConvSbbhToIdstr(dev.Sbbh), syrName, syrMobile}
		resps = append(resps, d)
	}
	records := make([]map[string]string, 0)
	for _, resp := range resps {
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
			"存放位置":   resp.Cfwz,
		})
	}
	// sort map key
	sortedKeys := make([]string, 21)
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
		case "存放位置":
			sortedKeys[20] = field
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
		mobile   = c.Query("mobile")
		sbbh     = c.Query("sbbh")
		property = c.Query("property")
		state    = c.Query("state")
		devtype  = c.Query("type")
		xlh      = c.Query("xlh")
		jgdm     = c.Query("jgdm")
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
		}
		ds, err := models.GetDevinfosGly(con)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
			return
		}
		devs = append(devs, ds...)
	}
	resps := make([]*Resp, 0)
	for _, dev := range devs {
		var syrName, syrMobile string
		if dev.Syr != "" {
			suser, err := models.GetUserdemoByUserid(dev.Syr)
			if err != nil {
				log.Println(fmt.Sprintf("根据userid[%s],获取设备使用人失败：%v", dev.Syr, err))
				syrName, syrMobile = dev.Syr, dev.Syr
			} else {
				syrName, syrMobile = suser.Name, suser.Mobile
			}
		}
		d := &Resp{dev, models.ConvSbbhToIdstr(dev.Sbbh), syrName, syrMobile}
		resps = append(resps, d)
	}
	data := make(map[string]interface{})
	data["lists"] = resps
	data["total"] = len(resps)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type Resp struct {
	*models.DevinfoResp
	Idstr     string `json:"idstr"` //6位短编号
	SyrName   string `json:"syr_name"`
	SyrMobile string `json:"syr_mobile"`
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
		var syrName, syrMobile string
		if dev.Syr != "" {
			suser, err := models.GetUserdemoByUserid(dev.Syr)
			if err != nil {
				log.Println(fmt.Sprintf("根据userid[%s],获取设备使用人失败：%v", dev.Syr, err))
				syrName, syrMobile = dev.Syr, dev.Syr
			} else {
				syrName, syrMobile = suser.Name, suser.Mobile
			}
		}
		d := &Resp{dev, models.ConvSbbhToIdstr(dev.Sbbh), syrName, syrMobile}
		appG.Response(http.StatusOK, e.SUCCESS, d)
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
		"rkrqq": "2000-01-01 00:00:00",
		"rkrqz": "2099-01-01 00:00:00",
		"syr":   "",
		"jgdm":  "",
		"mc":    mc,
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
	devResps := make([]*Resp, 0)
	data := make(map[string]interface{})
	for _, dev := range resps {
		var syrName, syrMobile string
		if dev.Syr != "" {
			suser, err := models.GetUserdemoByUserid(dev.Syr)
			if err != nil {
				log.Println(fmt.Sprintf("根据userid[%s]设备使用人失败：%v", dev.Syr, err))
				syrName, syrMobile = dev.Syr, dev.Syr
			}
			syrName, syrMobile = suser.Name, suser.Mobile
		}
		d := &Resp{dev, models.ConvSbbhToIdstr(dev.Sbbh), syrName, syrMobile}
		devResps = append(devResps, d)
	}
	data["lists"] = devResps
	data["total"] = len(devResps)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//设备下发
func Issued(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form OpForm
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
	if err := models.DevIssued(form.Ids, form.SrcJgdm, form.DstJgdm, czr.UserID, "2"); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//设备分配(管理员入库)&借出&收回&交回&上交
func Allocate(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form OpForm
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
	if form.Czlx == "10" { //上交
		if err := models.DevIssued(form.Ids, form.SrcJgdm, form.DstJgdm, czr, form.Czlx); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	} else {
		if err := models.DevAllocate(form.Ids, form.Dms, form.DstJgdm, syr, form.Cfwz, czr, form.Czlx, form.Lsh); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
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
	todolist, err := models.GetUpDevTodosOrDones(done)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	data := make([]interface{}, 0)
	for _, p := range todolist {
		if p.Gly == user.UserID && len(p.DevID) == 0 {
			data = append(data, p)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
