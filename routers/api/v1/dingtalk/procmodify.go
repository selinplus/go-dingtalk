package dingtalk

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/app"
	"github.com/selinplus/go-dingtalk/pkg/e"
	"net/http"
	"strconv"
	"time"
)

type DealForm struct {
	ID   uint
	Spyj string `json:"spyj"`
	Czr  string `json:"czr"`
}

//事件处理(退回&&通过)
func DealProc(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		flag = c.Query("flag")
		form DealForm
		node *models.Procnode
		t    string
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	t = time.Now().Format("2006-01-02 15:04:05")
	pm, err := models.GetProcMod(form.ID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_PROCMOD_FAIL, nil)
		return
	}
	pm.Spyj = form.Spyj
	pm.Czrq = t
	if err := models.UpdateProcMod(pm); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
		return
	}
	if flag == "yes" {
		if pm.Dm == "0" { //if 手工提报
			pm := models.Procmodify{
				ProcID: pm.ProcID,
				Dm:     pm.Dm,
				Tsr:    pm.Czr,
				Czr:    form.Czr,
			}
			if err := models.AddProcMod(&pm); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
				return
			}
			appG.Response(http.StatusOK, e.SUCCESS, nil)
			return
		}
		node, err = models.GetNextNode(pm.Dm, pm.Node)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		procMod := models.Procmodify{
			ProcID: pm.ProcID,
			Node:   node.Node,
			Dm:     pm.Dm,
			Tsr:    pm.Czr,
			Czr:    form.Czr,
		}
		for {
			if node.Flag == "0" {
				if node.Next == "-1" {
					procMod.Czr = procMod.Tsr
					procMod.Czrq = t
					procMod.Spyj = COMPLETE
				}
				if err := models.AddProcMod(&procMod); err != nil {
					appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
					return
				}
				appG.Response(http.StatusOK, e.SUCCESS, nil)
				return
			}
			if node.Flag == "1" {
				procMod.Spyj = AGREE
				procMod.Czrq = t
			}
			if err := models.AddProcMod(&procMod); err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
				return
			}
			node, err = models.GetNextNode(node.Dm, node.Node)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
				return
			}
		}
	}
	node, err = models.GetNode(pm.Dm, pm.Node)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	procMod := models.Procmodify{
		Czr: node.Rname,
	}
	//for {
	pmd := models.Procmodify{
		ProcID: pm.ProcID,
		Dm:     pm.Dm,
		Tsr:    procMod.Czr,
	}
	pmd.Node = node.Last
	if node.Node == "1" {
		p, err := models.GetProcDetail(pm.ProcID)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_PROC_FAIL, nil)
			return
		}
		pmd.Czr = p.Mobile
		if err = models.AddProcMod(&pmd); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	node, err = models.GetLastNode(node.Dm, node.Node)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	pmd.Czr = node.Rname
	//if node.Flag == "0" {
	if err := models.AddProcMod(&pmd); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
	return
	//}
	//if node.Flag == "1" {
	//	procMod.Czr = pmd.Czr
	//	pmd.Spyj = form.Spyj
	//	pmd.Czrq = t
	//	if err := models.AddProcMod(&pmd); err != nil {
	//		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_PROCMOD_FAIL, nil)
	//		return
	//	}
	//	if node.Last == "0" {
	//		//continue
	//	} else {
	//		node, err = models.GetLastNode(node.Dm, node.Node)
	//		if err != nil {
	//			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
	//			return
	//		}
	//	}
	//}
	//}
}

//事件处理流水记录查询
func GetProcModList(c *gin.Context) {
	appG := app.Gin{C: c}
	procid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	procMods, err := models.GetProcMods(uint(procid))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_PROCMOD_FAIL, nil)
		return
	}
	for _, procMod := range procMods {
		u, err := models.GetUserByMobile(procMod.Czr)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_GET_USERBYMOBILE_FAIL, nil)
			return
		}
		procMod.Czr = u.Name
	}
	appG.Response(http.StatusOK, e.SUCCESS, procMods)
}
