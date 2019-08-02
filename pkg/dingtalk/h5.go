package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"time"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresTime int64  `json:"expires_time"`
}
type UserInfo struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Avatar string `json:"avatar"`
}

var Token = &AccessToken{}
var goReq = gorequest.New()

func GetAccessToken() string {
	t := time.Now().UnixNano()
	if Token == nil || t-Token.ExpiresTime >= 0 {
		_, body, errs := goReq.Get(setting.DingtalkSetting.OapiHost+"/gettoken").
			Param("appKey", setting.MsgAppSetting.AppKey).
			Param("appSecret", setting.MsgAppSetting.AppSecret).End()
		if len(errs) > 0 {
			log.Printf("get dingtalk access token err:%v", errs[0])
		} else {
			err := json.Unmarshal([]byte(body), Token)
			util.ShowError("get token, unmarshall json", err)
		}
	}
	return Token.AccessToken
}
func GetUserId(code string) string {
	type UserID struct {
		UserID string `json:"user_id"`
		Errmsg string `json:"errmsg"`
	}
	var userId = UserID{}
	_, body, errs := goReq.Get(setting.DingtalkSetting.OapiHost+"/user/getuserinfo").
		Param("code", code).
		Param("access_token", GetAccessToken()).End()
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
	_, body, errs := goReq.Get(setting.DingtalkSetting.OapiHost+"/user/get").
		Param("userid", userId).
		Param("access_token", GetAccessToken()).End()
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
	_, body, errs := goReq.Get(setting.DingtalkSetting.OapiHost+"/get_jsapi_ticket").
		Param("access_token", GetAccessToken()).End()
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
	if ticket != "" {
		nonceStr := "dingtalk"
		timeStamp := string(time.Now().UnixNano())
		config = map[string]string{
			"url":       url,
			"nonceStr":  nonceStr,
			"agentId":   setting.MsgAppSetting.AgentID,
			"timeStamp": timeStamp,
			"corpId":    setting.DingtalkSetting.CorpID,
			"ticket":    ticket,
			"signature": genJsApiSign(ticket, nonceStr, timeStamp, url),
		}
		bytes, _ := json.Marshal(&config)
		return string(bytes)
	} else {
		return ""
	}
}
