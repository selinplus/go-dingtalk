package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Attachments struct {
	FileName string `json:"name" form:"name"`
	FileUrl  string `json:"url" form:"url"`
	FileSize int    `json:"size" form:"size"`
	FileType string `json:"type" form:"type"`
}
type MsgSendForm struct {
	ToID        string        `json:"to_id" form:"to_id"`
	ToName      string        `json:"to_name" form:"to_name"`
	FromID      string        `json:"from_id" form:"from_id"`
	FromName    string        `json:"from_name" form:"from_name"`
	Title       string        `json:"title" form:"title"`
	Content     string        `json:"content" form:"content"`
	Attachments []Attachments `json:"fileList" form:"fileList"`
}

//发信息
func MsgSend(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form MsgSendForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	var ats = make([]models.Attachment, 0)
	for _, at := range form.Attachments {
		a := models.Attachment{
			FileName: at.FileName,
			FileUrl:  at.FileUrl,
			FileSize: at.FileSize,
			FileType: at.FileType,
			Time:     t,
		}
		ats = append(ats, a)
	}
	msg := models.Msg{
		ToID:        form.ToID,
		ToName:      form.ToName,
		FromID:      form.FromID,
		FromName:    form.FromName,
		Title:       form.Title,
		Content:     form.Content,
		Time:        t,
		Attachments: ats,
	}
	agentID, _ := strconv.Atoi(setting.MsgAppSetting.AgentID)
	userIdList := strings.Split(msg.ToID, ",")
	link := map[string]interface{}{
		"messageUrl": "http://s.dingtalk.com/market/dingtalk/error_code.php",
		"picUrl":     "@lALOACZwe2Rk",
		"title":      msg.Title,
		"text":       msg.Content,
	}
	msgcotent := map[string]interface{}{
		"msgtype": "link",
		"link":    link,
	}
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": userIdList,
		//"to_all_user":  false,
		"msg": msgcotent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	err := models.AddMsgSend(&msg)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
		return
	}
	if msg.ID == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
		return
	}
	if msg.ID > 0 {
		//appG.Response(http.StatusOK, e.SUCCESS, msg.ID)
		models.AddMsgTag(msg.ID, msg.ToID, msg.FromID)
		asyncsendReturn, _ := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendReturn != nil {
			logging.Info(fmt.Sprintf("%v", asyncsendReturn))
			appG.Response(http.StatusOK, e.SUCCESS, asyncsendReturn)
		}
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
	}
}

//获取消息列表
func GetMsgs(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	tag := com.StrTo(c.Param("tag")).MustInt()
	pageNum := com.StrTo(c.Param("start")).MustInt()
	pageSize := com.StrTo(c.Param("size")).MustInt()
	session := sessions.Default(c)
	v := session.Get("userid")
	userID := v.(uint)
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	cnt, err := models.GetMsgCount(userID, uint(tag))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSGLIST_FAIL, nil)
		return
	}
	pages := 0
	if cnt%pageSize == 0 {
		pages = cnt / pageSize
	} else {
		pages = cnt/pageSize + 1
	}
	if pageNum > pages {
		pageNum = pages
	}
	msgs, err := models.GetMsgs(userID, uint(tag), pageNum, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSGLIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = msgs
	data["total"] = cnt
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取消息详情
func GetMsgByID(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, _ := strconv.Atoi(c.Param("id"))
	tag, _ := strconv.Atoi(c.Param("tag"))
	session := sessions.Default(c)
	v := session.Get("userid")
	userID := v.(uint)
	msg, err := models.GetMsgByID(uint(id), userID, uint(tag))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSG_FAIL, nil)
		return
	}
	if msg.ID > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, msg)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSG_FAIL, nil)
	}
}
