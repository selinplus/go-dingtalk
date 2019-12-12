package dingtalk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/parnurzeal/gorequest"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"log"
	"strconv"
	"strings"
	"sync"
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
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
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
	//log.Printf("access_token in getuserid is %s", GetAccessToken())
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
	//log.Printf("ticket body is %s\n", body)
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
	//log.Printf("ticket is :%s\n", ticket)
	if ticket != "" {
		nonceStr := "dingtalk"
		timeStamp := fmt.Sprintf("%d", time.Now().Unix())
		sign := genJsApiSign(ticket, nonceStr, timeStamp, url)
		//log.Printf("timeStamp is %s\n", timeStamp)
		//log.Printf("sign is %s\n", sign)
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
		log.Printf("%v", config)
		return string(bytes)
	} else {
		return ""
	}
}

//生成工作通知消息体
func MseesageToDingding(msg *models.Msg) string {
	agentID, _ := strconv.Atoi(setting.MsgAppSetting.AgentID)
	link := map[string]interface{}{
		"messageUrl": setting.AppSetting.DingtalkMsgUrl + "#/?id=" + strconv.Itoa(int(msg.ID)),
		"picUrl":     "@lALOACZwe2Rk",
		"title":      msg.Title,
		"text":       msg.Content,
	}
	msgcotent := map[string]interface{}{
		"msgtype": "link",
		"link":    link,
	}
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": msg.ToID,
		"to_all_user": false,
		"msg":         msgcotent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	//log.Println("tcmprJson is", tcmprJson)
	return tcmprJson
}

//生成流程提报待办通知消息体
func ProcessMseesageToDingding(p *models.ProcResponse, czr string) string {
	agentID, _ := strconv.Atoi(setting.MsgAppSetting.AgentID)
	var content string
	if p.Title == "" {
		content = fmt.Sprintf(
			"您有一条提报事项待办消息，推送人是%v，描述信息：%v", p.Tbr, p.Xq)
	} else {
		content = fmt.Sprintf(
			"您有一条提报事项待办消息，推送人是%v，标题为：%v，描述信息：%v", p.Tbr, p.Title, p.Xq)
	}
	text := map[string]interface{}{
		"content": content,
	}
	msgcotent := map[string]interface{}{
		"msgtype": "text",
		"text":    text,
	}
	user, _ := models.GetUserByMobile(czr)
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": user.UserID,
		"to_all_user": false,
		"msg":         msgcotent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	return tcmprJson
}

//生成流程提报补充描述通知消息体
func ProcessBcmsMseesageToDingding(p *models.ProcResponse) string {
	agentID, _ := strconv.Atoi(setting.MsgAppSetting.AgentID)
	link := map[string]interface{}{
		"messageUrl": fmt.Sprintf("eapp://pages/bcms/bcms?id=%v", p.ID),
		"picUrl":     "@lALOACZwe2Rk",
		"title":      p.Title,
		"text":       "您的提报描述不完整，请进行补充描述！",
	}
	msgcotent := map[string]interface{}{
		"msgtype": "link",
		"link":    link,
	}
	user, _ := models.GetUserByMobile(p.Mobile)
	tcmpr := map[string]interface{}{
		"agent_id":    agentID,
		"userid_list": user.UserID,
		"to_all_user": false,
		"msg":         msgcotent,
	}
	tcmprBytes, _ := json.Marshal(&tcmpr)
	tcmprJson := string(tcmprBytes)
	//log.Println("tcmprJson is", tcmprJson)
	return tcmprJson
}

// 企业会话消息异步发送
type AsyncsendReturn struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	TaskId  int    `json:"task_id"`
}

// 企业会话消息异步发送
func MessageCorpconversationAsyncsend(mpar string) *AsyncsendReturn {
	var asyncsendReturn AsyncsendReturn
	_, body, errs := gorequest.New().
		Post(setting.DingtalkSetting.OapiHost + "/topapi/message/corpconversation/asyncsend_v2?access_token=" + GetAccessToken()).
		Type("json").Send(mpar).End()
	if len(errs) > 0 {
		util.ShowError("MessageCorpconversationAsyncsend failed:", errs[0])
		return nil
	} else {
		err := json.Unmarshal([]byte(body), &asyncsendReturn)
		//log.Println("asyncsendReturn is", asyncsendReturn)
		if err != nil {
			log.Printf("unmarshall asyncsendReturn info error:%v", err)
			return nil
		}
	}
	return &asyncsendReturn
}

