package fsdj

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"strings"
)

type StudyMemberForm struct {
	ID     uint
	Jgdm   string `json:"dm"`
	Userid string `json:"userid"`
}

type StudyMemberResp struct {
	*models.StudyMember
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//增加学习小组成员
func AddGroupMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form StudyMemberForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	for _, userid := range strings.Split(form.Userid, ",") {
		if !models.IsStudyMemberExist(form.Jgdm, userid) {
			user := models.StudyMember{
				Dm:     form.Jgdm,
				UserID: userid,
			}
			if models.IsMemberExist(userid) {
				appG.Response(http.StatusOK, e.ERROR_ADD_USER_FAIL, "用户已存在!")
				return
			}
			if err := models.AddStudyMember(&user); err != nil {
				appG.Response(http.StatusOK, e.ERROR_ADD_USER_FAIL, err)
				return
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//移动学习小组成员
func UpdGroupMember(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form StudyMemberForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	user := models.StudyMember{
		ID:     form.ID,
		Dm:     form.Jgdm,
		UserID: form.Userid,
	}
	if err := models.UpdStudyMember(&user); err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEPARTMENT_FAIL, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取学习小组成员列表
func GetGroupMembers(c *gin.Context) {
	appG := app.Gin{C: c}
	members, err := models.GetStudyMembers(c.Query("dm"))
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_DEVUSER_FAIL, err)
		return
	}
	if len(members) > 0 {
		resp := make([]*StudyMemberResp, 0)
		for _, member := range members {
			user, err := models.GetUserByUserid(member.UserID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
				return
			}
			resp = append(resp, &StudyMemberResp{
				StudyMember: member,
				Name:        user.Name,
				Mobile:      user.Mobile,
			})
		}
		appG.Response(http.StatusOK, e.SUCCESS, resp)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除学习小组成员
func DelGroupMember(c *gin.Context) {
	appG := app.Gin{C: c}
	ids := c.Query("id")
	errs := ""
	for _, i := range strings.Split(ids, ",") {
		id, _ := strconv.Atoi(i)
		if err := models.DelStudyMember(uint(id)); err != nil {
			errs += i + ","
			continue
		}
	}
	if len(errs) > 0 {
		errId := strings.TrimRight(errs, ",")
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USER_FAIL, errId)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
