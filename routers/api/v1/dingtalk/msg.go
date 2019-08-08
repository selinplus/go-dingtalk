package dingtalk

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"net/http"
	"strconv"
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
func SendMsg(c *gin.Context) {
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
	err := models.AddSendMsg(&msg)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
		return
	}
	if msg.ID == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
		return
	}
	if msg.ID > 0 {
		err := models.AddMsgTag(msg.ID, msg.ToID, msg.FromID)
		if err != nil {
			logging.Info(fmt.Sprintf("%v", err))
		}
		tcmprJson := dingtalk.MseesageToDingding(msg.Title, msg.Content, msg.ToID)
		asyncsendReturn := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendReturn != nil {
			if asyncsendReturn.Errcode == 0 {
				models.UpdateMsgFlag(msg.ID)
			}
			appG.Response(http.StatusOK, e.SUCCESS, asyncsendReturn)
		}
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
	}
}

//获取消息列表
func GetMsgs(c *gin.Context) {
	appG := app.Gin{C: c}
	tag, _ := strconv.Atoi(c.Query("tag"))
	pageNum, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("size"))
	session := sessions.Default(c)
	v := session.Get("userid")
	userID := fmt.Sprintf("%v", v)
	msgs, err := models.GetMsgs(userID, uint(tag), pageNum, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSGLIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = msgs
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取消息详情
func GetMsgByID(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, _ := strconv.Atoi(c.Query("id"))
	tag, _ := strconv.Atoi(c.Query("tag"))
	session := sessions.Default(c)
	v := session.Get("userid")
	userID := fmt.Sprintf("%v", v)
	msg, err := models.GetMsgByID(uint(id), uint(tag), userID)
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

//删除消息
func DeleteMsg(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, _ := strconv.Atoi(c.Query("id"))
	tag, _ := strconv.Atoi(c.Query("tag"))
	session := sessions.Default(c)
	v := session.Get("userid")
	userID := fmt.Sprintf("%v", v)
	err := models.DeleteMsg(userID, uint(id), uint(tag))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_MSG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}