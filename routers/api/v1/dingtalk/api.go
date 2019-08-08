package dingtalk

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"log"
	"net/http"
)

type LoginForm struct {
	AuthCode string `json:"auth_code"`
}

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	session := sessions.Default(c)
	var form LoginForm
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if form.AuthCode == "" {
		log.Println("no auth code")
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	id := dingtalk.GetUserId(form.AuthCode)
	if id != "" {
		userInfo := dingtalk.GetUserInfo(id)
		session.Set("userid", userInfo.UserID)
		session.Save()
		appG.Response(http.StatusOK, e.SUCCESS, userInfo)
		return
	}
	log.Println("user id is empty:in Login")
}
func JsApiConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	url := c.Query("url")
	if url == "" {
		log.Println("no url")
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	sign := dingtalk.GetJsApiConfig(url)
	if sign != "" {
		appG.Response(http.StatusOK, e.SUCCESS, sign)
		return
	}
	log.Println("url is empty:in JsApiConfig")
}

//部门用户信息同步
func DepartmentUserSync(c *gin.Context) {
	appG := app.Gin{C: c}
	depIds, err := dingtalk.SubDepartmentList()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
		return
	}
	if depIds != nil {
		depIdChan := make(chan int, 100) //部门id
		exitChan := make(chan bool, 8)   //退出标志
		//开启线程，存入部门id
		go func() {
			for _, depId := range depIds {
				depIdChan <- depId
			}
			close(depIdChan)
		}()
		//开启8个线程，同时获取部门详情
		for i := 0; i < 8; i++ {
			go func() {
				for depId := range depIdChan {
					department := dingtalk.DepartmentDetail(depId)
					models.DepartmentSync(department)
					user := dingtalk.DepartmentUserDetail(depId)
					models.UserSync(user)
					exitChan <- true
				}
			}()
		}
		//开启一个线程，等待所有goroute全部退出
		go func() {
			for i := 0; i < 8; i++ {
				<-exitChan //不需要读值，仅计数
				log.Println("wait goroute", i, "exit")
			}
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}()
	}
}
