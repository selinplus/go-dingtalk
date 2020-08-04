package fsdj

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"log"
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

type LoginForm struct {
	AuthCode string `json:"auth_code"`
}

type UserInfoResp struct {
	*dingtalk.UserInfo
	GroupDm   string `json:"group_dm"`
	GroupName string `json:"group_name"`
}

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
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
		user, err := models.GetStudyMember(id)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err)
			return
		}
		if user != nil {
			group, err := models.GetStudyGroup(user.Dm)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR, err)
				return
			}
			appG.Response(http.StatusOK, e.SUCCESS, UserInfoResp{
				UserInfo:  dingtalk.GetUserInfo(id),
				GroupDm:   group.Dm,
				GroupName: group.Mc,
			})
			return
		}
		appG.Response(http.StatusOK, e.ERROR, "用户不在党小组中，请联系管理员添加")
		return
	}
	log.Println("userid is empty:in Login")
	appG.Response(http.StatusBadRequest, e.SUCCESS, nil)
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
			if member.UserID == "fsdj_admin" {
				resp = append(resp, &StudyMemberResp{
					StudyMember: member,
					Name:        "超级管理员",
					Mobile:      "0000",
				})
			} else {
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

//模糊查询福山区用户(简单排除)
func GetFsdjUserByMc(c *gin.Context) {
	appG := app.Gin{C: c}
	mc := c.Query("mc")
	if len(mc) == 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, "名称不能为空")
		return
	}
	users, err := models.GetFsdjUserByMc(mc)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err)
		return
	}
	if len(users) == 0 {
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	var data []*models.User
	for _, user := range users {
		if strings.Contains(user.Name, "_") ||
			strings.Contains(user.Name, "烟台") ||
			strings.Contains(user.Name, "芝罘") ||
			strings.Contains(user.Name, "莱山") ||
			strings.Contains(user.Name, "福山") ||
			strings.Contains(user.Name, "牟平") ||
			strings.Contains(user.Name, "保税") ||
			strings.Contains(user.Name, "开发区") ||
			strings.Contains(user.Name, "海阳") ||
			strings.Contains(user.Name, "栖霞") ||
			strings.Contains(user.Name, "招远") ||
			strings.Contains(user.Name, "莱阳") ||
			strings.Contains(user.Name, "莱州") ||
			strings.Contains(user.Name, "海阳") ||
			strings.Contains(user.Name, "龙口") ||
			strings.Contains(user.Name, "蓬莱") ||
			strings.Contains(user.Name, "长岛") ||
			strings.Contains(user.Name, "公司") ||
			strings.Contains(user.Name, "财务") ||
			strings.Contains(user.Name, "服务部") ||
			strings.Contains(user.Name, "商店") ||
			strings.Contains(user.Name, "超市") ||
			strings.Contains(user.Name, "法人") {
			continue
		}
		deptId, _ := strconv.Atoi(user.Department)
		//if deptId == 29464267 {
		//	data = append(data, user)
		//	continue
		//}
		if deptId < 29464267 {
			continue
		}
		dept, _ := models.GetDepartmentByID(deptId)
		if dept.OuterDept {
			continue
		}
		//if IsFsq(dept.Parentid) {
		user.Department = dept.Name
		data = append(data, user)
		//}
	}
	if len(data) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func IsFsq(deptId int) bool {
	dept, _ := models.GetDepartmentByID(deptId)
	if dept.ID == 29464267 {
		return true
	} else if dept.ID > 29464267 {
		return IsFsq(dept.Parentid)
	} else {
		return false
	}
}
