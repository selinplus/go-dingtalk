package sec

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strings"
)

func Sec() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userid")
		//log.Printf("session userid is :%v", userID)
		u := c.Request.URL.Path

		if strings.Index(u, "login") != -1 || strings.Index(u, "js_api_config") != -1 ||
			strings.Index(u, "callback/detail") != -1 {
			c.Next()
		} else {
			if userID == nil {
				code := e.ERROR_AUTH_CHECK_TOKEN_FAIL
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
}
