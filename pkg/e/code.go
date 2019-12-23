package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	INVALID_PARAMS_VERIFY = 11002
	INVALID_PARSE_FORM    = 11001

	ERROR_EXIST_LICENCE      = 10001
	ERROR_EXIST_LICENCE_FAIL = 10002
	ERROR_NOT_EXIST_LICENCE  = 10003
	ERROE_VERSION_LOW        = 10004

	ERROR_USERNAME_PASSWORD = 10011
	ERROR_USERNAME_EXIST    = 10012

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004

	ERROR_UPLOAD_SAVE_IMAGE_FAIL    = 30001
	ERROR_UPLOAD_CHECK_IMAGE_FAIL   = 30002
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT = 30003
	ERROR_UPLOAD_SAVE_FILE_FAIL     = 30004
	ERROR_UPLOAD_CHECK_FILE_FAIL    = 30005
	ERROR_UPLOAD_CHECK_FILE_FORMAT  = 30006

	ERROR_ADD_MSG_FAIL     = 40011
	ERROR_GET_MSGLIST_FAIL = 40012
	ERROR_GET_MSG_FAIL     = 40013
	ERROR_DELETE_MSG_FAIL  = 40014

	ERROR_ADD_NOTE_FAIL     = 40021
	ERROR_GET_NOTELIST_FAIL = 40022
	ERROR_GET_NOTE_FAIL     = 40023
	ERROR_DELETE_NOTE_FAIL  = 40024

	ERROR_UPLOAD_NDFILE_FAIL   = 40031
	ERROR_GET_NDFILELIST_FAIL  = 40032
	ERROR_MOVE_TO_TRASH_FAIL   = 40033
	ERROR_DELETE_NDFILE_FAIL   = 40034
	ERROR_GET_DIR_LIST_FAIL    = 40035
	ERROR_ADD_DIRL_FAIL        = 40036
	ERROR_UPDATE_DIR_FAIL      = 40037
	ERROR_DELETE_DIR_FAIL      = 40038
	ERROR_DELETE_DIR_IS_PARENT = 40039
	ERROR_DELETE_DIR_HAS_FILE  = 40040

	ERROR_GET_DEPARTMENT_NUMBER_FAIL = 50001
	ERROR_GET_USER_NUMBER_FAIL       = 50002
	ERROR_GET_DEPARTMENT_FAIL        = 50003
	ERROR_GET_USER_FAIL              = 50004
	ERROR_GET_USERBYMOBILE_FAIL      = 50005
	ERROR_ADD_DEPARTMENT_FAIL        = 50006

	ERROR_ADD_DEV_FAIL     = 60001
	ERROR_GET_DEV_FAIL     = 60002
	ERROR_GET_DEVLIST_FAIL = 60003
	ERROR_UPDATE_DEV_FAIL  = 60004
	ERROR_XLHEXIST_FAIL    = 60005

	ERROR_SAVE_PROC_FAIL         = 70001
	ERROR_ADD_PROC_FAIL          = 70002
	ERROR_GET_PROCLIST_TODO_FAIL = 70003
	ERROR_GET_PROCLIST_DONE_FAIL = 70004
	ERROR_GET_PROC_FAIL          = 70005
	ERROR_GET_PROCMOD_FAIL       = 70006
	ERROR_ADD_PROCMOD_FAIL       = 70007
)
