package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"log"
	"os"
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
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送记事本消息
		if err := c.AddFunc("*/30 * * * * *", NoteMessageDingding); err != nil {
			logging.Info(fmt.Sprintf("Send Process NoteMessageDingding failed：%v", err))
		}
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送提报事项待办工作通知消息
		if err := c.AddFunc("*/30 * * * * *", ProcessMessageDingding); err != nil {
			logging.Info(fmt.Sprintf("Send Process ProcessMessageDingding failed：%v", err))
		}
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送提报事项补充描述工作通知
		if err := c.AddFunc("*/30 * * * * *", ProcessBcmsMessageDingding); err != nil {
			logging.Info(fmt.Sprintf("Send ProcessBcms ProcessBcmsMessageDingding failed：%v", err))
		}
		// 每天半夜同步一次部门和人员信息
		if err := c.AddFunc("@midnight", Sync); err != nil {
			logging.Info(fmt.Sprintf("DepartmentUserSync failed：%v", err))
		}
		// 每个月执行一遍网盘回收站清理
		if err := c.AddFunc("@monthly", CleanUpNetdiskFile); err != nil {
			logging.Info(fmt.Sprintf("CleanUp NetdiskFile failed：%v", err))
		}
		// 开始
		c.Run()
	}()
}

//清理网盘回收站30天以上文件
func CleanUpNetdiskFile() {
	files, err := models.GetTrashFiles()
	if err != nil {
		logging.Error(fmt.Sprintf("get trash files err:%v", err))
		return
	}
	for _, fileInfo := range files {
		dirUrl := upload.GetImageFullPath() + fileInfo.FileUrl
		if err = os.Remove(dirUrl); err != nil {
			logging.Error(fmt.Sprintf("delete trash files err:%v", err))
		} else { //delete table column
			if err = models.DeleteNetdiskFile(fileInfo.ID); err != nil {
				logging.Error(fmt.Sprintf("delete column err:%v", err))
			}
		}
	}
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

//遍历一遍发送标志为0的信息，通知钉钉发送记事本消息
func NoteMessageDingding() {
	notes, err := models.GetNoteFlag()
	if err != nil {
		logging.Info(fmt.Sprintf("get note_flag err:%v", err))
		return
	}
	for _, note := range notes {
		tcmprJson := dingtalk.NoteMseesageToDingding(note)
		asyncsendReturn := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendReturn != nil && asyncsendReturn.Errcode == 0 {
			if err := models.UpdateNoteFlag(note.ID); err != nil {
				logging.Info(fmt.Sprintf("%v update note_flag err:%v", note.ID, err))
			}
		}
	}
}

//遍历一遍发送标志为0的信息，通知钉钉发送提报事项待办工作通知消息
func ProcessMessageDingding() {
	procs, err := models.GetProcessFlag()
	if err != nil {
		logging.Info(fmt.Sprintf("get process_flag err:%v", err))
		return
	}
	for _, proc := range procs {
		p, err := models.GetProcDetail(proc.ProcID)
		if err != nil {
			logging.Info(fmt.Sprintf("get process detail [id:%v] err:%v", proc.ID, err))
		}
		tcmprJson := dingtalk.ProcessMseesageToDingding(p, proc.Czr)
		asyncsendReturn := dingtalk.EappMessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendReturn != nil && asyncsendReturn.Errcode == 0 {
			if err := models.UpdateProcessFlag(proc.ID); err != nil {
				logging.Info(fmt.Sprintf("%v update process_flag err:%v", proc.ID, err))
			}
		}
	}
}

//遍历一遍发送标志为0的信息，通知钉钉发送提报事项补充描述工作通知
func ProcessBcmsMessageDingding() {
	procs, err := models.GetProcessBcmsFlag()
	if err != nil {
		logging.Info(fmt.Sprintf("get process_flag err:%v", err))
		return
	}
	for _, proc := range procs {
		p, err := models.GetProcDetail(proc.ProcID)
		if err != nil {
			logging.Info(fmt.Sprintf("get process detail [id:%v] err:%v", proc.ProcID, err))
		}
		tcmprJson := dingtalk.ProcessBcmsMseesageToDingding(p)
		asyncsendReturn := dingtalk.EappMessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendReturn != nil && asyncsendReturn.Errcode == 0 {
			if err := models.UpdateProcessBcmsFlag(proc.ID); err != nil {
				logging.Info(fmt.Sprintf("%v update process_flag err:%v", proc.ID, err))
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
