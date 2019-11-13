package api

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/file"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"golang.org/x/exp/errors/fmt"
	"net/http"
	"os"
)

func CleanUpFile(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		fNum int
	)
	dirpath := setting.AppSetting.RuntimeRootPath + setting.AppSetting.ImageSavePath
	files, err := file.FindFilesOlderThanDate(dirpath, int64(365))
	errNotExist := "open : The system cannot find the file specified."
	if err != nil && err.Error() != errNotExist {
		appG.Response(http.StatusOK, e.ERROR, err.Error())
		return
	}
	for _, fileInfo := range files {
		err = os.Remove(dirpath + "/" + fileInfo.Name())
		if err != nil {
			logging.Error(fmt.Sprintf("clean up files err:%v", err))
			fNum++
		}
	}
	data := map[string]int{
		"success": len(files) - fNum,
		"failed":  fNum,
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
