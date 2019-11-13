package dingtalk

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NoteForm struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Mobile  string `json:"mobile"` //inner useable
}

//新建记事本
func AddNote(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    NoteForm
		userID  string
		mobile  string
		err     error
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	mobile = form.Mobile
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
	if models.IsSameTitle(userID, form.Title) {
		n := models.SimilarTitle(userID, form.Title)
		if n != nil {
			if strings.Contains(n.Title, ")") {
				beg := strings.Index(form.Title, "(")
				end := strings.LastIndex(form.Title, ")")
				i, _ := strconv.Atoi(form.Title[beg+1 : end])
				form.Title = form.Title[:beg+1] + strconv.Itoa(i+1) + ")"
			} else {
				form.Title = form.Title + "(1)"
			}
		}
	}
	note := models.Note{
		Title:   form.Title,
		Content: form.Content,
		UserID:  userID,
		Xgrq:    t,
	}
	err = models.AddNote(&note)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_NOTE_FAIL, nil)
		return
	}
	if note.ID > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, note.ID)
	} else {
		appG.Response(http.StatusOK, e.ERROR_ADD_NOTE_FAIL, nil)
	}
}

//删除记事本
func DeleteNote(c *gin.Context) {
	appG := app.Gin{C: c}
	id, _ := strconv.Atoi(c.Query("id"))
	err := models.DeleteNote(uint(id))
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_NOTE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改记事本内容
func UpdateNote(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    NoteForm
		userID  string
		mobile  string
		err     error
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	mobile = form.Mobile
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
	t := time.Now().Format("2006-01-02 15:04:05")
	note := models.Note{
		ID:      form.ID,
		Title:   form.Title,
		Content: form.Content,
		UserID:  userID,
		Xgrq:    t,
	}
	err = models.UpdateNote(&note)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPDATE_DEV_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取记事本列表
func GetNoteList(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
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
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	notes, err := models.GetNoteList(userID, pageNum, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_NOTELIST_FAIL, nil)
		return
	}
	if len(notes) > 0 {
		data := map[string]interface{}{
			"lists": notes,
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

//查询记事本详情
func GetNoteDetail(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
	)
	id, _ := strconv.Atoi(c.Query("id"))
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		note, err := models.GetNoteDetail(uint(id))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_NOTE_FAIL, nil)
			return
		}
		if note.ID > 0 {
			appG.Response(http.StatusOK, e.SUCCESS, note)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	} else {
		userID := fmt.Sprintf("%v", session.Get("userid"))
		note, err := models.GetNoteDetail(uint(id))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_NOTE_FAIL, nil)
			return
		}
		if note.ID > 0 {
			if !strings.Contains(note.UserID, userID) {
				appG.Response(http.StatusUnauthorized, e.ERROR_GET_NOTE_FAIL, nil)
				return
			}
			appG.Response(http.StatusOK, e.SUCCESS, note)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	}
}
