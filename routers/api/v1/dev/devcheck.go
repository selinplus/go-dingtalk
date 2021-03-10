package dev

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/export"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//新增盘点任务
func GetDevCkTask(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form models.Devcheck
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	u, err := models.GetUserdemoByMobile(form.Fqr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL,
			fmt.Sprintf("新增盘点任务发起人获取失败：%s", form.Fqr))
		return
	}
	form.Fqr = u.UserID
	if err := models.AddDevCheckTask(&form); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	// 通知所有人员进行盘点任务
	go func(id uint, ckself string) {
		models.AddSendDevCkTasks(id, ckself)
	}(form.ID, form.Ckself)

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取盘点任务列表
func GetDevCkTasks(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		//是否自我盘点:Y,N;普通人员只能查看自我盘点任务;管理员可查看全部,或选择性查看
		flag     = c.Query("flag")
		sbdl     = c.Query("sbdl")
		pageSize int
		pageNo   int
	)
	if c.Query("pageNo") == "" {
		pageNo = 1
	} else {
		pageNo, _ = strconv.Atoi(c.Query("pageNo"))
	}
	if c.Query("pageSize") == "" {
		pageSize = 10
	} else {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}

	cond := "ckself like '%'"
	if flag != "" {
		cond = fmt.Sprintf("ckself = '%s'", c.Query("flag"))
	}
	if sbdl != "" {
		if sbdl == "1" {
			cond += ` and sbdl = 1`
		}
		if sbdl == "2" {
			cond += ` and sbdl = 2`
		}
	}
	ckTask, err := models.GetDevCheckTask(cond, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(ckTask) > 0 {
		for _, ck := range ckTask {
			u, err := models.GetUserdemoByUserid(ck.Fqr)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL,
					fmt.Sprintf("盘点任务id[%d],发起人获取失败：%s", ck.ID, ck.Fqr))
				return
			}
			ck.Fqr = u.Name
		}
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"list": ckTask,
			"cnt":  models.GetDevCheckTasksCnt(cond),
		})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取盘点任务清册明细
func GetDevCkDetail(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		checkId = c.Query("check_id")
		//是否仅显示所属人设备:Y:只查看名下设备,仅Eapp可用;不传:显示全部
		flag = c.Query("flag")
		ckBz = c.Query("ck_bz")
		jgdm = c.Query("jgdm")
		syr  = ""
	)
	if flag != "" {
		//使用人查看名下设备
		var userid string
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
		syr = userid
	}
	devckdetails, err := models.GetDevckdetails(checkId, ckBz, syr, jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(devckdetails) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"list": devckdetails,
			"cnt":  len(devckdetails),
		})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//导出盘点任务清册明细
func ExportDevCkDetail(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		checkId = c.Query("check_id")
		ckBz    = c.Query("ck_bz")
		jgdm    = c.Query("jgdm")
		syr     = ""
	)
	devckdetails, err := models.GetDevckdetails(checkId, ckBz, syr, jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(devckdetails) == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	records := make([]map[string]string, 0)
	for _, detail := range devckdetails {
		var ckBz = "未盘点"
		if detail.CkBz == 1 {
			ckBz = "已盘点"
		}
		records = append(records, map[string]string{
			"设备编号":    detail.Sbbh,
			"设备名称":    detail.Mc,
			"生产商":     detail.Scs,
			"设备型号":    detail.Xh,
			"序列号":     detail.Xlh,
			"使用人":     detail.Syr,
			"科室":      detail.SyrJgdm,
			"房间号":     detail.Cfwz,
			"盘点人":     detail.Pdr,
			"盘点日期":    detail.Cktime,
			"盘点状态":    ckBz,
			"资产编号":    detail.Zcbh,
			"设备类型":    detail.Lx,
			"设备使用人":   detail.Syr,
			"设备管理机构":  detail.Jgdm,
			"使用人所在机构": detail.SyrJgdm,
			"设备状态":    detail.Zt,
			"设备属性":    detail.Sx,
			"供应商":     detail.Gys,
			"保管人":     detail.Bgr,
		})
	}
	// sort map key
	sortedKeys := make([]string, 18)
	for field := range records[0] {
		switch field {
		case "盘点状态":
			sortedKeys[0] = field
		case "盘点人":
			sortedKeys[1] = field
		case "盘点日期":
			sortedKeys[2] = field
		case "设备名称":
			sortedKeys[3] = field
		case "设备型号":
			sortedKeys[4] = field
		case "设备编号":
			sortedKeys[5] = field
		case "资产编号":
			sortedKeys[6] = field
		case "设备类型":
			sortedKeys[7] = field
		case "序列号":
			sortedKeys[8] = field
		case "设备使用人":
			sortedKeys[9] = field
		case "设备管理机构":
			sortedKeys[10] = field
		case "使用人所在机构":
			sortedKeys[11] = field
		case "设备状态":
			sortedKeys[12] = field
		case "房间号":
			sortedKeys[13] = field
		case "设备属性":
			sortedKeys[14] = field
		case "生产商":
			sortedKeys[15] = field
		case "供应商":
			sortedKeys[16] = field
		case "保管人":
			sortedKeys[17] = field
		}
		//sorted_keys = append(sorted_keys, field)
	}
	fileName := "设备盘点清册" + time.Now().Format("150405")
	url, err := export.WriteIntoExecel(fileName, sortedKeys, records)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, url)
}

type DevCheckImgForm struct {
	DevId string `json:"devId"`
	Img   string `json:"img"`
}

//设备盘点拍照上传
func DevCheckImg(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form DevCheckImgForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	dev := &models.Devinfo{
		ID:  form.DevId,
		Img: form.Img,
	}
	if err := models.EditDevinfo(dev); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//设备盘点
func GetDevCheck(c *gin.Context) {
	var (
		appG      = app.Gin{C: c}
		DevinfoID = c.Query("devId")
	)
	CheckID, _ := strconv.Atoi(c.Query("checkId"))
	if models.IsChecked(uint(CheckID), DevinfoID) {
		appG.Response(http.StatusOK, e.ERROR, "设备已盘点,无需盘点！")
		return
	}
	var userid string
	u := c.Request.URL.Path
	if strings.Index(u, "api/v3") != -1 {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
		if c.Query("ckself") == "Y" {
			if !models.CheckSyrSelf(uint(CheckID), DevinfoID, userid) {
				appG.Response(http.StatusOK, e.ERROR, "当前盘点人和设备使用人不一致！")
				return
			}
		}
		err := models.DevCheck(uint(CheckID), DevinfoID, map[string]interface{}{
			"pdr": userid, "ck_bz": 1, "cktime": time.Now().Format("2006-01-02 15:04:05")})
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}
