package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"log"
	"sync"
	"time"
)

// 获取 ydks 项目 AccessToken
func GetYdksAccessToken() string {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	t := time.Now().UnixNano()
	if Token == nil || t-Token.ExpiresTime >= 0 {
		_, body, errs := gorequest.New().Get(setting.DingtalkSetting.OapiHost + "/gettoken").
			Query("appkey=" + setting.YdksAppSetting.AppKey).
			Query("appsecret=" + setting.YdksAppSetting.AppSecret).End()
		if len(errs) > 0 {
			log.Printf("get dingtalk access token err:%v", errs[0])
		} else {
			//log.Printf("Token is :%s", body)
			err := json.Unmarshal([]byte(body), Token)
			util.ShowError("get token, unmarshall json", err)
		}
	}
	return Token.AccessToken
}

// 创建待办任务
func YdksWorkrecordAdd(reqJson string) (*WorkrecordAddResponse, error) {
	_, body, errs := gorequest.New().
		Post(setting.DingtalkSetting.OapiHost + "/topapi/workrecord/add?access_token=" + GetYdksAccessToken()).
		Type(gorequest.TypeJSON).Send(reqJson).End()
	if len(errs) > 0 {
		util.ShowError("workrecord add err:", errs[0])
		return nil, errs[0]
	} else {
		resp := &WorkrecordAddResponse{}
		err := util.FormJson(body, resp)
		return resp, err
	}
}

// 更新任务状态
func YdksWorkrecordUpdate(userid, record_id string) (*WorkrecordUpdateResponse, error) {
	var req = WorkrecordUpdateRequest{
		UserID:   userid,
		RecordId: record_id,
	}
	reqJson, err := util.ToJson(req)
	fmt.Println(reqJson)
	if err != nil {
		return nil, err
	}
	_, body, errs := gorequest.New().
		Post(setting.DingtalkSetting.OapiHost + "/topapi/workrecord/update?access_token=" + GetYdksAccessToken()).
		Type(gorequest.TypeJSON).Send(reqJson).End()
	//log.Printf("body is %s\n", body)
	if len(errs) > 0 {
		util.ShowError("workrecord update err:", errs[0])
		return nil, errs[0]
	} else {
		resp := &WorkrecordUpdateResponse{}
		err = util.FormJson(body, resp)
		return resp, err
	}
}
