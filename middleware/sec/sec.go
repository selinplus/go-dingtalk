package sec

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var t = &TokenVertify{}

func Sec() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			session = sessions.Default(c)
			rkey    = "E5DOFhZl"
			userID  string
			code    int
		)

		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		log.Println("token=======", token)
		ts := strings.Split(token, ".")
		userID = fmt.Sprintf("%v", session.Get("userid"))

		u := c.Request.URL.Path
		if strings.Index(u, "login") != -1 || strings.Index(u, "js_api_config") != -1 ||
			strings.Index(u, "callback/detail") != -1 || strings.LastIndex(u, "sync") != -1 {
			code = e.SUCCESS
		} else {
			if userID == "" || token == "" {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else { //check token
				sign := ts[0] + rkey + ts[1]
				vertify := util.EncodeMD5(sign)
				if vertify != ts[2] {
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				} else { //judge if token is overtime
					tokenMsg := userID + "." + ts[0] + "." + ts[1]
					timeSmap, _ := strconv.Atoi(ts[0])
					if time.Now().Unix()-int64(timeSmap) < setting.AppSetting.TokenTimeout {
						if t.IsTokenExist(tokenMsg) {
							if strings.Index(u, "/upload/images/") == -1 {
								code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
							} else { //when download files ,token timeout over 30s
								code = e.SUCCESS
							}
						} else {
							code = e.SUCCESS
							t.AddToken(tokenMsg)
							t.DeleteToken()
						}
					} else {
						code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
					}
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

type TokenVertify struct {
	Lock   sync.Mutex
	Tokens []string
}

func (t *TokenVertify) IsTokenExist(tokenMsg string) bool {
	for _, tk := range t.Tokens {
		if tk == tokenMsg {
			return true
		}
	}
	return false
}

func (t *TokenVertify) AddToken(tokenMsg string) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Tokens = append(t.Tokens, tokenMsg)
}

func (t *TokenVertify) DeleteToken() {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	n := time.Now().Unix()
	for i := 0; i < len(t.Tokens); i++ {
		timeSmap, _ := strconv.Atoi(strings.Split(t.Tokens[i], ".")[1])
		if n-int64(timeSmap) > setting.AppSetting.TokenTimeout {
			t.Tokens = append(t.Tokens[:i], t.Tokens[i+1:]...)
			i--
		}
	}
}
