package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	INVALID_PARSE_FORM:              "解析绑定表单错误",
	INVALID_PARAMS_VERIFY:           "参数校验错误",
	ERROR_EXIST_LICENCE:             "已存在该LICENCE",
	ERROR_EXIST_LICENCE_FAIL:        "获取已存在LICENCE失败",
	ERROR_NOT_EXIST_LICENCE:         "该LICENCE不存在",
	ERROE_VERSION_LOW:               "版本过低，请更新新版本",
	ERROR_USERNAME_PASSWORD:         "用户名密码不正确",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时,请重新登录",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	ERROR_UPLOAD_SAVE_FILE_FAIL:    "保存文件失败",
	ERROR_UPLOAD_CHECK_FILE_FAIL:   "检查文件失败",
	ERROR_UPLOAD_CHECK_FILE_FORMAT: "校验文件错误，文件格式或大小有问题",

	ERROR_USERNAME_EXIST: "用户名已存在",

	ERROR_ADD_MSG_FAIL:     "发送信息失败",
	ERROR_GET_MSGLIST_FAIL: "获取消息列表失败",
	ERROR_GET_MSG_FAIL:     "获取消息失败",
	ERROR_DELETE_MSG_FAIL:  "删除消息失败",

	ERROR_GET_DEPARTMENT_NUMBER_FAIL: "获取部门同步条数失败",
	ERROR_GET_USER_NUMBER_FAIL:       "获取用户同步条数失败",
	ERROR_GET_DEPARTMENT_FAIL:        "获取部门列表失败",
	ERROR_GET_USER_FAIL:              "获取部门用户列表失败",
	ERROR_GET_USERBYMOBILE_FAIL:      "用户手机号不正确",

	ERROR_ADD_DEV_FAIL:     "设备登记失败",
	ERROR_GET_DEV_FAIL:     "获取设备登记信息失败",
	ERROR_GET_DEVLIST_FAIL: "获取设备登记信息列表失败",
	ERROR_UPDATE_DEV_FAIL:  "更新设备登记信息失败",
	ERROR_XLHEXIST_FAIL:    "序列号已存在",

	ERROR_SAVE_PROC_FAIL:         "提报事项保存失败",
	ERROR_ADD_PROC_FAIL:          "提报事项提交失败",
	ERROR_GET_PROCLIST_TODO_FAIL: "获取待办列表失败",
	ERROR_GET_PROCLIST_DONE_FAIL: "获取已办列表失败",
	ERROR_GET_PROC_FAIL:          "获取提报事项失败",
	ERROR_GET_PROCMOD_FAIL:       "获取审批流水记录失败",
	ERROR_ADD_PROCMOD_FAIL:       "审批失败",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
