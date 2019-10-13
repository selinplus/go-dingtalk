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
		//if err := c.AddFunc("0 */10 * * * *", Sync); err != nil {//test定时任务，10分钟一次
		if err := c.AddFunc("@midnight", Sync); err != nil {
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
		logging.Info(fmt.Sprintf("get msg_flag err:%v", err))
		return
	}
	for _, msg := range msgs {
		tcmprJson := dingtalk.MseesageToDingding(msg)
		asyncsendReturn := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		//log.Printf("asyncsendReturn is :%v", asyncsendReturn)
		if asyncsendReturn != nil && asyncsendReturn.Errcode == 0 {
			if err := models.UpdateMsgFlag(msg.ID); err != nil {
				logging.Info(fmt.Sprintf("%v update msg_flag err:%v", msg.ID, err))
			}
		}
	}
}

//同步信息
func Sync() {
	logging.Info(fmt.Sprintf("DepartmentUserSync start..."))
	var (
		t  = time.Now().Format("2006-01-02") + " 00:00:00"
		sn = 30 //goroutine数量
		wt = 20 //发送递归请求的次数
	)
	count, e := dingtalk.OrgUserCount(wt)
	if e != nil {
		logging.Error(e)
	}
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 90)
		DepartmentUserSync(wt, sn)
		if userNum, err := models.CountUserSyncNum(t); err == nil && count-userNum < 5 {
			break
		}
	}
	logging.Info(fmt.Sprintf("DepartmentUserSync success"))
}

//同步一次部门和人员信息
func DepartmentUserSync(wt, syncNum int) {
	var (
		cntChan = make(chan int)
		wg      = &sync.WaitGroup{}
	)
	deptIds, err := dingtalk.SubDepartmentList(wt)
	if err != nil {
		//logging.Info(fmt.Sprintf("%v", err))
	}
	if deptIds != nil {
		var seg int
		deptIdsNum := len(deptIds)
		if deptIdsNum%8 == 0 {
			seg = deptIdsNum / 8
		} else {
			seg = (deptIdsNum / 8) + 1
		}
		deptIdChan := make(chan int, 100) //部门id
		for j := 0; j < 8; j++ {
			beg := j * seg
			if beg > deptIdsNum-1 {
				break
			}
			var end int
			if (j+1)*seg < deptIdsNum {
				end = (j + 1) * seg
			} else {
				end = deptIdsNum
			}
			segIds := deptIds[beg:end]
			go func() {
				for i, deptId := range segIds {
					deptIdChan <- deptId
					cntChan <- i
				}
			}()
		}
		go func() {
			var num int
			for range cntChan {
				num++
				if num == deptIdsNum {
					close(deptIdChan)
				}
			}
		}()
		for k := 0; k < syncNum; k++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for deptId := range deptIdChan {
					t := time.Now().Format("2006-01-02") + " 00:00:00"
					if flag := models.IsDeptExist(deptId, t); !flag {
						if department := dingtalk.DepartmentDetail(deptId, wt); department != nil {
							department.SyncTime = time.Now().Format("2006-01-02 15:04:05")
							if err := models.DepartmentSync(department); err != nil {
								log.Printf("DepartmentSync err:%v", err)
							}
						}
					}
					if userids := dingtalk.DepartmentUserIdsDetail(deptId, wt); userids != nil {
						cnt := len(userids)
						var pageNumTotal int
						if cnt%100 == 0 {
							pageNumTotal = cnt / 100
						} else {
							pageNumTotal = cnt/100 + 1
						}
						for pageNum := 0; pageNum < pageNumTotal; pageNum++ {
							userlist := dingtalk.DepartmentUserDetail(deptId, pageNum, wt)
							if userlist != nil {
								for _, user := range *userlist {
									if user.UserID == "" {
										continue
									}
									if flag := models.IsUserExist(user.UserID, t); !flag {
										if err := models.UserSync(&user); err != nil {
											log.Printf("UserSync err:%v", err)
										}
									}
								}
							}
						}
					}
				}
			}()
		}
		wg.Wait()
	}
}
