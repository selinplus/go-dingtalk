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

func Setup() {
	// 定义一个cron运行器
	c := cron.New()
	// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送工作通知
	c.AddFunc("*/30 * * * * *", MessageDingding)
	// 每天同步一次部门和人员信息
	c.AddFunc("@midnight", DepartmentUserSync)

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
				models.UpdateMsgFlag(msg.ID)
			}
		}
	}
}

//同步一次部门和人员信息
func DepartmentUserSync() {
	depIds, err := dingtalk.SubDepartmentList()
	if err != nil {
		logging.Info(fmt.Sprintf("%v", err))
		return
	}
	if depIds != nil {
		var seg int
		depidsLen := len(depIds)
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
				for depId := range depIdChan {
					department := dingtalk.DepartmentDetail(depId)
					department.SyncTime = time.Now().Format("2006-01-02 15:04:05")
					if department.ID != 0 {
						if err := models.DepartmentSync(department); err != nil {
							log.Println("DepartmentSync err:%v", err)
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
							log.Println("UserSync err:%v", err)
						}
					}
				}
			}()
		}
		wg.Wait()
		logging.Info(fmt.Sprintf("DepartmentUserSync success"))
	}
}
