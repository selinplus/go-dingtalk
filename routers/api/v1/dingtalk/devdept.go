package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"log"
	"net/http"
	"strings"
	"time"
)

type DevdeptForm struct {
	Jgdm   string `json:"jgdm"`
	Jgmc   string `json:"jgmc"`
	Sjjgdm string `json:"sjjgdm"`
	Gly    string `json:"gly"`
	Mobile string `json:"mobile"`
}

type GlyResp struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

//增加设备管理机构
func AddDevdept(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   DevdeptForm
		userid string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if len(form.Mobile) > 0 {
		u, err := models.GetUserByMobile(form.Mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, err.Error())
			return
		}
		userid = u.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}
	jgdm, err := models.GenDevdeptDmBySjjgdm(form.Sjjgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	devdept := models.Devdept{
		Jgdm:   jgdm,
		Jgmc:   form.Jgmc,
		Sjjgdm: form.Sjjgdm,
		Gly:    form.Gly,
		Lrr:    userid,
		Lrrq:   t,
		Xgr:    userid,
		Xgrq:   t,
	}
	if err := models.AddDevdept(&devdept); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEPARTMENT_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//修改设备管理机构
func UpdateDevdept(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   DevdeptForm
		userid string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if len(form.Mobile) > 0 {
		u, err := models.GetUserByMobile(form.Mobile)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err.Error())
			return
		}
		userid = u.UserID
	} else {
		token := c.GetHeader("Authorization")
		auth := c.Query("token")
		if len(auth) > 0 {
			token = auth
		}
		ts := strings.Split(token, ".")
		userid = ts[3]
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	devdept := models.Devdept{
		Jgdm:   form.Jgdm,
		Jgmc:   form.Jgmc,
		Sjjgdm: form.Sjjgdm,
		Gly:    form.Gly,
		Xgr:    userid,
		Xgrq:   t,
	}
	if err := models.UpdateDevdept(&devdept); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取设备管理机构信息
func GetDept(c *gin.Context) {
	appG := app.Gin{C: c}
	d, err := models.GetDevdept(c.Query("jgdm"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
}

//获取设备管理机构上级有管理员的机构信息
func GetParentDept(c *gin.Context) {
	appG := app.Gin{C: c}
	jgdm := c.Query("jgdm")
	sjjgdm := jgdm[:len(jgdm)-2]
	d, err := models.GetDevdept(sjjgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	if len(d.Gly) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, d)
		return
	}
	d, err = models.GetDevdept(sjjgdm[:len(sjjgdm)-2])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, d)
}

//获取设备管理机构列表(树结构)
func GetDevdeptTree(c *gin.Context) {
	appG := app.Gin{C: c}
	jgdm := c.Query("jgdm")
	bz := c.Query("bz")
	jgdms := strings.Split(jgdm, ",")
	dms := make([]string, 0)
	for _, dmSrc := range jgdms {
		flag := true
		for _, dm := range jgdms {
			if len(dmSrc) > len(dm) {
				if strings.Contains(dmSrc, dm) {
					flag = false
					break
				}
			}
		}
		if flag {
			dms = append(dms, dmSrc)
		}
	}
	data := make([]models.DevdeptTree, 0)
	for _, dm := range dms {
		tree, err := models.GetDevdeptTree(dm, bz)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
			return
		}
		if len(tree) > 0 {
			data = append(data, tree...)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取设备管理机构列表(bz:0-管理员不可选;1-管理员可选)
func GetDevdeptGlyList(c *gin.Context) {
	var (
		appG    = app.Gin{C: c}
		flag    bool
		glyName string
	)
	bz := c.Query("bz")
	parentDt, err := models.GetDevdept(c.Query("jgdm"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	var dts []interface{}
	if parentDt.Gly != "" {
		gly, err := models.GetUserByUserid(parentDt.Gly)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err.Error())
			return
		}
		glyName = gly.Name
	} else {
		glyName = ""
	}
	if bz == "0" {
		flag = false
	} else {
		flag = true
	}
	data := map[string]interface{}{
		"jgdm":        parentDt.Jgdm,
		"sjjgdm":      parentDt.Sjjgdm,
		"jgmc":        parentDt.Jgmc,
		"gly":         glyName,
		"children":    dts,
		"scopedSlots": map[string]string{"title": "custom"},
		"disabled":    flag,
	}
	departments, err := models.GetDevdeptBySjjgdm(parentDt.Jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	if len(departments) > 0 {
		for _, department := range departments {
			var ds []interface{}
			leaf := models.IsLeafDevdept(department.Jgdm)
			if !leaf {
				departs, err := models.GetDevdeptBySjjgdm(department.Jgdm)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, department.Jgdm)
					return
				}
				if len(departs) > 0 {
					for _, depart := range departs {
						if bz == "0" {
							if depart.Gly != "" {
								flag = true
							} else {
								flag = false
							}
						}
						if bz == "1" {
							if depart.Gly != "" && department.Gly == "" {
								flag = false
							} else {
								flag = true
							}
						}
						if depart.Gly != "" {
							gly, err := models.GetUserByUserid(depart.Gly)
							if err != nil {
								appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err.Error())
								return
							}
							glyName = gly.Name
						} else {
							glyName = ""
						}
						d := map[string]interface{}{
							"jgdm":        depart.Jgdm,
							"sjjgdm":      depart.Sjjgdm,
							"jgmc":        depart.Jgmc,
							"gly":         glyName,
							"scopedSlots": map[string]string{"title": "custom"},
							"children":    nil,
							"disabled":    flag,
						}
						ds = append(ds, d)
					}
				}
			}
			if bz == "0" {
				if department.Gly != "" {
					flag = true
				} else {
					flag = false
				}
			}
			if bz == "1" {
				if department.Gly != "" {
					flag = false
				} else {
					flag = true
				}
			}
			if department.Gly != "" {
				gly, err := models.GetUserByUserid(department.Gly)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err.Error())
					return
				}
				glyName = gly.Name
			} else {
				glyName = ""
			}
			dt := map[string]interface{}{
				"jgdm":        department.Jgdm,
				"sjjgdm":      department.Sjjgdm,
				"jgmc":        department.Jgmc,
				"gly":         glyName,
				"scopedSlots": map[string]string{"title": "custom"},
				"children":    ds,
				"disabled":    flag,
			}
			dts = append(dts, dt)
		}
		data["children"] = dts
		appG.Response(http.StatusOK, e.SUCCESS, data)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, data)
	}
}

