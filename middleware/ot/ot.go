package ot

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/middleware/sec"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//judge if token is overtime
func OT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			t    = &sec.TokenVertify{}
			rkey = "E5DOFhZl"
			code int
		)

		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		//log.Printf("token is: %s", token)
		ts := strings.Split(token, ".")
		userID := ts[3]

		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			sign := ts[0] + rkey + ts[1] + userID
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
						code = e.SUCCESS
						t.AddToken(tokenMsg)
						t.DeleteToken()
					}
				} else {
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
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
