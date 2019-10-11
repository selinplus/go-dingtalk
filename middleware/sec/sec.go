package sec

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Sec() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userid")
		u := c.Request.URL.Path

		if strings.Index(u, "login") != -1 || strings.Index(u, "js_api_config") != -1 {
			c.Next()
		} else {
			if userID == nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 401,
					"msg":  "鉴权失败，请联系管理员",
					"data": nil,
				})
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
