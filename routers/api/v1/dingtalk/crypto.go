package dingtalk

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"log"
	"net/http"
)

type CallbackDetail struct {
	Encrypt string `form:"encrypt" json:"encrypt"`
}

// 获取回调的结果
func GetCallbacks(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	var cbd CallbackDetail
	err := c.BindJSON(&cbd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errcode": 400, "description": "Post Data Err"})
	} else {
		log.Println(cbd)
	}
	//token := "123456"
	//aeskey := "4g5j64qlyl3zvetqxz5jiocdr586fn2zvjpa8zls3ij"
	//corpid := "suite4xxxxxxxxxxxxxxx"
	secretMsg := cbd.Encrypt
	token := "ytsw3706"
	aeskey := "35KWSX7PyqS3RRvjyMf17SpSD0m03StbqWR3x6wmY8y"
	corpid := setting.DingtalkSetting.CorpID
	dc := dingtalk.NewDingTalkCrypto(token, aeskey, corpid)
	replyMsg, err := dc.GetDecryptMsg(signature, timestamp, nonce, secretMsg)
	if err != nil {
		logging.Info(fmt.Sprintf("GetDecryptMsg failed:%v", err))
	} else {
		logging.Info(fmt.Sprintf("replyMsg is:%v", replyMsg))
	}
	res := "success"
	crypt, sign, _ := dc.GetEncryptMsg(res, timestamp, nonce)
	data := map[string]interface{}{
		"msg_signature": sign,
		"timeStamp":     timestamp,
		"nonce":         nonce,
		"encrypt":       crypt,
	}
	c.JSON(http.StatusOK, data)
}