// 获取子部门Id列表
func SubDepartmentList(wt int) ([]int, error) {
	var (
		subDeptIdList = map[string]interface{}{}
		deptIds       []int
		err           error
	)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/department/list?access_token=" + GetAccessToken()).End()
	if len(errs) > 0 {
		util.ShowError("get department list_ids failed:", errs[0])
		wt = wt - 1
		if wt >= 0 {
			time.Sleep(time.Second * 10)
			deptIds, err = SubDepartmentList(wt)
			return deptIds, err
		} else {
			return nil, err
		}
	}
	err = json.Unmarshal([]byte(body), &subDeptIdList)
	if err != nil {
		if strings.Contains(body, "<") {
			wt = wt - 1
			if wt >= 0 {
				time.Sleep(time.Second * 10)
				deptIds, err = SubDepartmentList(wt)
				return deptIds, err
			} else {
				return nil, err
			}
		}
		log.Printf("unmarshall SubDeptIdList error:%v", err)
		return nil, err
	}
	if subDeptIdList["department"] != nil {
		depts := subDeptIdList["department"].([]interface{})
		for _, v := range depts {
			vv := v.(map[string]interface{})
			for k, val := range vv {
				if k == "id" {
					deptIds = append(deptIds, int(val.(float64)))
					break
				}
			}
		}
		//logging.Info(fmt.Sprintf("deptIds length is %d", len(deptIds)))
	}
	return deptIds, nil
}

// 获取部门详情
func DepartmentDetail(id, wt int) *models.Department {
	var department models.Department
	deptId := strconv.Itoa(id)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/department/get?access_token=" + GetAccessToken() + "&id=" + deptId).End()
	if len(errs) > 0 {
		util.ShowError("get department failed:", errs[0])
		wt = wt - 1
		if wt >= 0 {
			time.Sleep(time.Second * 10)
			dt := DepartmentDetail(id, wt)
			return dt
		} else {
			return nil
		}
	}
	err := json.Unmarshal([]byte(body), &department)
	if err != nil {
		if strings.Contains(body, "<") {
			wt = wt - 1
			if wt >= 0 {
				time.Sleep(time.Second * 10)
				dt := DepartmentDetail(id, wt)
				return dt
			} else {
				return nil
			}
		}
		log.Printf("unmarshall department info error:%v", err)
	}
	return &department
}

// 获取部门用户详情
func DepartmentUserDetail(id, pageNum, wt int) *[]models.User {
	var (
		usersList []models.User
		user      models.User
		userlist  = map[string]interface{}{}
	)
	offset := strconv.Itoa(pageNum * 100)
	deptId := strconv.Itoa(id)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/user/listbypage").
		Query("access_token=" + GetAccessToken()).Query("department_id=" + deptId).
		Query("offset=" + offset).Query("size=100").
		End()
	if len(errs) > 0 {
		util.ShowError("get user failed:", errs[0])
		wt = wt - 1
		if wt >= 0 {
			time.Sleep(time.Second * 10)
			dt := DepartmentUserDetail(id, pageNum, wt)
			return dt
		} else {
			return nil
		}
	}
	err := json.Unmarshal([]byte(body), &userlist)
	if err != nil {
		if strings.Contains(body, "<") {
			wt = wt - 1
			if wt >= 0 {
				time.Sleep(time.Second * 10)
				dt := DepartmentUserDetail(id, pageNum, wt)
				return dt
			} else {
				return nil
			}
		}
		log.Printf("unmarshall userlist error:%v", err)
		return nil
	}
	if userlist["userlist"] != nil {
		users := userlist["userlist"].([]interface{})
		for _, v := range users {
			vv := v.(map[string]interface{})
			_ = mapstructure.Decode(vv, &user)
			user.SyncTime = time.Now().Format("2006-01-02 15:04:05")
			for k, val := range vv {
				if k == "department" {
					var paramSlice []string
					for _, d := range val.([]interface{}) {
						v := strconv.Itoa(int(d.(float64)))
						paramSlice = append(paramSlice, v)
					}
					deptIds := strings.Join(paramSlice, ",")
					user.Department = deptIds
				}
			}
			usersList = append(usersList, user)
		}
	}
	return &usersList
}

