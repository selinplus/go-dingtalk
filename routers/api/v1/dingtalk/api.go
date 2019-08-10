package dingtalk

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
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
		if err := session.Save(); err != nil {
			log.Println("session.Save() err:%v", err)
		}
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
		appG.Response(http.StatusBadRequest, e.SUCCESS, err)
		return
	}
	if depIds != nil {
		depIdChan := make(chan int, 100) //部门id
		for j := 0; j < 8; j++ {
			var num int
			go func() {
				for _, depId := range depIds {
					depIdChan <- depId
					num++
				}
			}()
			if num == len(depIds) {
				close(depIdChan)
			}
		}
		syncNum := 15
		wg := &sync.WaitGroup{}
		wg.Add(syncNum)
		for k := 0; k < syncNum; k++ {
			wg.Done()
			go func() {
				for depId := range depIdChan {
					department := dingtalk.DepartmentDetail(depId)
					department.SyncTime = time.Now().Format("2006-01-02 15:04:05")
					log.Printf("departmen is %v", department)
					if department.ID != 0 {
						_ = models.DepartmentSync(department)
					}
					userids := dingtalk.DepartmentUserIdsDetail(depId)
					log.Printf("userids is %v", userids)
					pageNum := int(math.Ceil(float64(len(userids) % 100)))
					for l := 0; l < pageNum; l++ {
						userlist := dingtalk.DepartmentUserDetail(depId, l)
						//models.UserSync(userlist)
						for _, user := range userlist {
							if models.IsUseridExist(user.UserID) {
								ue := models.EditUser(user)
								log.Printf("update user err %v", ue)
							}
							er := models.AddUser(user)
							log.Printf("add user err %v", er)
						}
					}
				}
			}()
		}
		wg.Wait()
		appG.Response(http.StatusOK, e.SUCCESS, "请求发送成功，数据同步中...")
		return
	}
}
