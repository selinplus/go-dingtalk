package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SAVENOTSUB = "已保存未提交"
	COMPLETE   = "处理完成"
	BACKTOBEG  = "退回发起人"
	SUBMIT     = "发起事件"
	SUBAGAIN   = "重新上报"
	AGREE      = "同意"
	INVALID    = "作废"
)

type ProcForm struct {
	ID       uint
	Dm       string `json:"dm"`
	Tbr      string `json:"tbr"`
	Tsr      string `json:"tsr"`
	Mobile   string `json:"mobile"`
	DevID    string `json:"devid"`
	Title    string `json:"title"`
	Xq       string `json:"xq"`
	Bcms     string `json:"bcms"`
	Zp       string `json:"zp"`
	Czr      string `json:"czr"`
	Modifyid uint   `json:"modifyid"`
	SyrName  string `json:"syr_name"`
	Syr      string `json:"syr"`
	Cfwz     string `json:"cfwz"`
	Spyj     string `json:"spyj"`
}

//提报事项保存&&提交
func AddProc(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		flag   = c.Query("flag")
		form   ProcForm
		tsr, t string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	sj := time.Now().Format("2006-01-02")
	proc := models.Process{
		Dm:      form.Dm,
		Tbr:     form.Tbr,
		Mobile:  form.Mobile,
		DevID:   form.DevID,
		Title:   form.Title,
		Xq:      form.Xq,
		Zp:      form.Zp,
		Tbsj:    sj,
		SyrName: form.SyrName,
		Syr:     form.Syr,
		Cfwz:    form.Cfwz,
	}
	b, procid, err := models.IsProcExist(proc)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if !b {
		if err := models.AddProc(&proc); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_SAVE_PROC_FAIL, nil)
			return
		}
		procid = proc.ID
	}
	if flag == "submit" {
		t = time.Now().Format("2006-01-02 15:04:05")
	}
	if t != "" {
		if form.Dm == "0" { //if 手工提报
			pm := models.Procmodify{
				ProcID:     procid,
				Node:       "0",
				Dm:         form.Dm,
				Tsr:        form.Mobile,
				Czr:        form.Mobile,
				Spyj:       SUBMIT,
				Czrq:       t,
				FlagNotice: 1,
			}
			if err := models.AddProcMod(&pm); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
				return
			}
			pmnext := models.Procmodify{
				ProcID: procid,
				Node:   "1",
				Dm:     form.Dm,
				Tsr:    form.Mobile,
				Czr:    form.Czr,
			}
			if err := models.AddProcMod(&pmnext); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
				return
			}
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}
		procMod := models.Procmodify{
			ProcID:     procid,
			Node:       "0",
			Dm:         form.Dm,
			Tsr:        form.Mobile,
			Czr:        form.Mobile,
			Spyj:       SUBMIT,
			Czrq:       t,
			FlagNotice: 1,
		}
		if err := models.AddProcMod(&procMod); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
			return
		}
		next, err := models.GetNextNode(form.Dm, "0")
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		for {
			if next.Node == "1" {
				tsr = procMod.Czr
			} else {
				last, err := models.GetLastNode(form.Dm, next.Node)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR, nil)
					return
				}
				tsr = last.Rname
			}
			pm := models.Procmodify{
				ProcID: procMod.ProcID,
				Node:   next.Node,
				Dm:     form.Dm,
				Tsr:    tsr,
				Czr:    next.Rname,
			}
			if next.Flag == "0" {
				if err := models.AddProcMod(&pm); err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
					return
				}
				break
			}
			if next.Flag == "1" {
				pm.Spyj = AGREE
				pm.Czrq = t
				pm.FlagNotice = 1
				if err := models.AddProcMod(&pm); err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
					return
				}
				procMod.Czr = pm.Czr
				next, err = models.GetNextNode(form.Dm, next.Node)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR, nil)
					return
				}
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//(未)提交事项修改&&提交
func UpdateProc(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		flag   = c.Query("flag")
		form   ProcForm
		tsr, t string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	sj := time.Now().Format("2006-01-02")
	proc := models.Process{
		ID:     form.ID,
		Dm:     form.Dm,
		Tbr:    form.Tbr,
		Mobile: form.Mobile,
		DevID:  form.DevID,
		Title:  form.Title,
		Xq:     form.Xq,
		Zp:     form.Zp,
		Tbsj:   sj,
	}
	if err := models.UpdateProc(&proc); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_SAVE_PROC_FAIL, nil)
		return
	}
	if flag == "submit" {
		t = time.Now().Format("2006-01-02 15:04:05")
	}
	if t != "" {
		if form.Dm == "0" { //if 手工提报
			pm := models.Procmodify{
				ProcID:     proc.ID,
				Node:       "0",
				Dm:         form.Dm,
				Tsr:        form.Mobile,
				Czr:        form.Mobile,
				Spyj:       SUBMIT,
				Czrq:       t,
				FlagNotice: 1,
			}
			if err := models.AddProcMod(&pm); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
				return
			}
			pmnext := models.Procmodify{
				ProcID: proc.ID,
				Node:   "1",
				Dm:     form.Dm,
				Tsr:    form.Mobile,
				Czr:    form.Czr,
			}
			if err := models.AddProcMod(&pmnext); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
				return
			}
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}
		procMod := models.Procmodify{
			ProcID:     proc.ID,
			Node:       "0",
			Dm:         form.Dm,
			Czr:        form.Mobile,
			Czrq:       t,
			FlagNotice: 1,
		}
		if form.Modifyid > 0 {
			procMod.ID = form.Modifyid
			procMod.Spyj = SUBAGAIN
			if err := models.UpdateProcMod(&procMod); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
				return
			}
		} else {
			procMod.Tsr = form.Mobile
			procMod.Spyj = SUBMIT
			if err := models.AddProcMod(&procMod); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
				return
			}
		}
		next, err := models.GetNextNode(form.Dm, "0")
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		for {
			if next.Node == "1" {
				tsr = procMod.Czr
			} else {
				last, err := models.GetLastNode(form.Dm, next.Node)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR, nil)
					return
				}
				tsr = last.Rname
			}
			pm := models.Procmodify{
				ProcID: procMod.ProcID,
				Node:   next.Node,
				Dm:     form.Dm,
				Tsr:    tsr,
				Czr:    next.Rname,
			}
			if next.Flag == "0" {
				pm.Spyj = form.Spyj
				if err := models.AddProcMod(&pm); err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
					return
				}
				break
			}
			if next.Flag == "1" {
				pm.Spyj = AGREE
				pm.Czrq = t
				pm.FlagNotice = 1
				if err := models.AddProcMod(&pm); err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
					return
				}
				next, err = models.GetNextNode(form.Dm, next.Node)
				if err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR, nil)
					return
				}
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//作废&&删除提报事项
func DeleteProc(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	procid := uint(id)
	proc, err := models.GetProcDetail(procid)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
		return
	}
	if proc.Node == "" {
		if err := models.DeleteProc(procid); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	p := models.Process{
		ID:   procid,
		Zfbz: "1",
	}
	if err := models.UpdateProc(&p); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_SAVE_PROC_FAIL, nil)
		return
	}
	procMod := models.Procmodify{
		ID:         proc.Modifyid,
		Spyj:       INVALID,
		FlagNotice: 1,
		Czrq:       time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := models.UpdateProcMod(&procMod); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//发起补充描述
func ProcBcms(c *gin.Context) {
	appG := app.Gin{C: c}
	var form ProcForm
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	proc := models.Process{
		ID:   form.ID,
		Bcms: form.Bcms,
	}
	if err := models.UpdateProc(&proc); err != nil {
		appG.Response(http.StatusOK, e.ERROR_SAVE_PROC_FAIL, nil)
		return
	}
	procTag := models.ProcessTag{ProcID: form.ID}
	if err := models.AddProcessBcms(&procTag); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//补充描述提交
func UpdateProcBcms(c *gin.Context) {
	appG := app.Gin{C: c}
	var form ProcForm
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	proc := models.Process{
		ID:   form.ID,
		Bcms: form.Bcms,
	}
	if err := models.UpdateProc(&proc); err != nil {
		appG.Response(http.StatusOK, e.ERROR_SAVE_PROC_FAIL, nil)
		return
	}
	if err := models.UpdateProcessModFlag(form.Modifyid); err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//查询提报事项详情
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
	if proc.Syr != "" {
		syr, err := models.GetUserByMobile(proc.Syr)
		if err != nil {
			appG.Response(http.StatusOK, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		proc.Syr = syr.Name
	}
	if proc.Dm == "0" {
		if proc.Czr == "" {
			proc.Zt = SAVENOTSUB
			appG.Response(http.StatusOK, e.SUCCESS, proc)
			return
		}
		if models.IsProcManualDone(proc.ID) {
			proc.Zt = "已提交至" + proc.Czr + "处理"
		} else {
			proc.Zt = "已由" + proc.Czr + "办结"
		}
		appG.Response(http.StatusOK, e.SUCCESS, proc)
		return
	}
	if proc.Node == "" {
		proc.Zt = SAVENOTSUB
		appG.Response(http.StatusOK, e.SUCCESS, proc)
		return
	}
	if proc.Node == "0" {
		proc.Zt = BACKTOBEG
		appG.Response(http.StatusOK, e.SUCCESS, proc)
		return
	}
	node, err := models.GetNode(proc.Dm, proc.Node)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
		return
	}
	if node.Next == "-1" {
		proc.Zt = COMPLETE
	} else {
		proc.Zt = "已提交至" + node.Role
	}
	appG.Response(http.StatusOK, e.SUCCESS, proc)
}

//获取待办列表
func GetProcTodoList(c *gin.Context) {
	var data []interface{}
	appG := app.Gin{C: c}
	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]
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
		for _, p := range procList {
			if p.Dm == "0" {
				if models.IsProcManualDone(p.ID) {
					p.Zt = "已提交至" + p.Czr + "处理"
				} else {
					p.Zt = "已由" + p.Czr + "办结"
				}
			} else if p.Node == "0" {
				p.Zt = BACKTOBEG
			} else {
				node, err := models.GetNode(p.Dm, p.Node)
				if err != nil {
					appG.Response(http.StatusOK, e.ERROR, nil)
					return
				}
				if node.Next == "-1" {
					p.Zt = COMPLETE
				} else {
					p.Zt = "已提交至" + node.Role
				}
			}
			data = append(data, p)
		}
		appG.Response(http.StatusOK, e.SUCCESS, data)
		return
	}
	//已保存未提交
	psList, err := models.GetProcSaveList(czr.Mobile)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
		return
	}
	for _, proc := range psList {
		proc.Zt = "已保存未提交"
	}
	appG.Response(http.StatusOK, e.SUCCESS, psList)
}

