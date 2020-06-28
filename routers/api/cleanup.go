package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/file"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"net/http"
	"os"
)

func CleanUpFile(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		fNum int
	)
	dirpath := upload.GetImageFullPath()
	files, err := file.FindFilesOlderThanDate(dirpath, 365)
	errNotExist := "open : The system cannot find the file specified."
	if err != nil && err.Error() != errNotExist {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	for _, fileInfo := range files {
		err = os.Remove(dirpath + fileInfo.Name())
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