//获取部门用户userid列表
func DepartmentUserIdsDetail(id, wt int) []string {
	var (
		useridslice []string
		useridslist = map[string]interface{}{}
	)
	deptId := strconv.Itoa(id)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/user/getDeptMember").
		Query("access_token=" + GetAccessToken()).Query("deptId=" + deptId).
		End()
	if len(errs) > 0 {
		util.ShowError("get userids failed:", errs[0])
		wt = wt - 1
		if wt >= 0 {
			time.Sleep(time.Second * 10)
			dt := DepartmentUserIdsDetail(id, wt)
			return dt
		} else {
			return nil
		}
	}
	err := json.Unmarshal([]byte(body), &useridslist)
	if err != nil {
		if strings.Contains(body, "<") {
			wt = wt - 1
			if wt >= 0 {
				time.Sleep(time.Second * 10)
				dt := DepartmentUserIdsDetail(id, wt)
				return dt
			} else {
				return nil
			}
		}
		log.Printf("unmarshall useridslist error:%v", err)
		return nil
	}
	if useridslist["userIds"] != nil {
		userids := useridslist["userIds"].([]interface{})
		for _, param := range userids {
			useridslice = append(useridslice, param.(string))
		}
	}
	return useridslice
}

//获取用户详情
func UserDetail(userid string, wt int) *models.User {
	var (
		user     models.User
		userlist = map[string]interface{}{}
	)
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/user/get").
		Query("access_token=" + GetAccessToken()).Query("userid=" + userid).
		End()
	if len(errs) > 0 {
		util.ShowError("get user failed:", errs[0])
		wt = wt - 1
		if wt >= 0 {
			time.Sleep(time.Second * 10)
			u := UserDetail(userid, wt)
			return u
		} else {
			return nil
		}
	}
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		if strings.Contains(body, "<") {
			wt = wt - 1
			if wt >= 0 {
				time.Sleep(time.Second * 10)
				u := UserDetail(userid, wt)
				return u
			} else {
				return nil
			}
		}
		log.Printf("convert struct error:%v", errs)
		return nil
	}
	user.SyncTime = time.Now().Format("2006-01-02 15:04:05")
	err = json.Unmarshal([]byte(body), &userlist)
	if err != nil {
		log.Printf("unmarshall user info error:%v", err)
		return nil
	}
	var deptIds string
	if len(userlist) > 0 {
		for k, val := range userlist {
			if k == "department" {
				var paramSlice []string
				for _, d := range val.([]interface{}) {
					v := strconv.Itoa(int(d.(float64)))
					paramSlice = append(paramSlice, v)
				}
				deptIds = strings.Join(paramSlice, ",")
				break
			}
		}
		user.Department = deptIds
	}
	return &user
}

// 获取企业员工人数
func OrgUserCount(wt int) (int, error) {
	var data map[string]interface{}
	_, body, errs := gorequest.New().
		Get(setting.DingtalkSetting.OapiHost + "/user/get_org_user_count?access_token=" + GetAccessToken() + "&onlyActive=0").End()
	if len(errs) > 0 {
		util.ShowError("get user count failed:", errs[0])
		wt = wt - 1
		if wt >= 0 {
			time.Sleep(time.Second * 10)
			count, err := OrgUserCount(wt)
			return count, err
		} else {
			return 0, errs[0]
		}
	}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		if strings.Contains(body, "<") {
			wt = wt - 1
			if wt >= 0 {
				time.Sleep(time.Second * 10)
				count, e := OrgUserCount(wt)
				return count, e
			} else {
				return 0, err
			}
		}
		log.Printf("unmarshall OrgUserCount error:%v", err)
	}
	if int(data["errcode"].(float64)) != 0 {
		return 0, errors.New(fmt.Sprintf("%v", data["errmsg"]))
	}
	return int(data["count"].(float64)), nil
}
