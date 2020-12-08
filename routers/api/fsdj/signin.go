package fsdj

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strings"
	"time"
)

type SigninResp struct {
	*models.StudySignin
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//签到
func StudySignin(c *gin.Context) {
	appG := app.Gin{C: c}
	var userid string
	if len(c.Query("mobile")) > 0 {
		user, err := models.GetUserdemoByMobile(c.Query("mobile"))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", c.Query("mobile"), err))
			return
		}
		userid = user.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	signIn := &models.StudySignin{
		UserID: userid,
		Qdrq:   time.Now().Format("2006-01-02"),
	}
	if models.IsSinin(signIn) {
		appG.Response(http.StatusOK, e.ERROR, "用户已签到!")
		return
	}
	if err := models.AddStudySignin(signIn); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//党员查看签到情况
func GetSigninsByUserid(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		user   *models.Userdemo
		err    error
		userid string
	)
	if len(c.Query("mobile")) > 0 {
		user, err = models.GetUserdemoByMobile(c.Query("mobile"))
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", c.Query("mobile"), err))
			return
		}
		userid = user.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}

	signins, err := models.GetSigninByUserid(userid)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(signins) > 0 {
		data := make([]*SigninResp, 0)
		if user == nil {
			user, err = models.GetUserdemoByUserid(userid)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL,
					fmt.Sprintf("根据userid：%s 获取人员信息错误：%v", userid, err))
				return
			}
		}
		for _, signin := range signins {
			data = append(data, &SigninResp{
				StudySignin: signin,
				Name:        user.Name,
				Mobile:      user.Mobile,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//查看某天全局签到情况
func GetSigninsByQdrq(c *gin.Context) {
	appG := app.Gin{C: c}
	rq := c.Query("qdrq")
	signins, err := models.GetSigninsByQdrq(rq)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}
	if len(signins) > 0 {
		data := make([]*SigninResp, 0)
		for _, signin := range signins {
			user, err := models.GetUserdemoByUserid(signin.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL,
					fmt.Sprintf("根据userid：%s 获取人员信息错误：%v", signin.UserID, err))
				return
			}
			data = append(data, &SigninResp{
				StudySignin: signin,
				Name:        user.Name,
				Mobile:      user.Mobile,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS,
			map[string]interface{}{
				"list": data,
				"cnt":  len(signins),
			})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
