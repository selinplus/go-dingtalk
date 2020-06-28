package ydksjwt

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/middleware/h5m"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var t = &h5m.TokenVertify{}

func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			rkey = "ytsw-ydks"
			code int
		)

		token := c.GetHeader("Authorization")
		ts := strings.Split(token, ".")

		u := c.Request.URL.Path
		if strings.Index(u, "inner/") != -1 { //skip inner url
			code = e.SUCCESS
		} else {
			if token == "" {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else { //check token
				sign := ts[0] + rkey + ts[1]
				vertify := util.EncodeMD5(sign)
				if vertify != ts[2] {
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				} else { //judge if token is overtime
					tokenMsg := ts[2] + "." + ts[0] + "." + ts[1]
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
