package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"log"
)

type OpenAPIResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type CallBackResponse struct {
	OpenAPIResponse
}

type QueryCallbackResponse struct {
	OpenAPIResponse
	CallbackTag []string `json:"call_back_tag"`
	Token       string   `json:"token"`
	AesKey      string   `json:"aes_key"`
	URL         string   `json:"url"`
}

type GetFailedCallbackResponse struct {
	OpenAPIResponse
	HasMore    bool              `json:"has_more"`
	FailedList []FailedCallbacks `json:"failed_list"`
}

type FailedCallbacks struct {
	EventTime   int      `json:"event_time"`
	CallbackTag string   `json:"call_back_tag"`
	UserID      []string `json:"userid"`
	CorpID      string   `json:"corpid"`
}

// 注册事件回调接口
func RegisterCallback(request map[string]interface{}) (*CallBackResponse, error) {
	var data CallBackResponse
	_, body, errs := gorequest.New().
		Post(setting.DingtalkSetting.OapiHost + "/call_back/register_call_back?access_token=" + GetAccessToken()).
		Type("json").Send(request).End()
	if len(errs) > 0 {
		util.ShowError("CBRegisterCallback failed:", errs[0])
		return nil, nil
	} else {
		err := json.Unmarshal([]byte(body), &data)
		//logging.Info(fmt.Sprintf("%v", data))
		if err != nil {
			log.Printf("unmarshall CallBackResponse info error:%v", err)
			return nil, err
		}
	}
	return &data, nil
}

// 查询事件回调接口
func QueryCallback() (*QueryCallbackResponse, error) {
	var data QueryCallbackResponse
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/call_back/get_call_back?access_token=" + GetAccessToken()).End()
	if len(errs) > 0 {
		util.ShowError("CBQueryCallback failed:", errs[0])
		return nil, nil
	} else {
		err := json.Unmarshal([]byte(body), &data)
		logging.Info(fmt.Sprintf("%v", data))
		if err != nil {
			log.Printf("unmarshall CBQueryCallback info error:%v", err)
			return nil, err
		}
	}
	return &data, nil
}

// 更新事件回调接口
func UpdateCallback(request map[string]interface{}) (*CallBackResponse, error) {
	var data CallBackResponse
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/call_back/delete_call_back?access_token=" + GetAccessToken()).
		Type("json").Send(request).End()
	if len(errs) > 0 {
		util.ShowError("CBUpdateCallback failed:", errs[0])
		return nil, nil
	} else {
		err := json.Unmarshal([]byte(body), &data)
		logging.Info(fmt.Sprintf("%v", data))
		if err != nil {
			log.Printf("unmarshall CBUpdateCallback info error:%v", err)
			return nil, err
		}
	}
	return &data, nil
}

// 删除事件回调接口
func DeleteCallback() (*CallBackResponse, error) {
	var data CallBackResponse
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/call_back/delete_call_back?access_token=" + GetAccessToken()).End()
	if len(errs) > 0 {
		util.ShowError("CBDeleteCallback failed:", errs[0])
		return nil, nil
	} else {
		err := json.Unmarshal([]byte(body), &data)
		logging.Info(fmt.Sprintf("%v", data))
		if err != nil {
			log.Printf("unmarshall CBDeleteCallback info error:%v", err)
			return nil, err
		}
	}
	return &data, nil
}

// 获取回调失败的结果
func GetFailedCallbacks() (*GetFailedCallbackResponse, error) {
	var data GetFailedCallbackResponse
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/call_back/get_call_back_failed_result?access_token=" + GetAccessToken()).End()
	if len(errs) > 0 {
		util.ShowError("CBGetFailedCallbacks failed:", errs[0])
		return nil, nil
	} else {
		err := json.Unmarshal([]byte(body), &data)
		logging.Info(fmt.Sprintf("%v", data))
		if err != nil {
			log.Printf("unmarshall CBGetFailedCallbacks info error:%v", err)
			return nil, err
		}
	}
	return &data, nil
}

//main方法启动时注册回调接口
func RegCallbackInit() {
	callbacks := []string{"user_add_org", "user_modify_org", "user_leave_org", "org_dept_create", "org_dept_modify", "org_dept_remove"}
	callbackURL := setting.DingtalkSetting.CallBackHost + "/api/v2/callback/detail"
	request := map[string]interface{}{
		"call_back_tag": callbacks,
		"token":         setting.DingtalkSetting.Token,
		"aes_key":       setting.DingtalkSetting.AesKey,
		"url":           callbackURL,
	}
	response, err := RegisterCallback(request)
	if err != nil {
		return
	}
	if response.ErrCode == 0 {
		logging.Info(fmt.Sprintf("RegisterCallback success!"))
	} else {
		logging.Info(fmt.Sprintf("RegisterCallback failed:%v!", response.ErrMsg))
	}
}
