package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"log"
	"strconv"
	"time"
)

// 小程序Token(用于发送工作通知)
func GetEappAccessToken() string {
	var eappToken = &AccessToken{}
	_, body, errs := gorequest.New().Get(setting.DingtalkSetting.OapiHost + "/gettoken").
		Query("appkey=" + setting.EAppSetting.AppKey).
		Query("appsecret=" + setting.EAppSetting.AppSecret).End()
	if len(errs) > 0 {
		log.Printf("get dingtalk access token err:%v", errs[0])
	} else {
		err := json.Unmarshal([]byte(body), eappToken)
		util.ShowError("get token, unmarshall json", err)
	}
	return eappToken.AccessToken
}

//生成流程提报待办通知消息体
func ProcessMseesageToDingding(p *models.ProcResponse, czr string) string {
	agentID, _ := strconv.Atoi(setting.EAppSetting.AgentID)
	var content string
	if p.Title == "" {
		if p.Bcms == "" {
			content = fmt.Sprintf(
				"您有一条提报事项待办消息，推送人是%v，详情：%v", p.Tbr, p.Xq)
		} else {
			content = fmt.Sprintf(
				"您有一条提报事项待办消息，推送人是%v，补充描述为：%v", p.Tbr, p.Bcms)
		}
	} else {
		if p.Bcms == "" {
			content = fmt.Sprintf(
				"您有一条提报事项待办消息，推送人是%v，标题为：%v，详情：%v", p.Tbr, p.Title, p.Xq)
		} else {
			content = fmt.Sprintf(
				"您有一条提报事项待办消息，推送人是%v，标题为：%v，补充描述为：%v", p.Tbr, p.Title, p.Bcms)
		}
	}
	text := map[string]interface{}{
		"content": content,
	}
	msgcontent := map[string]interface{}{
		"msgtype": "text",
		"text":    text,
	}
	user, _ := models.GetUserByMobile(czr)
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": user.UserID,
		"to_all_user": false,
		"msg":         msgcontent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	//log.Println("tcmprJson is", tcmprJson)
	return tcmprJson
}

//生成流程提报补充描述通知消息体
func ProcessBcmsMseesageToDingding(p *models.ProcResponse) string {
	agentID, _ := strconv.Atoi(setting.EAppSetting.AgentID)
	var text string
	if p.Title == "" {
		text = ":请对提报事项进行补充描述"
	} else {
		text = fmt.Sprintf(":请对标题为%s的提报事项进行补充描述", p.Title)
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	link := map[string]interface{}{
		"messageUrl": fmt.Sprintf("eapp://pages/myreport/myreport?id=%v", p.ID),
		"picUrl":     "@lALOACZwe2Rk",
		"title":      "您的提报描述不够准确，请进行补充描述！",
		"text":       t + text,
	}
	msgcontent := map[string]interface{}{
		"msgtype": "link",
		"link":    link,
	}
	user, _ := models.GetUserByMobile(p.Mobile)
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": user.UserID,
		"to_all_user": false,
		"msg":         msgcontent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	//log.Println("tcmprJson is", tcmprJson)
	return tcmprJson
}

//生成交回设备信息通知消息体
func DeviceDingding(devid, gly, done string) string {
	agentID, _ := strconv.Atoi(setting.EAppSetting.AgentID)
	link := map[string]interface{}{
		"messageUrl": fmt.Sprintf("eapp://pages/myreport/myreport?sbid=%s&done=%s", devid, done),
		"picUrl":     "@lALOACZwe2Rk",
		"title":      "交回设备待入库",
		"text":       "请将交回设备入库",
	}
	msgcontent := map[string]interface{}{
		"msgtype": "link",
		"link":    link,
	}
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": gly,
		"to_all_user": false,
		"msg":         msgcontent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	//log.Println("tcmprJson is", tcmprJson)
	return tcmprJson
}

//生成上交设备信息通知消息体
func UpDeviceDingding(num int, jgmc, gly string) string {
	agentID, _ := strconv.Atoi(setting.EAppSetting.AgentID)
	text := map[string]interface{}{
		"content": fmt.Sprintf("%s上交了%d台设备，请在内网管理平台确认入库！",
			jgmc, num),
	}
	msgcontent := map[string]interface{}{
		"msgtype": "text",
		"text":    text,
	}
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": gly,
		"to_all_user": false,
		"msg":         msgcontent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	//log.Println("tcmprJson is", tcmprJson)
	return tcmprJson
}

// 企业会话消息异步发送
func EappMessageCorpconversationAsyncsend(mpar string) *AsyncsendReturn {
	var asyncsendReturn AsyncsendReturn
	_, body, errs := gorequest.New().
		Post(setting.DingtalkSetting.OapiHost + "/topapi/message/corpconversation/asyncsend_v2?access_token=" + GetEappAccessToken()).
		Type("json").Send(mpar).End()
	if len(errs) > 0 {
		util.ShowError("MessageCorpconversationAsyncsend failed:", errs[0])
		return nil
	} else {
		err := json.Unmarshal([]byte(body), &asyncsendReturn)
		log.Println("asyncsendReturn is", asyncsendReturn)
		if err != nil {
			log.Printf("unmarshall asyncsendReturn info error:%v", err)
			return nil
		}
	}
	return &asyncsendReturn
}