//获取设备管理机构列表(循环遍历)
func GetDevdeptBySjjgdm(c *gin.Context) {
	var appG = app.Gin{C: c}
	jgdm := c.Query("jgdm")
	parentDt, err := models.GetDevdept(jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	var dts []interface{}
	data := map[string]interface{}{
		"key":      parentDt.Jgdm,
		"value":    parentDt.Jgdm,
		"title":    parentDt.Jgmc,
		"children": dts,
	}
	departments, err := models.GetDevdeptBySjjgdm(jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	if len(departments) > 0 {
		for _, department := range departments {
			leaf := models.IsLeafDevdept(department.Jgdm)
			dt := map[string]interface{}{
				"key":    department.Jgdm,
				"value":  department.Jgdm,
				"title":  department.Jgmc,
				"isLeaf": leaf,
			}
			dts = append(dts, dt)
		}
		data["children"] = dts
		appG.Response(http.StatusOK, e.SUCCESS, data)
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, data)
	}
}

//获取设备管理机构及人员列表(Epp循环遍历)
func GetDevdeptEppTree(c *gin.Context) {
	var appG = app.Gin{C: c}
	jgdm := c.Query("jgdm")
	var data []interface{}
	departments, err := models.GetDevdeptBySjjgdm(jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	if len(departments) > 0 {
		for _, department := range departments {
			if department.Gly == "" {
				dt := map[string]interface{}{
					"dm":     department.Jgdm,
					"mc":     department.Jgmc,
					"isDept": true,
				}
				data = append(data, dt)
			}
		}
	}
	dus, err := models.GetDevuser(jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEVUSER_FAIL, err.Error())
		return
	}
	if len(dus) > 0 {
		for _, du := range dus {
			user, err := models.GetUserByUserid(du.Syr)
			if err != nil {
				log.Println(err)
			}
			u := map[string]interface{}{
				"dm":     user.UserID,
				"mc":     user.Name,
				"isDept": false,
			}
			data = append(data, u)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//删除设备管理机构
func DeleteDevdept(c *gin.Context) {
	appG := app.Gin{C: c}
	jgdm := c.Query("jgdm")
	if models.IsSjjg(jgdm) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DEVDETP_IS_PARENT, nil)
		return
	}
	if models.IsDevdeptUserExist(jgdm) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DEVDETP_NOT_NULL, nil)
		return
	}
	if models.IsDevdeptGylExist(jgdm) {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DEVDETPGYL_NOT_NULL, nil)
		return
	}
	devs, err := models.GetDevinfosByJgdm(jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USERDEV_FAIL, err.Error())
		return
	}
	if len(devs) > 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USERDEV_FAIL, nil)
		return
	}
	if err := models.DeleteDevdept(jgdm); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_DEPARTMENT_FAIL, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除当前机构管理员
func DelDevdeptGly(c *gin.Context) {
	appG := app.Gin{C: c}
	jgdm := c.Query("jgdm")
	devdept := map[string]interface{}{
		"jgdm": jgdm,
		"gly":  "",
		"xgrq": time.Now().Format("2006-01-02 15:04:05"),
	}
	devs, err := models.GetDevinfosByJgdm(jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USERDEV_FAIL, nil)
		return
	}
	if len(devs) > 0 {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USERDEV_FAIL, nil)
		return
	}
	if err := models.DelDevdeptGly(devdept); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//获取当前机构管理员信息
func GetDevdeptGly(c *gin.Context) {
	appG := app.Gin{C: c}
	jgdm := c.Query("jgdm")
	ddept, err := models.GetDevdept(jgdm)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	resps := make([]*GlyResp, 0)
	if ddept.Gly == "" {
		appG.Response(http.StatusOK, e.SUCCESS, resps)
		return
	}
	user, err := models.GetUserByUserid(ddept.Gly)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, err.Error())
		return
	}
	resp := &GlyResp{
		UserID: user.UserID,
		Name:   user.Name,
		Mobile: user.Mobile,
	}
	resps = append(resps, resp)
	appG.Response(http.StatusOK, e.SUCCESS, resps)
}

//获取当前用户为机构管理员的所有机构列表
func GetDevGly(c *gin.Context) {
	appG := app.Gin{C: c}
	gly, err := models.GetUserByMobile(c.Query("gly"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, err.Error())
		return
	}
	depts, err := models.GetDevdeptsHasGlyByUserid(gly.UserID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_DEPARTMENT_FAIL, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["lists"] = depts
	data["total"] = len(depts)
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
