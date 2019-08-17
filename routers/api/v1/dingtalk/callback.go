package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Callbacks struct {
	Callbacks string `form:"callbacks" json:"callbacks"`
}

type CallbackDetail struct {
	Encrypt string `form:"encrypt" json:"encrypt"`
}

//注册事件回调
func RegisterCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	callbacks := []string{"user_add_org", "user_modify_org", "user_leave_org", "org_dept_create", "org_dept_modify", "org_dept_remove"}
	callbackURL := setting.DingtalkSetting.CallBackHost + "/api/v1/callback/detail"
	request := map[string]interface{}{
		"call_back_tag": callbacks,
		"token":         setting.DingtalkSetting.Token,
		"aes_key":       setting.DingtalkSetting.AesKey,
		"url":           callbackURL,
	}
	response, err := dingtalk.RegisterCallback(request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
		return
	}
	appG.Response(http.StatusBadRequest, e.SUCCESS, response)
}

// 查询事件回调
func QueryCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	response, err := dingtalk.QueryCallback()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
		return
	}
	appG.Response(http.StatusBadRequest, e.SUCCESS, response)
}

// 更新事件回调
func UpdateCallback(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		cbs  Callbacks
	)
	err := c.BindJSON(&cbs)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errcode": 400, "description": "Post Data Err"})
	}
	callbacks := strings.Split(cbs.Callbacks, ",")
	//callbacks := []string{"user_add_org", "user_modify_org", "user_leave_org", "org_dept_create", "org_dept_modify", "org_dept_remove"}
	callbackURL := setting.DingtalkSetting.CallBackHost + "/api/v1/callback/detail"
	request := map[string]interface{}{
		"call_back_tag": callbacks,
		"token":         setting.DingtalkSetting.Token,
		"aes_key":       setting.DingtalkSetting.AesKey,
		"url":           callbackURL,
	}
	response, err := dingtalk.UpdateCallback(request)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
		return
	}
	appG.Response(http.StatusBadRequest, e.SUCCESS, response)
}

// 删除事件回调
func DeleteCallback(c *gin.Context) {
	appG := app.Gin{C: c}
	response, err := dingtalk.DeleteCallback()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
		return
	}
	appG.Response(http.StatusBadRequest, e.SUCCESS, response)
}

// 获取回调失败的结果
func GetFailedCallbacks(c *gin.Context) {
	appG := app.Gin{C: c}
	response, err := dingtalk.GetFailedCallbacks()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if response.ErrCode == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, response)
		return
	}
	appG.Response(http.StatusBadRequest, e.SUCCESS, response)
}

// 获取回调的结果
func GetCallbacks(c *gin.Context) {
	var (
		reply = map[string]interface{}{}
		cbd   CallbackDetail
		res   = "success"
	)
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	err := c.BindJSON(&cbd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errcode": 400, "description": "Post Data Err"})
		return
	}
	var (
		secretMsg = cbd.Encrypt
		token     = setting.DingtalkSetting.Token
		aeskey    = setting.DingtalkSetting.AesKey
		corpid    = setting.DingtalkSetting.CorpID
		dc        = dingtalk.NewDingTalkCrypto(token, aeskey, corpid)
	)
	replyMsg, err := dc.GetDecryptMsg(signature, timestamp, nonce, secretMsg)
	if err != nil {
		logging.Info(fmt.Sprintf("GetDecryptMsg failed:%v", err))
		return
	}
	//logging.Info(fmt.Sprintf("replyMsg is:%v", replyMsg))
	errJson := json.Unmarshal([]byte(replyMsg), &reply)
	if errJson != nil {
		log.Printf("unmarshall replyMsg info error:%v", errJson)
		return
	}
	switch reply["EventType"] {
	case "user_add_org", "user_modify_org":
		for _, userid := range reply["UserId"].([]string) {
			user := dingtalk.UserDetail(userid, 10)
			if err := models.UserDetailSync(user); err != nil {
				log.Printf("UserSync err:%v", err)
				return
			}
		}
	case "user_leave_org":
		for _, userid := range reply["UserId"].([]string) {
			b, err := models.IsUseridExist(userid)
			if err != nil {
				log.Printf("Get IsUseridExist err:%v", err)
				return
			} else {
				if b {
					if err := models.DeleteUser(userid); err != nil {
						log.Printf("DeleteUser err:%v", err)
						return
					}
				} else {
					log.Println("User not exist")
				}
			}
		}
	case "org_dept_create", "org_dept_modify":
		for _, deptIds := range reply["DeptId"].([]string) {
			deptId, _ := strconv.Atoi(deptIds)
			dt := dingtalk.DepartmentDetail(deptId, 10)
			if err := models.DepartmentSync(dt); err != nil {
				log.Printf("DepartmentSync err:%v", err)
				return
			}
		}
	case "org_dept_remove":
		for _, deptIds := range reply["DeptId"].([]string) {
			deptId, _ := strconv.Atoi(deptIds)
			b, err := models.IsDeptIdExist(deptId)
			if err != nil {
				log.Printf("Get IsDeptIdExist err:%v", err)
				return
			} else {
				if b {
					if err := models.DeleteDepartment(deptId); err != nil {
						log.Printf("DepartmentSync err:%v", err)
						return
					}
				} else {
					log.Println("Department not exist")
				}
			}
		}
	case "check_url":
		logging.Info("RegisterCallback check_url")
	default:
		logging.Info(fmt.Sprintf("Unregister Callback,replyMsg is %v", reply))
		return
	}
	crypt, sign, _ := dc.GetEncryptMsg(res, timestamp, nonce)
	data := map[string]interface{}{
		"msg_signature": sign,
		"timeStamp":     timestamp,
		"nonce":         nonce,
		"encrypt":       crypt,
	}
	c.JSON(http.StatusOK, data)
}
