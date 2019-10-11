package ot

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
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
		var code int

		code = e.SUCCESS
		session := sessions.Default(c)
		userID := fmt.Sprintf("%v", session.Get("userid"))
		token := c.GetHeader("X-Access-Token")
		ts := strings.Split(token, ".")

		u := c.Request.URL.Path
		if strings.Index(u, "login") != -1 || strings.Index(u, "js_api_config") != -1 {
			c.Next()
		} else {
			if token == "" {
				code = e.INVALID_PARAMS
			} else {
				sign := ts[0] + ts[1] + setting.AppSetting.JwtSecret
				vertify := util.EncodeMD5(sign)
				if vertify != ts[2] {
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				} else {
					tokenMsg := userID + "." + ts[0] + "." + ts[1]
					timeSmap, _ := strconv.Atoi(ts[0])
					if time.Now().Unix()-int64(timeSmap) > setting.AppSetting.TokenTimeout {
						code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
						_ = models.DeleteToken(tokenMsg)
					} else {
						if models.IsTokenExist(tokenMsg) {
							code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
						} else {
							_ = models.AddToken(tokenMsg)
						}
					}
				}
			}
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
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
