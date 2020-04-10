package dingtalk

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"log"
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
		i := strings.LastIndex(at.FileUrl, "/")
		fileUrl := at.FileUrl[i+1:]
		a := models.Attachment{
			FileName: at.FileName,
			FileUrl:  fileUrl,
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
		appG.Response(http.StatusOK, e.SUCCESS, msg.ID)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
	}
}

//发信息(内网)
func SendMsgMobile(c *gin.Context) {
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
		i := strings.LastIndex(at.FileUrl, "/")
		fileUrl := at.FileUrl[i+1:]
		a := models.Attachment{
			FileName: at.FileName,
			FileUrl:  fileUrl,
			FileSize: at.FileSize,
			FileType: at.FileType,
			Time:     t,
		}
		ats = append(ats, a)
	}
	Mobile := form.FromID
	user, errm := models.GetUserByMobile(Mobile)
	if errm != nil {
		logging.Info(fmt.Sprintf("%v", errm))
	}
	msg := models.Msg{
		ToID:        form.ToID,
		ToName:      form.ToName,
		FromID:      user.UserID,
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
		appG.Response(http.StatusOK, e.SUCCESS, msg.ID)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
	}
}

type MsgResp struct {
	models.Msg
	DeptName string `json:"dept_name"`
}

//获取消息列表
func GetMsgs(c *gin.Context) {
	var (
		data     = make(map[string]interface{})
		session  = sessions.Default(c)
		appG     = app.Gin{C: c}
		msgs     []*models.Msg
		msgResps []*MsgResp
	)
	tag, _ := strconv.Atoi(c.Query("tag"))
	pageNum, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("size"))
	mobile := c.Query("mobile")
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		msgs, err = models.GetMsgs(user.UserID, uint(tag), pageNum, pageSize)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSGLIST_FAIL, nil)
			return
		}
		if len(msgs) > 0 {
			for _, msg := range msgs {
				var ats []models.Attachment
				for _, at := range msg.Attachments {
					at.FileUrl = upload.GetImageFullUrl(at.FileUrl)
					ats = append(ats, at)
				}
				msg.Attachments = ats
				if msg.ToName == user.Name {
					msg.ToID = mobile
				}
				if msg.FromName == user.Name {
					msg.FromID = mobile
				}
			}
			data["lists"] = msgs
			appG.Response(http.StatusOK, e.SUCCESS, data)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	} else {
		userID := fmt.Sprintf("%v", session.Get("userid"))
		msgs, err = models.GetMsgs(userID, uint(tag), pageNum, pageSize)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSGLIST_FAIL, nil)
			return
		}
		if len(msgs) > 0 {
			for _, msg := range msgs {
				var ats []models.Attachment
				for _, at := range msg.Attachments {
					at.FileUrl = upload.GetAppImageFullUrl(at.FileUrl)
					ats = append(ats, at)
				}
				msg.Attachments = ats
				//add deptName
				var userid string
				switch tag {
				case 0:
					if userID == msg.FromID {
						userid = msg.ToID
					}
					if userID == msg.ToID {
						userid = msg.FromID
					}
				case 1:
					userid = msg.FromID
				case 2:
					userid = msg.ToID
				}
				usr, err := models.GetUserByUserid(userid)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
					return
				}
				deptids := strings.Split(usr.Department, ",")
				id, _ := strconv.Atoi(deptids[0])
				department, err := models.GetDepartmentByID(id)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
					return
				}
				msgResp := MsgResp{
					Msg:      *msg,
					DeptName: department.Name,
				}
				msgResps = append(msgResps, &msgResp)
			}
			data["lists"] = msgResps
			appG.Response(http.StatusOK, e.SUCCESS, data)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	}
}

