package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"log"
	"sync"
	"time"
)

var flag bool
var depidsNum int

func Setup() {
	// 定义一个cron运行器
	c := cron.New()
	// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送工作通知
	if err := c.AddFunc("*/30 * * * * *", MessageDingding); err != nil {
		logging.Info(fmt.Sprintf("Send MessageDingding failed：%v", err))
	}
	// 每天半夜同步一次部门和人员信息
	if err := c.AddFunc("@midnight", func() {
		DepartmentUserSync()
		for {
			if flag {
				time.Sleep(time.Second * 60)
				DepartmentUserSync()
				depNum, _ := models.CountDepartmentSyncNum()
				if depidsNum == depNum {
					flag = false
				}
			}
			break
		}
		logging.Info(fmt.Sprintf("DepartmentUserSync success"))
	}); err != nil {
		logging.Info(fmt.Sprintf("DepartmentUserSync failed：%v", err))
	}

	// 开始
	c.Start()
	defer c.Stop()
}

//遍历一遍发送标志为0的信息，通知钉钉发送工作通知
func MessageDingding() {
	msgs, _ := models.GetMsgFlag()
	for _, msg := range msgs {
		tcmprJson := dingtalk.MseesageToDingding(msg.Title, msg.Content, msg.ToID)
		asyncsendReturn := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendReturn != nil {
			if asyncsendReturn.Errcode == 0 {
				if err := models.UpdateMsgFlag(msg.ID); err != nil {
					logging.Info(fmt.Sprintf("UpdateMsgFlag err:%v", err))
				}
			}
		}
	}
}

//同步一次部门和人员信息
func DepartmentUserSync() {
	defer func() {
		if r := recover(); r != nil {
			flag = true
		}
	}()
	depIds, err := dingtalk.SubDepartmentList()
	if err != nil {
		logging.Info(fmt.Sprintf("%v", err))
		return
	}
	if depIds != nil {
		var seg int
		depidsLen := len(depIds)
		depidsNum = depidsLen //用于判断信息同步是否完成
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
						flag = true
					}
				}()
				for depId := range depIdChan {
					department := dingtalk.DepartmentDetail(depId)
					department.SyncTime = time.Now().Format("2006-01-02 15:04:05")
					if department.ID != 0 {
						if err := models.DepartmentSync(department); err != nil {
							log.Printf("DepartmentSync err:%v", err)
						}
					}
					userids := dingtalk.DepartmentUserIdsDetail(depId)
					cnt := len(userids)
					var pageNumTotal int
					if cnt%100 == 0 {
						pageNumTotal = cnt / 100
					} else {
						pageNumTotal = cnt/100 + 1
					}
					for pageNum := 0; pageNum < pageNumTotal; pageNum++ {
						userlist := dingtalk.DepartmentUserDetail(depId, pageNum)
						if err := models.UserSync(userlist); err != nil {
							log.Printf("UserSync err:%v", err)
						}
					}
				}
			}()
		}
		wg.Wait()
	}
}
