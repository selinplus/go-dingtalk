package upload

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/selinplus/go-dingtalk/pkg/file"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
)

// GetFileFullUrl get the full access path
func GetFileFullUrl(name string, yaodianID string, mendianID string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetFilePath(yaodianID, mendianID) + name
}

// GetFileName get File name
func GetFileName(name string) string {
	ext := path.Ext(name)
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

// GetFilePath get save path
func GetFilePath(yaodianID string, mendianID string) string {
	today := time.Now().Format("2006-01-02")
	return setting.AppSetting.ImageSavePath + yaodianID + "_" + mendianID + "/" + today + "/"
}

// GetFileFullPath get full save path
func GetFileFullPath(yaodianID string, mendianID string) string {
	return setting.AppSetting.RuntimeRootPath + GetFilePath(yaodianID, mendianID)
}

// CheckFileExt check File file ext
func CheckFileExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// CheckFileSize check File size
func CheckFileSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

// CheckFile check if the file exists
func CheckFile(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
