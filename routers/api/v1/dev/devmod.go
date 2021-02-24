package dev

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"strings"
)

type DevmodForm struct {
	DevID string `json:"devid"`
	Czlx  string `json:"czlx"`
	Sydw  string `json:"sydw"`
	Syks  string `json:"syks"`
	Syr   string `json:"syr"`
	Cfwz  string `json:"cfwz"`
	Czr   string `json:"czr"`
}

//获取当前操作人所有流水记录
func GetDevMods(c *gin.Context) {
	appG := app.Gin{C: c}
	czr, err := models.GetUserdemoByMobile(c.Query("czr"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
			fmt.Sprintf("根据手机号[%s]获取人员信息失败：%v", c.Query("czr"), err))
		return
	}
	devs, err := models.GetDevMods(czr.UserID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = len(devs)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type LsResp struct {
	*models.DevmodResp
	Xh int `json:"xh"`
}

//根据流水号查询记录
func GetDevModetails(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		pageNo   int
		pageSize int
	)
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
	devs, err := models.GetDevModetails(c.Query("lsh"), pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	resps := make([]*LsResp, 0)
	if len(devs) > 0 {
		for i, dev := range devs {
			dev.Idstr = models.ConvSbbhToIdstr(dev.Sbbh)
			resp := &LsResp{dev, i}
			resps = append(resps, resp)
		}
	}
	data := make(map[string]interface{})
	data["lists"] = resps
	data["total"] = len(resps)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//设备流水记录查询
func GetDevModList(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		devId = c.Query("devid")
		devid string
	)
	if len(devId) > 10 {
		devid = strings.Split(devId, "$")[0]
	} else {
		id, _ := strconv.Atoi(devId)
		devinfo := models.GetDevinfoBySbbh(uint(id))
		devid = devinfo.ID
	}
	devs, err := models.GetDevModsByDevid(devid)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVLIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = len(devs)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
