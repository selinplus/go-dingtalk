package h5

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
)

type AddressbookForm struct {
	ID     uint
	Name   string `json:"name" valid:"Required"`
	Mobile string `json:"mobile"`
}

type ContactersForm struct {
	Bookid    uint     `json:"bookid" valid:"Required"`
	UserIDs   []string `json:"userids"`
	Deptnames []string `json:"deptname"`
}

//增加通讯录组
func AddAddressbook(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    AddressbookForm
		userID  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	mobile := form.Mobile
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	book := models.MsgAddressbook{
		UserID: userID,
		Name:   form.Name,
	}
	err = models.AddAddressbook(&book)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除通讯录组
func DeleteAddressbook(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)

	mobile := c.Query("mobile")
	bookid, _ := strconv.Atoi(c.Query("bookid"))
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	//删除通讯录下所有联系人
	err = models.DeleteContacters(uint(bookid))
	//删除通讯录
	err = models.DeleteAddressbook(uint(bookid), userID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改通讯录组名称
func UpdateAddressbook(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		form    AddressbookForm
		userID  string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	mobile := form.Mobile
	var err error
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	book := models.MsgAddressbook{
		ID:     form.ID,
		UserID: userID,
		Name:   form.Name,
	}
	err = models.UpdateAddressbook(&book)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取通讯组列表
func GetAddressbooks(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	mobile := c.Query("mobile")
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	books, err := models.GetAddressbooks(userID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, books)
}

//增加联系人
func AddContacter(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ContactersForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	fail := make([]string, 0)
	for i, userID := range form.UserIDs {
		book := models.MsgContacter{
			BookeID:  form.Bookid,
			UserID:   userID,
			DeptName: form.Deptnames[i],
		}
		err := models.AddContacter(&book)
		if err != nil {
			fail = append(fail, userID)
		}
	}
	if len(fail) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, fail)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除联系人
func DeleteContacter(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ContactersForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	fail := make([]string, 0)
	for _, userID := range form.UserIDs {
		err := models.DeleteContacter(userID, form.Bookid)
		if err != nil {
			fail = append(fail, userID)
		}
	}
	if len(fail) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, fail)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取最近联系人列表
func GetRecentContacter(c *gin.Context) {
	var (
		session = sessions.Default(c)
		appG    = app.Gin{C: c}
		userID  string
	)
	mobile := c.Query("mobile")
	pageSize := 10
	if len(c.Query("pageSize")) > 0 {
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
	}
	if len(mobile) > 0 {
		user, err := models.GetUserByMobile(mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL,
				fmt.Sprintf("根据手机号：%s 获取人员信息错误：%v", mobile, err))
			return
		}
		userID = user.UserID
	} else {
		userID = fmt.Sprintf("%v", session.Get("userid"))
	}
	users, err := models.GetRecentContacter(userID, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, users)
}

//获取通讯组联系人列表
func GetContacters(c *gin.Context) {
	appG := app.Gin{C: c}
	bookid, _ := strconv.Atoi(c.Query("bookid"))
	contacters, err := models.GetContacters(uint(bookid))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if len(contacters) > 0 {
		data := make([]*models.User, 0)
		for _, contacter := range contacters {
			user, err := models.GetUserByUserid(contacter.UserID)
			if err != nil {
				continue
			}
			user.Department = contacter.DeptName
			data = append(data, user)
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
