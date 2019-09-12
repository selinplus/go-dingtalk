package dingtalk

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"time"
)

type ProcForm struct {
	ID     uint
	Dm     string `json:"dm"`
	Tbr    string `json:"tbr"`
	Mobile string `json:"mobile"`
	DevID  string `json:"devid"`
	Xq     string `json:"xq"`
	Zp     string `json:"zp"`
	Tbsj   string `json:"tbsj"`
	Czr    string `json:"czr"`
}

//提报事项保存&&提交
func AddProc(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		flag = c.Query("flage")
		form ProcForm
		t    string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if flag == "submit" {
		t = time.Now().Format("2006-01-02 15:04:05")
	}
	proc := models.Process{
		Dm:     form.Dm,
		Tbr:    form.Tbr,
		Mobile: form.Mobile,
		DevID:  form.DevID,
		Xq:     form.Xq,
		Zp:     form.Zp,
		Tbsj:   t,
	}
	if err := models.AddProc(&proc); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_SAVE_PROC_FAIL, nil)
		return
	}
	if t != "" {
		procMod := models.Procmodify{
			ProcID: proc.ID,
			Node:   "0",
			Tsr:    form.Mobile,
			Czr:    form.Mobile,
			Czrq:   t,
		}
		if err := models.AddProcMod(&procMod); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		next, err := models.GetNextNode(form.Dm, "0")
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		procMod.Node = next.Node
		procMod.Czr = form.Czr
		if err := models.AddProcMod(&procMod); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		for {
			if next.Flag == "0" {
				break
			}
			if next.Flag == "1" {
				next, err := models.GetNextNode(form.Dm, next.Node)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
					return
				}
				procMod.Node = next.Node
				procMod.Czr = next.Rname
				if err := models.AddProcMod(&procMod); err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
					return
				}
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//未提交事项修改&&提交
func UpdateProc(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		flag = c.Query("flag")
		form ProcForm
		t    string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if flag == "submit" {
		t = time.Now().Format("2006-01-02 15:04:05")
	}
	proc := models.Process{
		ID:     form.ID,
		Dm:     form.Dm,
		Tbr:    form.Tbr,
		Mobile: form.Mobile,
		DevID:  form.DevID,
		Xq:     form.Xq,
		Zp:     form.Zp,
		Tbsj:   t,
	}
	if err := models.UpdateProc(&proc); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_SAVE_PROC_FAIL, nil)
		return
	}
	if t != "" {
		procMod := models.Procmodify{
			ProcID: proc.ID,
			Node:   "0",
			Tsr:    form.Mobile,
			Czr:    form.Mobile,
			Czrq:   t,
		}
		if err := models.AddProcMod(&procMod); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		next, err := models.GetNextNode(form.Dm, "0")
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		procMod.Node = next.Node
		procMod.Czr = form.Czr
		if err := models.AddProcMod(&procMod); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		for {
			if next.Flag == "0" {
				break
			}
			if next.Flag == "1" {
				next, err := models.GetNextNode(form.Dm, next.Node)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
					return
				}
				procMod.Node = next.Node
				procMod.Czr = next.Rname
				if err := models.AddProcMod(&procMod); err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
					return
				}
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GetProcDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	procid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	proc, err := models.GetProcDetail(uint(procid))
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
		return
	}
	if proc.ID > 0 {
		node, err := models.GetNode(proc.Dm, proc.Node)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
			return
		}
		if node.Role == "" {
			if node.Node == "0" {
				proc.Zt = "已保存未提交"
			}
			if node.Next == "-1" {
				proc.Zt = "处理完成"
			}
		} else {
			proc.Zt = "已提交至" + node.Role
		}
		appG.Response(http.StatusOK, e.SUCCESS, proc)
	} else {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
	}
}

func GetProcTodoList(c *gin.Context) {
	appG := app.Gin{C: c}
	session := sessions.Default(c)
	userid := fmt.Sprintf("%v", session.Get("userid"))
	czr, uerr := models.GetUserByUserid(userid)
	if uerr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
		return
	}
	procList, err := models.GetProcTodoList(czr.Mobile)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
		return
	}
	if len(procList) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, procList)
	} else {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
	}
}

func GetProcDoneList(c *gin.Context) {
	appG := app.Gin{C: c}
	session := sessions.Default(c)
	userid := fmt.Sprintf("%v", session.Get("userid"))
	czr, uerr := models.GetUserByUserid(userid)
	if uerr != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
		return
	}
	procList, err := models.GetProcDoneList(czr.Mobile)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
		return
	}
	if len(procList) > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, procList)
	} else {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
	}
}
