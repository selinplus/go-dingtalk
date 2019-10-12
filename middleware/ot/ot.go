package ot

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

//judge if token is overtime
func OT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var t = &TokenVertify{}
		var rkey = "E5DOFhZl"
		var code int

		session := sessions.Default(c)
		userID := fmt.Sprintf("%v", session.Get("userid"))
		token := c.GetHeader("X-Access-Token")
		ts := strings.Split(token, ".")

		u := c.Request.URL.Path
		if strings.Index(u, "login") != -1 || strings.Index(u, "js_api_config") != -1 {
			code = e.SUCCESS
		} else {
			if token == "" {
				code = e.INVALID_PARAMS
			} else {
				sign := ts[0] + rkey + ts[1]
				vertify := util.EncodeMD5(sign)
				if vertify != ts[2] {
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				} else {
					tokenMsg := userID + "." + ts[0] + "." + ts[1]
					timeSmap, _ := strconv.Atoi(ts[0])
					if time.Now().Unix()-int64(timeSmap) < setting.AppSetting.TokenTimeout {
						if t.IsTokenExist(tokenMsg) {
							code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
						} else {
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
	Lock   *sync.Mutex
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
	var tokens []string
	n := time.Now().Unix()
	for i := 0; i < len(t.Tokens); i++ {
		timeSmap, _ := strconv.Atoi(strings.Split(t.Tokens[i], ".")[1])
		if n-int64(timeSmap) > setting.AppSetting.TokenTimeout {
			t.Tokens = append(t.Tokens[:i], t.Tokens[i+1:]...)
			i--
		}
	}
	t.Tokens = tokens
}