//获取消息详情
func GetMsgByID(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		msg     *models.Msg
	)
	id, _ := strconv.Atoi(c.Query("id"))
	tag, _ := strconv.Atoi(c.Query("tag"))
	mobile := c.Query("mobile")
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		msg, err = models.GetMsgByID(uint(id), uint(tag), user.UserID)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSG_FAIL, nil)
			return
		}
		if msg.ID > 0 {
			var ats []models.Attachment
			for _, at := range msg.Attachments {
				at.FileUrl = upload.GetImageFullUrl(at.FileUrl)
				ats = append(ats, at)
			}
			msg.Attachments = ats
			if msg.ToName == user.Name {
				msg.ToID = mobile
			}
			if msg.FromName == user.Name {
				msg.FromID = mobile
			}
			appG.Response(http.StatusOK, e.SUCCESS, msg)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	} else {
		userID := fmt.Sprintf("%v", session.Get("userid"))
		msg, err = models.GetMsgByID(uint(id), uint(tag), userID)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSG_FAIL, nil)
			return
		}
		if msg.ID > 0 {
			if !strings.Contains(msg.FromID, userID) && !strings.Contains(msg.ToID, userID) {
				appG.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
				return
			}
			var ats []models.Attachment
			for _, at := range msg.Attachments {
				at.FileUrl = upload.GetAppImageFullUrl(at.FileUrl)
				ats = append(ats, at)
			}
			msg.Attachments = ats
			//add deptName
			var userid string
			switch tag {
			case 0:
				if userID == msg.FromID {
					userid = msg.ToID
				}
				if userID == msg.ToID {
					userid = msg.FromID
				}
			case 1:
				userid = msg.FromID
			case 2:
				userid = msg.ToID
			}
			usr, err := models.GetUserByUserid(userid)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
				return
			}
			deptids := strings.Split(usr.Department, ",")
			id, _ := strconv.Atoi(deptids[0])
			department, err := models.GetDepartmentByID(id)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, nil)
				return
			}
			msgResp := MsgResp{
				Msg:      *msg,
				DeptName: department.Name,
			}
			appG.Response(http.StatusOK, e.SUCCESS, msgResp)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	}
}

//删除消息
func DeleteMsg(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	ids := c.Query("id")
	tag, _ := strconv.Atoi(c.Query("tag"))
	mobile := c.Query("mobile")
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	for _, id := range strings.Split(ids, ",") {
		i, _ := strconv.Atoi(id)
		err = models.DeleteMsg(userID, uint(i), uint(tag))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_MSG_FAIL, id+"删除失败")
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取最近联系人列表
func GetRecentContacter(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	mobile := c.Query("mobile")
	pageSize := 10
	if len(c.Query("pageSize")) > 0 {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	users, err := models.GetRecentContacter(userID, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, users)
}

type AddressbookForm struct {
	ID     uint
	Name   string `json:"name" valid:"Required"`
	Mobile string `json:"mobile"`
}

//增加通讯录组
func AddAddressbook(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    AddressbookForm
		userID  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	mobile := form.Mobile
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	book := models.MsgAddressbook{
		UserID: userID,
		Name:   form.Name,
	}
	err = models.AddAddressbook(&book)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除通讯录组
func DeleteAddressbook(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)

	mobile := c.Query("mobile")
	bookid, _ := strconv.Atoi(c.Query("bookid"))
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	//删除通讯录下所有联系人
	err = models.DeleteContacters(uint(bookid))
	//删除通讯录
	err = models.DeleteAddressbook(uint(bookid), userID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改通讯录组名称
func UpdateAddressbook(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    AddressbookForm
		userID  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	mobile := form.Mobile
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	book := models.MsgAddressbook{
		ID:     form.ID,
		UserID: userID,
		Name:   form.Name,
	}
	err = models.UpdateAddressbook(&book)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取通讯组列表
func GetAddressbooks(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	books, err := models.GetAddressbooks(userID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, books)
}

type ContactersForm struct {
	Bookid    uint     `json:"bookid" valid:"Required"`
	UserIDs   []string `json:"userids"`
	Deptnames []string `json:"deptname"`
}

//增加联系人
func AddContacter(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ContactersForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	fail := make([]string, 0)
	for i, userID := range form.UserIDs {
		book := models.MsgContacter{
			BookeID:  form.Bookid,
			UserID:   userID,
			DeptName: form.Deptnames[i],
		}
		err := models.AddContacter(&book)
		if err != nil {
			fail = append(fail, userID)
		}
	}
	if len(fail) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, fail)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除联系人
func DeleteContacter(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ContactersForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	fail := make([]string, 0)
	for _, userID := range form.UserIDs {
		err := models.DeleteContacter(userID, form.Bookid)
		if err != nil {
			fail = append(fail, userID)
		}
	}
	if len(fail) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, fail)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取通讯组联系人列表
func GetContacters(c *gin.Context) {
	appG := app.Gin{C: c}
	bookid, _ := strconv.Atoi(c.Query("bookid"))
	contacters, err := models.GetContacters(uint(bookid))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(contacters) > 0 {
		data := make([]*models.User, 0)
		for _, contacter := range contacters {
			user, err := models.GetUserByUserid(contacter.UserID)
			if err != nil {
				log.Println(contacter.UserID)
				continue
			}
			user.Department = contacter.DeptName
			data = append(data, user)
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
