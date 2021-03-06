package h5

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
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
	OrID     int    `json:"orid"`
	FileName string `json:"file_name"`
	FileUrl  string `json:"url"`
	FileSize int    `json:"file_size"`
}

//上传网盘文件
func AddNetdiskFile(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    NetdiskForm
		userID  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	mobile := form.Mobile
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	spareCap, err := models.GetNetdiskSpareCap(userID)
	if spareCap != -1 && err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_NDFILE_FAIL, err)
		return
	}
	if spareCap == -1 { //Initialize the capacity of the Netdisk
		capacity := 2 * 1024 * 1024 //KB
		if err = models.ModNetdiskCap(userID, capacity); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
		}
		spareCap = capacity
	}
	i := strings.LastIndex(form.FileUrl, "/")
	fileUrl := form.FileUrl[i+1:]
	if form.FileSize > spareCap {
		dirUrl := upload.GetImageFullPath() + fileUrl
		if err = os.Remove(dirUrl); err != nil {
			logging.Error(fmt.Sprintf("delete files:[%v] err:%v",
				fileUrl, err))
		}
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_NDFILE_FAIL, "空间不足!")
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
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
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_NDFILE_FAIL, err)
		return
	}
	if nd.ID > 0 {
		if err = models.ModNetdiskCap(userID, spareCap-form.FileSize); err != nil {
			msg := fmt.Sprintf("上传成功，网盘容量大小修改失败：%v", err.Error())
			appG.Response(http.StatusOK, e.ERROR, msg)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

//修改网盘文件&从回收站恢复
func UpdateNetdiskFile(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    NetdiskForm
		userID  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	mobile := form.Mobile
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	i := strings.LastIndex(form.FileUrl, "/")
	fileUrl := form.FileUrl[i+1:]
	t := time.Now().Format("2006-01-02 15:04:05")
	nd := models.Netdisk{
		ID:       form.ID,
		UserID:   userID,
		FileName: form.FileName,
		FileUrl:  fileUrl,
		FileSize: form.FileSize,
		Xgrq:     t,
	}
	msg := ""
	if form.TreeID == 0 { //if recovery file from trash
		if models.IsDirExist(userID, form.OrID) {
			nd.TreeID = form.OrID
		} else {
			nd.TreeID = 1
			msg = fmt.Sprintf("源文件夹已删除，%s恢复到根目录下", form.FileName)
		}
	}
	if err := models.UpdateNetdiskFile(&nd); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	} else {
		if msg != "" {
			appG.Response(http.StatusOK, e.SUCCESS, msg)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
	}
}

//获取当前文件夹文件列表
func GetFileListByDir(c *gin.Context) {
	var (
		data    = make(map[string]interface{})
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
	)
	treeid, _ := strconv.Atoi(c.Query("treeid"))
	pageNum, _ := strconv.Atoi(c.Query("start"))
	pageSize, _ := strconv.Atoi(c.Query("size"))
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID := user.UserID
		nds, err := models.GetNetdiskFileList(userID, treeid, pageNum, pageSize)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_NDFILELIST_FAIL, err)
			return
		}
		if len(nds) > 0 {
			for _, nd := range nds {
				nd.FileUrl = upload.GetImageFullUrl(nd.FileUrl)
			}
			data["lists"] = nds
			appG.Response(http.StatusOK, e.SUCCESS, data)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	} else {
		userID := fmt.Sprintf("%v", session.Get("userid"))
		nds, err := models.GetNetdiskFileList(userID, treeid, pageNum, pageSize)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_NDFILELIST_FAIL, err)
			return
		}
		if len(nds) > 0 {
			for _, nd := range nds {
				nd.FileUrl = upload.GetAppImageFullUrl(nd.FileUrl)
			}
			data["lists"] = nds
			appG.Response(http.StatusOK, e.SUCCESS, data)
		} else {
			appG.Response(http.StatusOK, e.SUCCESS, nil)
		}
	}
}

//移动到回收站
func MoveToTrash(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	id, _ := strconv.Atoi(c.Query("id"))
	file, _ := models.GetNetdiskFileDetail(id)
	if !strings.Contains(file.UserID, userID) {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	file.TreeID = 0 //回收站id=0
	orid, _ := strconv.Atoi(c.Query("orid"))
	file.OrID = orid
	file.Xgrq = time.Now().Format("2006-01-02 15:04:05")
	if err := models.UpdateNetdiskFile(file); err != nil {
		msg := fmt.Sprintf("%s删除失败,err:%v", file.FileName, err.Error())
		appG.Response(http.StatusOK, e.ERROR, msg)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除文件
func DeleteNetdiskFile(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	fail := make([]string, 0)
	ids := strings.Split(strings.TrimRight(c.Query("ids"), ","), ",")
	for _, id := range ids {
		i, _ := strconv.Atoi(id)
		file, _ := models.GetNetdiskFileDetail(i)
		if !strings.Contains(file.UserID, userID) {
			msg := fmt.Sprintf("%s删除失败,操作人对文件无操作权限", file.FileName)
			fail = append(fail, msg)
			continue
		}
		dirUrl := upload.GetImageFullPath() + file.FileUrl
		if err := os.Remove(dirUrl); err != nil {
			msg := fmt.Sprintf("%s删除失败,err:%v", file.FileName, err.Error())
			fail = append(fail, msg)
		} else {
			if err = models.DeleteNetdiskFile(i); err != nil {
				msg := fmt.Sprintf("%s删除失败,err:%v", file.FileName, err.Error())
				fail = append(fail, msg)
				continue
			}
			spareCap, err := models.GetNetdiskSpareCap(file.UserID)
			if err = models.ModNetdiskCap(file.UserID, spareCap+file.FileSize); err != nil {
				msg := fmt.Sprintf("文件删除成功，网盘容量大小修改失败：%v", err.Error())
				fail = append(fail, msg)
				continue
			}
		}
	}
	data := map[string]interface{}{
		"success":   len(ids) - len(fail),
		"failed":    len(fail),
		"fail_list": fail,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
