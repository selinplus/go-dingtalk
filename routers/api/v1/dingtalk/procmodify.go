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
	Flag string `json:"flag"`
}

//事件处理(退回&&通过)
func DealProc(c *gin.Context) {
	var (
		appG   = app.Gin{C: c}
		form   DealForm
		node   *models.Procnode
		t, Czr string
	)
	pm, err := models.GetProcMod(form.ID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if form.Flag == "yes" {
		node, err = models.GetNextNode(pm.Dm, pm.Node)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		if node.Flag == "1" {
			t = time.Now().Format("2006-01-02 15:04:05")
		}
		procMod := models.Procmodify{
			ProcID: pm.ProcID,
			Node:   node.Node,
			Tsr:    pm.Czr,
			Czr:    form.Czr,
			Czrq:   t,
		}
		if err := models.AddProcMod(&procMod); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		for {
			if node.Flag == "0" {
				break
			}
			if node.Flag == "1" {
				next, err := models.GetNextNode(node.Dm, node.Node)
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
	} else {
		node, err = models.GetLastNode(pm.Dm, pm.Node)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, nil)
			return
		}
		if node.Flag == "1" {
			t = time.Now().Format("2006-01-02 15:04:05")
		}
		if node.Last == "0" {
			p, err := models.GetProcDetail(pm.ProcID)
			if err != nil {
				appG.Response(http.StatusInternalServerError, e.ERROR_GET_PROC_FAIL, nil)
				return
			}
			Czr = p.Mobile
		} else {
			Czr = node.Rname
		}
		procMod := models.Procmodify{
			ProcID: pm.ProcID,
			Node:   node.Node,
			Tsr:    pm.Czr,
			Czr:    Czr,
			Czrq:   t,
		}
		if err := models.AddProcMod(&procMod); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
			return
		}
		for {
			if node.Last != "0" {
				if node.Flag == "0" {
					break
				}
				if node.Flag == "1" {
					last, err := models.GetLastNode(node.Dm, node.Node)
					if err != nil {
						appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
						return
					}
					procMod.Node = last.Node
					if last.Last == "0" {
						p, err := models.GetProcDetail(pm.ProcID)
						if err != nil {
							appG.Response(http.StatusInternalServerError, e.ERROR_GET_PROC_FAIL, nil)
							return
						}
						procMod.Czr = p.Mobile
					} else {
						procMod.Czr = last.Rname
					}
					if err := models.AddProcMod(&procMod); err != nil {
						appG.Response(http.StatusInternalServerError, e.ERROR_ADD_DEV_FAIL, nil)
						return
					}
				}
			}
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//事件处理流水记录查询
func GetProcModList(c *gin.Context) {
	appG := app.Gin{C: c}
	procid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	proMods, err := models.GetProcMods(uint(procid))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_PROCMOD_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, proMods)
}
