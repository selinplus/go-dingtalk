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
	go func() {
		log.Println("Cron starting...")
		// 定义一个cron运行器
		c := cron.New()
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送工作通知
		if err := c.AddFunc("*/30 * * * * *", MessageDingding); err != nil {
			logging.Info(fmt.Sprintf("Send MessageDingding failed：%v", err))
		}
		// 每天半夜同步一次部门和人员信息
		//if err := c.AddFunc("0 */10 * * * *", func() { //test定时任务，10分钟一次
		if err := c.AddFunc("@midnight", func() {
			logging.Info(fmt.Sprintf("DepartmentUserSync start..."))
			wt := 20      //发生网页劫持后，发送递归请求的次数
			syncNum := 30 //goroutine数量
			for i := 0; i < 10; i++ {
				time.Sleep(time.Second * 90)
				useridsNum, depidsNum := DepartmentUserSync(wt, syncNum)
				if useridsNum > 0 && depidsNum > 0 {
					userNum, _ := models.CountUserSyncNum()
					depNum, _ := models.CountDepartmentSyncNum()
					if userNum == useridsNum && depNum == depidsNum {
						goto Loop
					}
				}
			}
		Loop:
			logging.Info(fmt.Sprintf("DepartmentUserSync stopped"))
		}); err != nil {
			logging.Info(fmt.Sprintf("DepartmentUserSync failed：%v", err))
		}
		// 开始
		c.Run()
	}()
}

//遍历一遍发送标志为0的信息，通知钉钉发送工作通知
func MessageDingding() {
	msgs, err := models.GetMsgFlag()
	if err != nil {
		return
	}
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
func DepartmentUserSync(wt, syncNum int) (int, int) {
	var (
		userIdsNum = 0
		depidsLen  = 0
	)
	depIds, err := dingtalk.SubDepartmentList(wt)
	if err != nil {
		logging.Info(fmt.Sprintf("%v", err))
		return userIdsNum, depidsLen
	}
	if depIds != nil {
		var seg int
		depidsLen = len(depIds)
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
		wg := &sync.WaitGroup{}
		wg.Add(syncNum)
		for k := 0; k < syncNum; k++ {
			go func() {
				for depId := range depIdChan {
					department := dingtalk.DepartmentDetail(depId, wt)
					department.SyncTime = time.Now().Format("2006-01-02 15:04:05")
					if department.ID != 0 {
						if err := models.DepartmentSync(department); err != nil {
							log.Printf("DepartmentSync err:%v", err)
						}
					}
					if userids := dingtalk.DepartmentUserIdsDetail(depId, wt); userids != nil {
						cnt := len(userids)
						userIdsNum += cnt
						var pageNumTotal int
						if cnt%100 == 0 {
							pageNumTotal = cnt / 100
						} else {
							pageNumTotal = cnt/100 + 1
						}
						for pageNum := 0; pageNum < pageNumTotal; pageNum++ {
							if userlist := dingtalk.DepartmentUserDetail(depId, pageNum, wt); userlist != nil {
								if err := models.UserSync(userlist); err != nil {
									log.Printf("UserSync err:%v", err)
								}
							}
						}
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
	return userIdsNum, depidsLen
}
