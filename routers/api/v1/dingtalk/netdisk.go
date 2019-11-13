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

type NetdiskForm struct {
	ID       int    `json:"id"`
	Mobile   string `json:"mobile"` //inner useful
	TreeID   int    `json:"tree_id"`
	FileName string `json:"file_name"`
	FileUrl  string `json:"url" form:"url"`
	FileSize int    `json:"file_size"`
	Tag      int    `json:"tag"`
}

//上传文件
func AddNetdiskFile(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    NetdiskForm
		mobile  string
		userID  string
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
	i := strings.LastIndex(form.FileUrl, "/")
	fileUrl := form.FileUrl[i+1:]
	nd := models.Netdisk{
		UserID:   userID,
		TreeID:   form.TreeID,
		FileName: form.FileName,
		FileUrl:  fileUrl,
		FileSize: form.FileSize,
		Xgrq:     t,
	}
	err = models.AddNetdiskFile(&nd)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
		return
	}
	if nd.ID == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
		return
	}
	if nd.ID > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	} else {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_MSG_FAIL, nil)
	}
}

//获取当前文件夹文件列表
func GetFileListByDir(c *gin.Context) {
	var (
		data    = make(map[string]interface{})
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
		err     error
	)
	tag, _ := strconv.Atoi(c.Query("tag"))
	treeid, _ := strconv.Atoi(c.Query("treeid"))
	pageNum, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("size"))
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
	nds, err := models.GetNetdiskFileList(userID, tag, treeid, pageNum, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_MSGLIST_FAIL, nil)
		return
	}
	if len(nds) > 0 {
		data["lists"] = nds
		appG.Response(http.StatusOK, e.SUCCESS, data)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

//移动到回收站
func MoveToTrash(c *gin.Context) {}

//删除文件
func DeleteNetdiskFile(c *gin.Context) {
	var (
		//session = sessions.Default(c)
		appG = app.Gin{C: c}
		//userID  string
	)
	ids := c.Query("id")
	//mobile := c.Query("mobile")
	var err error
	//if len(mobile) > 0 {
	//	user, err := models.GetUserByMobile(mobile)
	//	if err != nil {
	//		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
	//		return
	//	}
	//	userID = user.UserID
	//} else {
	//	userID = fmt.Sprintf("%v", session.Get("userid"))
	//}
	for _, id := range strings.Split(ids, ",") {
		i, _ := strconv.Atoi(id)
		err = models.DeleteNetdiskFile(uint(i))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_MSG_FAIL, id+"删除失败")
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
