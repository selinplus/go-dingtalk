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
	"strings"
	"time"
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
	callbackURL := setting.DingtalkSetting.CallBackHost + "/api/v2/callback/detail"
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
	errJson := json.Unmarshal([]byte(replyMsg), &reply)
	//log.Printf("replyMsg is:%v", reply)
	if errJson != nil {
		log.Printf("unmarshall replyMsg info error:%v", errJson)
		return
	}
	switch reply["EventType"] {
	case "user_add_org", "user_modify_org":
		for _, userid := range reply["UserId"].([]interface{}) {
			if user := dingtalk.UserDetail(userid.(string), 10); user != nil {
				if err := models.UserSync(user); err != nil {
					logging.Info(fmt.Sprintf("sync userid:%v err:%v", userid, err))
				}
			} else {
				logging.Info(fmt.Sprintf("%v：get user detail failed!", userid))
			}
		}
	case "user_leave_org":
		for _, userid := range reply["UserId"].([]interface{}) {
			if err := models.DeleteUser(userid.(string)); err != nil {
				logging.Info(fmt.Sprintf("delete userid:%v err:%v", userid, err))
			}
		}
	case "org_dept_create", "org_dept_modify":
		for _, deptId := range reply["DeptId"].([]interface{}) {
			deptid := int(deptId.(float64))
			if dt := dingtalk.DepartmentDetail(deptid, 10); dt != nil {
				dt.SyncTime = time.Now().Format("2006-01-02 15:04:05")
				if err := models.DepartmentSync(dt); err != nil {
					logging.Info(fmt.Sprintf("sync deptid:%v err:%v", deptid, err))
				}
			} else {
				logging.Info(fmt.Sprintf("%v：get department detail failed!", deptId))
			}
		}
	case "org_dept_remove":
		for _, deptId := range reply["DeptId"].([]interface{}) {
			deptid := int(deptId.(float64))
			if err := models.DeleteDepartment(deptid); err != nil {
				logging.Info(fmt.Sprintf("delete deptid:%v err:%v", deptid, err))
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
