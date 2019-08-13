package dingtalk

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"github.com/selinplus/go-gin-web/pkg/logging"
	"log"
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
			log.Printf("session.Save() err:%v", err)
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
	defer func() {
		if r := recover(); r != nil {
			logging.Info(fmt.Sprintf("recover panic in SubDepartmentList:%v", r))
		}
	}()
	appG := app.Gin{C: c}
	var wt = 10 //发生网页劫持后，发送递归请求的次数
	depIds, err := dingtalk.SubDepartmentList(wt)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.SUCCESS, "获取部门id失败，请稍候重试")
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
		syncNum := 30
		wg := &sync.WaitGroup{}
		wg.Add(syncNum)
		for k := 0; k < syncNum; k++ {
			wg.Done()
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logging.Info(fmt.Sprintf("recover panic in DepartmentUserSync:%v", r))
					}
				}()
				for depId := range depIdChan {
					department := dingtalk.DepartmentDetail(depId, wt)
					department.SyncTime = time.Now().Format("2006-01-02 15:04:05")
					if department.ID != 0 {
						if err := models.DepartmentSync(department); err != nil {
							log.Printf("DepartmentSync err:%v", err)
						}
					}
					userids := dingtalk.DepartmentUserIdsDetail(depId, wt)
					cnt := len(userids)
					//log.Printf("userids lenth is:%v", cnt)
					var pageNumTotal int
					if cnt%100 == 0 {
						pageNumTotal = cnt / 100
					} else {
						pageNumTotal = cnt/100 + 1
					}
					for pageNum := 0; pageNum < pageNumTotal; pageNum++ {
						userlist := dingtalk.DepartmentUserDetail(depId, pageNum, wt)
						if err := models.UserSync(userlist); err != nil {
							log.Printf("UserSync err:%v", err)
						}
					}
				}
			}()
		}
		appG.Response(http.StatusOK, e.SUCCESS, "请求发送成功，数据同步中...")
		wg.Wait()
		return
	}
}

//获取部门用户信息同步条数
func DepartmentUserSyncNum(c *gin.Context) {
	appG := app.Gin{C: c}
	depNum, deperr := models.CountDepartmentSyncNum()
	if deperr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_NUMBER_FAIL, nil)
		return
	}
	userNum, usererr := models.CountUserSyncNum()
	if usererr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_NUMBER_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["syncTime"] = time.Now().Format("2006-01-02")
	data["depNum"] = depNum
	data["userNum"] = userNum
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
