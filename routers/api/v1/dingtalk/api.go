package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	authCode := c.PostForm("authCode")
	id := dingtalk.GetUserId(authCode)
	if id != "" {
		userInfo := dingtalk.GetUserInfo(id)
		appG.Response(http.StatusOK, e.SUCCESS, userInfo)
		return
	}
	log.Println("user id is empty:in Login")
}
