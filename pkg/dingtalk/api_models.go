package dingtalk

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

type OpenAPIResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// 会话消息异步发送
type AsyncsendResponse struct {
	OpenAPIResponse
	TaskId int `json:"task_id"`
}

// 回调事件Models
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

//待办任务Models
type FormItemList struct {
	Title   string `json:"title"`   //表单标题
	Content string `json:"content"` //表单内容
}
type WorkrecordAddRequest struct {
	UserID       string         `json:"userid"`
	CreateTime   int64          `json:"create_time"`  //待办事项发起时间
	Title        string         `json:"title"`        //待办标题
	Url          string         `json:"url"`          //待办事项的跳转链接
	PcUrl        string         `json:"pcUrl"`        //pc端跳转url,不传则使用url参数
	FormItemList []FormItemList `json:"formItemList"` //待办事项表单
	PcOpenType   int            `json:"pc_open_type"` //可选,待办的pc打开方式。2表示在pc端打开，4表示在浏览器打开
}
type WorkrecordAddResponse struct {
	OpenAPIResponse
	RecordId string `json:"record_id"`
}
type WorkrecordUpdateRequest struct {
	UserID   string `json:"userid"`
	RecordId string `json:"record_id"`
}
type WorkrecordUpdateResponse struct {
	OpenAPIResponse
	Result bool `json:"result"`
}
type WorkrecordQueryRequest struct {
	UserID string `json:"userid"`
	Offset int    `json:"offset"` //分页游标，从0开始，如返回结果中has_more为true，则表示还有数据，offset再传上一次的offset+limit
	Limit  int    `json:"limit"`  //分页大小，最多50
	Status int    `json:"status"` //待办事项状态，0表示未完成，1表示完成
}
type WorkrecordQueryResponse struct {
	OpenAPIResponse
	Records WorkrecordRecord `json:"records"`
}
type WorkrecordRecord struct {
	HasMore bool             `json:"has_more"` //true表示还有多余的数据
	List    []WorkrecordList `json:"list"`
}
type WorkrecordList struct {
	RecordId   string         `json:"record_id"`   //待办事项id，可用此id调用更新待办的接口
	CreateTime int64          `json:"create_time"` //待办事项发起时间
	Title      string         `json:"title"`       //待办标题
	Url        string         `json:"url"`         //待办事项的跳转链接
	Forms      []FormItemList `json:"forms"`       //待办表单列表
}
