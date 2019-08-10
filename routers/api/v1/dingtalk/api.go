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
		var seg int
		depidsLen := len(depIds)
		if depidsLen%8 == 0 {
			seg = depidsLen / 8
		} else {
			seg = (depidsLen / 8) + 1
		}
		depIdChan := make(chan int, 100) //部门id
		for j := 0; j < 8; j++ {
			segIds := depIds[j*seg : (j+1)*seg]
			var num int
			go func() {
				for _, depId := range segIds {
					depIdChan <- depId
					num++
				}
			}()
			if num == depidsLen {
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
						models.DepartmentSync(department)
					}
					userids := dingtalk.DepartmentUserIdsDetail(depId)
					log.Printf("userids is %v", userids)
					len := math.Ceil(float64(len(userids) % 100))
					for l := 0; l < int(len); l++ {
						userlist := dingtalk.DepartmentUserDetail(depId, l)
						models.UserSync(userlist)
					}
				}
			}()
		}
		wg.Wait()
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
}
