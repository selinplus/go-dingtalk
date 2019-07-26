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
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	ERROR_NOT_EXIST_YAODIAN:        "该药店不存在",
	ERROR_CHECK_EXIST_YAODIAN_FAIL: "检查药店是否存在失败",

	ERROR_NOT_EXIST_MENDIAN:        "该门店不存在",
	ERROR_CHECK_EXIST_MENDIAN_FAIL: "检查门店是否存在失败",
	ERROR_ADD_MENDIAN_FAIL:         "增加门店失败",
	ERROR_DELETE_MENDIAN_FAIL:      "删除门店失败",
	ERROR_EDIT_MENDIAN_FAIL:        "编辑门店失败",
	ERROR_COUNT_MENDIAN_FAIL:       "统计门店失败",
	ERROR_GET_MENDIANS_FAIL:        "获取门店列表失败",
	ERROR_GET_MENDIAN_FAIL:         "获取门店失败",

	ERROR_NOT_EXIST_PATIENT:        "不存在该病人信息",
	ERROR_CHECK_EXIST_PATIENT_FAIL: "",
	ERROR_ADD_PATIENT_FAIL:         "增加病人失败",
	ERROR_DELETE_PATIENT_FAIL:      "删除病人信息失败",
	ERROR_EDIT_PATIENT_FAIL:        "修改病人信息失败",
	ERROR_COUNT_PATIENT_FAIL:       "",
	ERROR_GET_PATIENTS_FAIL:        "获取病人列表失败",
	ERROR_GET_PATIENT_FAIL:         "获取病人信息失败",

	ERROR_NOT_EXIST_PREMEDICINE:        "不存在处方药信息",
	ERROR_CHECK_EXIST_PREMEDICINE_FAIL: "",
	ERROR_ADD_PREMEDICINE_FAIL:         "",
	ERROR_DELETE_PREMEDICINE_FAIL:      "",
	ERROR_EDIT_PREMEDICINE_FAIL:        "",
	ERROR_COUNT_PREMEDICINE_FAIL:       "",
	ERROR_GET_PREMEDICINES_FAIL:        "获取处方药列表失败",
	ERROR_GET_PREMEDICINE_FAIL:         "获取处方药信息失败",

	ERROR_NOT_EXIST_PRESCRIPTION:        "不存在处方信息",
	ERROR_CHECK_EXIST_PRESCRIPTION_FAIL: "当日不存在审核完成处方",
	ERROR_ADD_PRESCRIPTION_FAIL:         "添加处方单失败",
	ERROR_DELETE_PRESCRIPTION_FAIL:      "删除处方失败",
	ERROR_EDIT_PRESCRIPTION_FAIL:        "修改处方单失败",
	ERROR_COUNT_PRESCRIPTION_FAIL:       "",
	ERROR_GET_PRESCRIPTIONS_FAIL:        "获取处方单列表失败",
	ERROR_GET_PRESCRIPTION_FAIL:         "获取处方单信息失败",

	ERROR_USERNAME_EXIST:          "用户名已存在",
	ERROR_NOT_EXIST_YAOSHI:        "药师不存在",
	ERROR_CHECK_EXIST_YAOSHI_FAIL: "",
	ERROR_ADD_YAOSHI_FAIL:         "增加药师失败",
	ERROR_DELETE_YAOSHI_FAIL:      "删除药师失败",
	ERROR_EDIT_YAOSHI_FAIL:        "更新药师信息失败",
	ERROR_COUNT_YAOSHI_FAIL:       "",
	ERROR_GET_YAOSHIS_FAIL:        "获取药师列表失败",
	ERROR_GET_YAOSHI_FAIL:         "获取药师信息失败",

	ERROR_NOT_EXIST_YISHI:        "医师不存在",
	ERROR_CHECK_EXIST_YISHI_FAIL: "",
	ERROR_ADD_YISHI_FAIL:         "增加医师失败",
	ERROR_DELETE_YISHI_FAIL:      "",
	ERROR_EDIT_YISHI_FAIL:        "更新医师失败",
	ERROR_COUNT_YISHI_FAIL:       "",
	ERROR_GET_YISHIS_FAIL:        "获取医师列表失败",
	ERROR_GET_YISHI_FAIL:         "获取医师信息失败",

	ERROR_NOT_EXIST_MEDICINE:        "药品不存在",
	ERROR_CHECK_EXIST_MEDICINE_FAIL: "",
	ERROR_ADD_MEDICINE_FAIL:         "",
	ERROR_DELETE_MEDICINE_FAIL:      "",
	ERROR_EDIT_MEDICINE_FAIL:        "",
	ERROR_COUNT_MEDICINE_FAIL:       "",
	ERROR_GET_MEDICINES_FAIL:        "",
	ERROR_GET_MEDICINE_FAIL:         "",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