//获取已办列表(全部)
func GetProcDoneList(c *gin.Context) {
	var data []interface{}
	appG := app.Gin{C: c}
	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]
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
		for _, proc := range procList {
			p, err := models.GetProcDetail(proc.ID)
			if err != nil {
				appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
				return
			}
			if p.Dm == "0" {
				if models.IsProcManualDone(p.ID) {
					p.Zt = "已提交至" + p.Czr + "处理"
				} else {
					p.Zt = "已由" + p.Czr + "办结"
				}
			} else if p.Node == "0" {
				p.Zt = BACKTOBEG
			} else {
				node, err := models.GetNode(p.Dm, p.Node)
				if err != nil {
					appG.Response(http.StatusOK, e.ERROR, nil)
					return
				}
				if node.Next == "-1" {
					p.Zt = COMPLETE
				} else {
					p.Zt = "已提交至" + node.Role
				}
			}
			data = append(data, p)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取已办列表(已办结)
func GetProcDoneListEnd(c *gin.Context) {
	var data []interface{}
	appG := app.Gin{C: c}
	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]
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
		for _, proc := range procList {
			p, err := models.GetProcDetail(proc.ID)
			if err != nil {
				appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
				return
			}
			if p.Dm == "0" {
				if models.IsProcManualDone(p.ID) {
					continue
				} else {
					p.Zt = "已由" + p.Czr + "办结"
				}
			} else if p.Node == "0" {
				continue
			} else {
				node, err := models.GetNode(p.Dm, p.Node)
				if err != nil {
					appG.Response(http.StatusOK, e.ERROR, nil)
					return
				}
				if node.Next == "-1" {
					p.Zt = COMPLETE
				} else {
					continue
				}
			}
			data = append(data, p)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//获取已办列表(未办结)
func GetProcDoneListDoing(c *gin.Context) {
	var data []interface{}
	appG := app.Gin{C: c}
	token := c.GetHeader("Authorization")
	auth := c.Query("token")
	if len(auth) > 0 {
		token = auth
	}
	ts := strings.Split(token, ".")
	userid := ts[3]
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
		for _, proc := range procList {
			p, err := models.GetProcDetail(proc.ID)
			if err != nil {
				appG.Response(http.StatusOK, e.ERROR_GET_PROC_FAIL, nil)
				return
			}
			if p.Dm == "0" {
				if models.IsProcManualDone(p.ID) {
					p.Zt = "已提交至" + p.Czr + "处理"
				} else {
					continue
				}
			} else if p.Node == "0" {
				p.Zt = BACKTOBEG
			} else {
				node, err := models.GetNode(p.Dm, p.Node)
				if err != nil {
					appG.Response(http.StatusOK, e.ERROR, nil)
					return
				}
				if node.Next == "-1" {
					continue
				} else {
					p.Zt = "已提交至" + node.Role
				}
			}
			data = append(data, p)
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
