package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
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
	czr, err := models.GetUserByMobile(c.Query("czr"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
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
	data := make(map[string]interface{})
	data["lists"] = devs
	data["total"] = len(devs)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}