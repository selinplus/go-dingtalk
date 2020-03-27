package api

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"time"
)

type DutyForm struct {
	Mobile string `json:"mobile"`
	Text   string `json:"text"`
}

//接收值班通知推送消息
func OnDuty(c *gin.Context) {
	appG := app.Gin{C: c}
	var form DutyForm
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	user, err := models.GetUserByMobile(form.Mobile)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
		return
	}
	d := models.Onduty{
		UserID:  user.UserID,
		Content: form.Text,
		Tsrq:    time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.AddOnduty(&d); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
