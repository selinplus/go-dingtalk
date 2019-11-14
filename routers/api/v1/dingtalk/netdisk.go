package dingtalk

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"net/http"
	"os"
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
}

//上传网盘文件
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
	fileUrl := "netdisk" + form.FileUrl[i+1:]
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
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_NDFILE_FAIL, nil)
		return
	}
	if nd.ID > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
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
	nds, err := models.GetNetdiskFileList(userID, treeid, pageNum, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_NDFILELIST_FAIL, nil)
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
func MoveToTrash(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	ids := c.Query("ids")
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
	fail := make([]string, 0)
	for _, id := range strings.Split(ids, ",") {
		i, _ := strconv.Atoi(id)
		file, _ := models.GetNetdiskFileDetail(i)
		if !strings.Contains(file.UserID, userID) {
			appG.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
			return
		}
		file.TreeID = 0 //回收站id=0
		file.Xgrq = time.Now().Format("2006-01-02 15:04:05")
		err = models.UpdateNetdiskFile(file)
		if err != nil {
			fail = append(fail, file.FileName+"删除失败")
		}
	}
	data := map[string]interface{}{
		"success":   len(ids) - len(fail),
		"failed":    len(fail),
		"fail_list": fail,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//删除文件
func DeleteNetdiskFile(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	ids := c.Query("ids")
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
	fail := make([]string, 0)
	for _, id := range strings.Split(ids, ",") {
		i, _ := strconv.Atoi(id)
		file, _ := models.GetNetdiskFileDetail(i)
		if !strings.Contains(file.UserID, userID) {
			appG.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
			return
		}
		err = os.Remove(upload.GetImageFullPath() + file.FileName)
		if err != nil {
			fail = append(fail, file.FileName+"删除失败")
		} else {
			_ = models.DeleteNetdiskFile(i)
		}
	}
	data := map[string]interface{}{
		"success":   len(ids) - len(fail),
		"failed":    len(fail),
		"fail_list": fail,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
