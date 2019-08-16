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

	ERROR_ADD_MSG_FAIL     = 40011
	ERROR_GET_MSGLIST_FAIL = 40012
	ERROR_GET_MSG_FAIL     = 40013
	ERROR_DELETE_MSG_FAIL  = 40014

	ERROR_GET_DEPARTMENT_NUMBER_FAIL = 50001
	ERROR_GET_USER_NUMBER_FAIL       = 50002
	ERROR_GET_DEPARTMENT_FAIL        = 50003
	ERROR_GET_USER_FAIL              = 50004
	ERROR_GET_USERBYMOBILE_FAIL      = 50005
)
