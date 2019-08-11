package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"log"
	"strconv"
	"strings"
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
			//log.Printf("Token is :%s", body)
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
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/get_jsapi_ticket?access_token=" + GetAccessToken()).End()
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

//发送工作通知
func MseesageToDingding(title, text, userid_list string) string {
	agentID, _ := strconv.Atoi(setting.MsgAppSetting.AgentID)
	link := map[string]interface{}{
		"messageUrl": "http://s.dingtalk.com/market/dingtalk/error_code.php",
		"picUrl":     "@lALOACZwe2Rk",
		"title":      title,
		"text":       text,
	}
	msgcotent := map[string]interface{}{
		"msgtype": "link",
		"link":    link,
	}
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": userid_list,
		"to_all_user": false,
		"msg":         msgcotent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	return tcmprJson
}

// 企业会话消息异步发送
type AsyncsendReturn struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Task_id int    `json:"task_id"`
}

// 企业会话消息异步发送
func MessageCorpconversationAsyncsend(mpar string) *AsyncsendReturn {
	var asyncsendReturn *AsyncsendReturn
	logging.Info(fmt.Sprintf("%v", mpar))
	_, body, errs := gorequest.New().
		Post(setting.DingtalkSetting.OapiHost + "/topapi/message/corpconversation/asyncsend_v2?access_token=" + GetAccessToken()).
		Type("json").Send(mpar).End()
	if len(errs) > 0 {
		util.ShowError("MessageCorpconversationAsyncsend failed:", errs[0])
		return nil
	} else {
		err := json.Unmarshal([]byte(body), &asyncsendReturn)
		logging.Info(fmt.Sprintf("%v", asyncsendReturn))
		if err != nil {
			log.Printf("unmarshall asyncsendReturn info error:%v", err)
			return nil
		}
	}
	return asyncsendReturn
}

// 获取子部门Id列表
func SubDepartmentList() ([]int, error) {
	var depIds []int
	var subDeptIdList = map[string]interface{}{}
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/department/list?access_token=" + GetAccessToken()).End()
	if len(errs) > 0 {
		util.ShowError("get department list_ids failed:", errs[0])
		return nil, errs[0]
	} else {
		err := json.Unmarshal([]byte(body), &subDeptIdList)
		if err != nil {
			log.Printf("unmarshall SubDeptIdList info error_body is:%v", body)
			log.Printf("unmarshall SubDeptIdList info error:%v", err)
			return nil, err
		}
	}
	depts := subDeptIdList["department"].([]interface{})
	for _, v := range depts {
		vv := v.(map[string]interface{})
		for k, val := range vv {
			if k == "id" {
				depIds = append(depIds, int(val.(float64)))
				break
			}
		}
	}
	log.Printf("depIds length is %d", len(depIds))
	return depIds, nil
}

// 获取部门详情
func DepartmentDetail(id int) *models.Department {
	var department *models.Department
	depId := strconv.Itoa(id)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/department/get?access_token=" + GetAccessToken() + "&id=" + depId).End()
	if len(errs) > 0 {
		util.ShowError("get department failed:", errs[0])
		return nil
	} else {
		err := json.Unmarshal([]byte(body), &department)
		if err != nil {
			log.Printf("unmarshall department info error_body is:%v", body)
			log.Printf("unmarshall department info error:%v", err)
		}
	}
	return department
}

// 获取部门用户详情
func DepartmentUserDetail(id, pageNum int) []*models.User {
	var usersList []*models.User
	var user models.User
	var userlist = map[string]interface{}{}
	//offset := strconv.Itoa(pageNum * 100)
	depId := strconv.Itoa(id)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/user/listbypage").
		Query("access_token=" + GetAccessToken()).Query("department_id=" + depId).
		Query("offset=0").Query("size=100").
		End()
	if len(errs) > 0 {
		util.ShowError("get user failed:", errs[0])
		return nil
	} else {
		log.Println("body is %v", body)
		err := json.Unmarshal([]byte(body), &userlist)
		if err != nil {
			log.Printf("unmarshall userlist info error_body is:%v", body)
			log.Printf("unmarshall userlist info error:%v", err)
		}
		if userlist["userlist"] != nil {
			users := userlist["userlist"].([]interface{})
			for _, v := range users {
				vv := v.(map[string]interface{})
				log.Println("vv is %v", vv)
				mapstructure.Decode(vv, &user)
				user.SyncTime = time.Now().Format("2006-01-02 15:04:05")
				for k, val := range vv {
					if k == "department" {
						var paramSlice []string
						for _, d := range val.([]interface{}) {
							v := strconv.Itoa(int(d.(float64)))
							paramSlice = append(paramSlice, v)
						}
						depIds := strings.Join(paramSlice, ",")
						user.Department = depIds
					}
				}
				log.Printf("user is:%v", user)
			}
		}
		usersList = append(usersList, &user)
	}
	return usersList
}

//获取部门用户userid列表
func DepartmentUserIdsDetail(id int) []string {
	var useridslist = map[string]interface{}{}
	depId := strconv.Itoa(id)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/user/getDeptMember").
		Query("access_token=" + GetAccessToken()).Query("deptId=" + depId).
		End()
	if len(errs) > 0 {
		util.ShowError("get userids failed:", errs[0])
		return nil
	} else {
		err := json.Unmarshal([]byte(body), &useridslist)
		if err != nil {
			log.Printf("unmarshall useridslist info error_body is:%v", body)
			log.Printf("unmarshall useridslist info error:%v", err)
		}
		if useridslist["userIds"] != nil {
			userids := useridslist["userIds"].([]interface{})
			var useridslice []string
			for _, param := range userids {
				useridslice = append(useridslice, param.(string))
			}
			return useridslice
		}
	}
	return nil
}

//获取用户详情
func UserDetail(userid string) *models.User {
	var user *models.User
	var userlist = map[string]interface{}{}
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/user/get").
		Query("access_token=" + GetAccessToken()).Query("userid=" + userid).
		End()
	if len(errs) > 0 {
		util.ShowError("get user failed:", errs[0])
		return nil
	} else {
		err := json.Unmarshal([]byte(body), &userlist)
		if err != nil {
			log.Printf("unmarshall user info error_body is:%v", body)
			log.Printf("unmarshall user info error:%v", err)
		}
		errs := json.Unmarshal([]byte(body), &user)
		if errs != nil {
			log.Printf("convert struct error:%v", err)
		}
		var depIds string
		if len(userlist) > 0 {
			for k, val := range userlist {
				if k == "department" {
					var paramSlice []string
					for _, d := range val.([]interface{}) {
						v := strconv.Itoa(int(d.(float64)))
						paramSlice = append(paramSlice, v)
					}
					depIds = strings.Join(paramSlice, ",")
					break
				}
			}
			user.Department = depIds
		}
		user.SyncTime = time.Now().Format("2006-01-02 15:04:05")
		//log.Printf("users is:%v", user)
	}
	return user
}
