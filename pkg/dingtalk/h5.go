package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"log"
	"strconv"
	"time"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresTime int64  `json:"expires_time"`
}
type UserInfo struct {
	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Department []int  `json:"department"`
	Mobile     string `json:"mobile"`
}

// 企业会话消息异步发送
type TopCorpMessageCorpconversationAsyncsendRequest struct {
	AgentId    int
	UserIdList []string
	DeptIdList []int
	ToAllUser  bool
	Msgcontent interface{}
}
type AsyncsendReturn struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Task_id int    `json:"task_id"`
}

var Token = &AccessToken{}

func GetAccessToken() string {
	t := time.Now().UnixNano()
	if Token == nil || t-Token.ExpiresTime >= 0 {
		_, body, errs := gorequest.New().Get(setting.DingtalkSetting.OapiHost + "/gettoken").
			Query("appkey=" + setting.MsgAppSetting.AppKey).
			Query("appsecret=" + setting.MsgAppSetting.AppSecret).End()
		if len(errs) > 0 {
			log.Printf("get dingtalk access token err:%v", errs[0])
		} else {
			log.Printf("Token is :%s", body)
			err := json.Unmarshal([]byte(body), Token)
			util.ShowError("get token, unmarshall json", err)
		}
	}
	return Token.AccessToken
}
func GetUserId(code string) string {
	type UserID struct {
		UserID string `json:"userid"`
		Errmsg string `json:"errmsg"`
	}
	var userId = UserID{}
	_, body, errs := gorequest.New().Get(setting.DingtalkSetting.OapiHost + "/user/getuserinfo").
		Query("code=" + code).
		Query("access_token=" + GetAccessToken()).End()
	log.Printf("access_token in getuserid is %s", GetAccessToken())
	log.Printf("body in getuserid is %s", body)

	if len(errs) > 0 {
		util.ShowError("get userinfo", errs[0])
		return ""
	} else {
		err := json.Unmarshal([]byte(body), &userId)
		if err != nil {
			log.Printf("unmarshall userid info error:%v", err)
			return ""
		}
		return userId.UserID
	}
}
func GetUserInfo(userId string) *UserInfo {
	var userInfo = UserInfo{}
	_, body, errs := gorequest.New().Get(setting.DingtalkSetting.OapiHost + "/user/get").
		Query("userid=" + userId).
		Query("access_token=" + GetAccessToken()).End()
	if len(errs) > 0 {
		util.ShowError("get userinfo", errs[0])
		return nil
	} else {
		err := json.Unmarshal([]byte(body), &userInfo)
		if err != nil {
			log.Printf("unmarshall userid info error:%v", err)
			return nil
		}
		return &userInfo
	}
}
func getJsApiTicket() string {
	type ApiTicket struct {
		Ticket string `json:"ticket"`
	}
	var apiTicket = ApiTicket{}
	_, body, errs := gorequest.New().Get(setting.DingtalkSetting.OapiHost + "/get_jsapi_ticket?access_token=" + GetAccessToken()).End()
	log.Printf("ticket body is %s\n", body)
	if len(errs) > 0 {
		util.ShowError("GetJsApiTicket:", errs[0])
		return ""
	} else {
		err := json.Unmarshal([]byte(body), &apiTicket)
		if err != nil {
			log.Printf("unmarshall GetJsApiTicket info error:%v", err)
			return ""
		}
		return apiTicket.Ticket
	}
}
func genJsApiSign(ticket string, nonceStr string, timeStamp string, url string) string {
	s := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonceStr, timeStamp, url)
	return util.Sha1Sign(s)
}
func GetJsApiConfig(url string) string {
	var config map[string]string
	ticket := getJsApiTicket()
	log.Printf("ticket is :%s\n", ticket)
	if ticket != "" {
		nonceStr := "dingtalk"
		timeStamp := strconv.Itoa(int(time.Now().UnixNano()))
		sign := genJsApiSign(ticket, nonceStr, timeStamp, url)
		log.Printf("timeStamp is %s\n", timeStamp)
		log.Printf("sign is %s\n", sign)
		config = map[string]string{
			"url":       url,
			"nonceStr":  nonceStr,
			"agentId":   setting.MsgAppSetting.AgentID,
			"timeStamp": timeStamp,
			"corpId":    setting.DingtalkSetting.CorpID,
			"ticket":    ticket,
			"signature": sign,
		}
		bytes, _ := json.Marshal(&config)
		return string(bytes)
	} else {
		return ""
	}
}

func MessageCorpconversationAsyncsend(tcmpar string) (*AsyncsendReturn, error) {
	var asyncsendReturn *AsyncsendReturn
	logging.Info(fmt.Sprintf("%v", tcmpar))
	_, body, errs := gorequest.New().
		Post(setting.DingtalkSetting.OapiHost + "/topapi/message/corpconversation/asyncsend_v2?access_token=" + GetAccessToken()).
		Type("json").Send(tcmpar).End()
	if len(errs) > 0 {
		util.ShowError("MessageCorpconversationAsyncsend failed:", errs[0])
		return nil, nil
	} else {
		err := json.Unmarshal([]byte(body), &asyncsendReturn)
		logging.Info(fmt.Sprintf("%v", asyncsendReturn))
		if err != nil {
			log.Printf("unmarshall asyncsendReturn info error:%v", err)
			return nil, nil
		}
		return asyncsendReturn, nil
	}
	return nil, nil
}
