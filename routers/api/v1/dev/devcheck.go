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
	u, err := models.GetUserByMobile(form.Fqr)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	form.Fqr = u.UserID
	if err := models.AddDevCheckTask(&form); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if form.Ckself == "Y" {
		//todo:通知所有人员进行盘点
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取盘点任务列表
func GetDevCkTasks(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		//是否自我盘点:Y,N;普通人员只能查看自我盘点任务;管理员可查看全部,或选择性查看
		flag     = c.Query("flag")
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

	ckTask, err := models.GetDevCheckTask(cond, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(ckTask) > 0 {
		for _, ck := range ckTask {
			u, err := models.GetUserByUserid(ck.Fqr)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
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
		flag     = c.Query("flag")
		ckBz     = c.Query("ck_bz")
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
	cond := fmt.Sprintf("check_id = '%s'", checkId)
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
		cond += fmt.Sprintf(" and pdr = '%s'", userid)
	}
	if ckBz != "" {
		cond += fmt.Sprintf(" and ck_bz = %s", ckBz)
	}

	devckdetails, err := models.GetDevckdetails(cond, pageNo, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(devckdetails) > 0 {
		for _, detail := range devckdetails {
			if detail.Pdr != "" {
				u, err := models.GetUserByUserid(detail.Pdr)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
					return
				}
				detail.Pdr = u.Name
			}
			if detail.Czr != "" {
				u, err := models.GetUserByUserid(detail.Czr)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
					return
				}
				detail.Czr = u.Name
			}
			if detail.Syr != "" {
				u, err := models.GetUserByUserid(detail.Syr)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
					return
				}
				detail.Syr = u.Name
			}
			if detail.SyrJgdm != "" {
				devdept, err := models.GetDevdept(detail.SyrJgdm)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
					return
				}
				detail.SyrJgdm = devdept.Jgmc
			}
		}
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"list": devckdetails,
			"cnt":  models.GetDevckdetailsCnt(cond),
		})
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
	}
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
		if models.CheckSyrSelf(uint(CheckID), DevinfoID, userid) {
			appG.Response(http.StatusOK, e.ERROR, "当前盘点人和设备使用人不一致！")
		}
		err := models.DevCheck(uint(CheckID), map[string]interface{}{
			"pdr": userid, "ck_bz": 1, "cktime": time.Now().Format("2006-01-02 15:04:05")})
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}
