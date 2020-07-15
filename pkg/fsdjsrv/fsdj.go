package fsdjsrv

import "github.com/selinplus/go-dingtalk/pkg/setting"

// GetFsdjImageFullUrl get the full access path
func GetFsdjImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/fsdj/" + GetFsdjImagePath() + name
}

// GetFsdjEappImageFullUrl get the full access path for internet eapp:without token
func GetFsdjEappImageFullUrl(name string) string {
	return setting.AppSetting.AppPrefixUrl + "/api/fsdj/" + GetFsdjImagePath() + name
}

// GetFsdjEappImageFullPath get full save path
func GetFsdjEappImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetFsdjImagePath()
}

// GetFsdjImagePath get save path
func GetFsdjImagePath() string {
	return setting.FsdjEappSetting.FsdjSavePath
}
