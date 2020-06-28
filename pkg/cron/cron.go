package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/export"
	"github.com/selinplus/go-dingtalk/pkg/file"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"github.com/selinplus/go-dingtalk/pkg/ydksrv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// dmz区 crontab
func DmzSetup() {
	go func() {
		log.Println("dmz crontab starting...")
		// 定义一个cron运行器
		c := cron.New()
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送工作通知
		if err := c.AddFunc("*/30 * * * * *", MessageDingding); err != nil {
			log.Printf("Send MessageDingding failed：%v", err)
		}
		// 遍历一遍发送标志为0的交回设备信息，通知钉钉发送工作通知管理员
		if err := c.AddFunc("*/30 * * * * *", DeviceDingding); err != nil {
			log.Printf("Send DeviceDingding failed：%v", err)
		}
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送记事本消息
		if err := c.AddFunc("*/30 * * * * *", NoteMessageDingding); err != nil {
			log.Printf("Send Process NoteMessageDingding failed：%v", err)
		}
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送值班通知消息
		if err := c.AddFunc("*/30 * * * * *", OndutyMessageDingding); err != nil {
			log.Printf("Send Process OndutyMessageDingding failed：%v", err)
		}
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送提报事项待办工作通知消息
		if err := c.AddFunc("*/30 * * * * *", ProcessMessageDingding); err != nil {
			log.Printf("Send Process ProcessMessageDingding failed：%v", err)
		}
		// 每30秒遍历一遍发送标志为0的信息，通知钉钉发送提报事项补充描述工作通知
		if err := c.AddFunc("*/30 * * * * *", ProcessBcmsMessageDingding); err != nil {
			log.Printf("Send ProcessBcms ProcessBcmsMessageDingding failed：%v", err)
		}
		// 每30秒遍历一遍 ydks 消息，通知钉钉发送待办任务
		if err := c.AddFunc("*/30 * * * * *", ydksrv.Ydksworkrecord); err != nil {
			log.Printf("Send Ydksworkrecord failed：%v", err)
		}
		// 每天半夜同步一次部门和人员信息
		if err := c.AddFunc("@midnight", Sync); err != nil {
			log.Printf("DepartmentUserSync failed：%v", err)
		}
		// 每个月执行一遍网盘回收站清理
		if err := c.AddFunc("@monthly", CleanUpNetdiskFile); err != nil {
			log.Printf("CleanUp NetdiskFile failed：%v", err)
		}

		// 开始
		c.Run()
	}()
}

// app区 crontab
func AppSetup() {
	go func() {
		log.Println("app crontab starting...")
		// 定义一个cron运行器
		c := cron.New()
		// 每天半夜将前一天 ydks 数据写入文件
		if err := c.AddFunc("@midnight", func() {
			date := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
			ydksrv.WriteIntoFile(date)
		}); err != nil {
			log.Printf("WriteIntoFile crontab failed：%v", err)
		}
		// 每天1点清理超过1天的导出记录
		if err := c.AddFunc("0 0 1 * * *", CleanUpExportFiles); err != nil {
			log.Printf("WriteIntoFile crontab failed：%v", err)
		}
		// 每天2点清理同步时间超过7天的部门用户记录
		if err := c.AddFunc("0 0 2 * * *", CleanUpDepartmentUser); err != nil {
			log.Printf("CleanUp Department&User crontab failed：%v", err)
		}
		// 每个月清理一遍超过30天的以地控税文件
		if err := c.AddFunc("@monthly", ydksrv.CleanUpYdksFiles); err != nil {
			log.Printf("CleanUp YdksFiles failed：%v", err)
		}
		// 开始
		c.Run()
	}()
}

//清理超过一天的导出文件
func CleanUpExportFiles() {
	dirpath := export.GetExcelFullPath()
	files, err := file.FindFilesOlderThanDate(dirpath, 1)
	errNotExist := "open : The system cannot find the file specified."
	if err != nil && err.Error() != errNotExist {
		log.Println("CleanUp ExportFiles err:", err)
		return
	}
	for _, fileInfo := range files {
		if strings.Contains(fileInfo.Name(), "device.xlsx") {
			continue
		}
		err = os.Remove(dirpath + fileInfo.Name())
		if err != nil {
			log.Println("CleanUp ExportFiles err:", err)
		}
	}
}

//清理同步时间超过7天的部门用户记录
func CleanUpDepartmentUser() {
	if err := models.CleanUpDepartment(); err != nil {
		logging.Error(fmt.Sprintf("CleanUp Department err:%v", err))
	}
	if err := models.CleanUpUser(); err != nil {
		logging.Error(fmt.Sprintf("CleanUp User err:%v", err))
	}
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
		logging.Error(fmt.Sprintf("get msg_flag err:%v", err))
		return
	}
	for _, msg := range msgs {
		tcmprJson := dingtalk.MseesageToDingding(msg)
		asyncsendResponse := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		//log.Printf("asyncsendResponse is :%v", asyncsendResponse)
		if asyncsendResponse != nil && asyncsendResponse.ErrCode == 0 {
			if err := models.UpdateMsgFlag(msg.ID); err != nil {
				logging.Error(fmt.Sprintf("%v update msg_flag err:%v", msg.ID, err))
			}
		}
	}
}

//遍历一遍发送标志为0的交回设备信息，通知钉钉发送工作通知管理员
func DeviceDingding() {
	todos, err := models.GetDevFlag()
	if err != nil {
		logging.Error(fmt.Sprintf("get dev_flag err:%v", err))
		return
	}
	for _, todo := range todos {
		var tcmprJson string
		if todo.Czlx == "8" { //交回,发送link
			tcmprJson = dingtalk.DeviceDingding(todo.DevID, todo.Gly, strconv.Itoa(todo.Done))
		}
		if todo.Czlx == "10" { //上交,发送text
			dept, err := models.GetDevdept(todo.SrcJgdm)
			if err != nil {
				logging.Error(fmt.Sprintf("%v GetDevdept err:%v", todo.ID, err))
			}
			tcmprJson = dingtalk.UpDeviceDingding(todo.Num, dept.Jgmc, todo.Gly)
		}
		asyncsendResponse := dingtalk.EappMessageCorpconversationAsyncsend(tcmprJson)
		//log.Printf("asyncsendResponse is :%v", asyncsendResponse)
		if asyncsendResponse != nil && asyncsendResponse.ErrCode == 0 {
			if err := models.UpdateDevtodoFlag(todo.ID); err != nil {
				logging.Error(fmt.Sprintf("%v update dev_flag err:%v", todo.ID, err))
			}
		}
	}
}

//遍历一遍发送标志为0的信息，通知钉钉发送记事本消息
func NoteMessageDingding() {
	notes, err := models.GetNoteFlag()
	if err != nil {
		logging.Error(fmt.Sprintf("get note_flag err:%v", err))
		return
	}
	for _, note := range notes {
		tcmprJson := dingtalk.NoteMseesageToDingding(note)
		asyncsendResponse := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendResponse != nil && asyncsendResponse.ErrCode == 0 {
			if err := models.UpdateNoteFlag(note.ID); err != nil {
				logging.Error(fmt.Sprintf("%v update note_flag err:%v", note.ID, err))
			}
		}
	}
}

//遍历一遍发送标志为0的信息，通知钉钉发送值班通知消息
func OndutyMessageDingding() {
	ods, err := models.GetOndutyFlag()
	if err != nil {
		logging.Error(fmt.Sprintf("get onduty_flag err:%v", err))
		return
	}
	for _, od := range ods {
		tcmprJson := dingtalk.OndutyMseesageToDingding(od)
		asyncsendResponse := dingtalk.MessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendResponse != nil && asyncsendResponse.ErrCode == 0 {
			if err := models.UpdateOndutyFlag(od.ID); err != nil {
				logging.Error(fmt.Sprintf("%v update onduty_flag err:%v", od.ID, err))
			}
		}
	}
}

//遍历一遍发送标志为0的信息，通知钉钉发送提报事项待办工作通知消息
func ProcessMessageDingding() {
	procs, err := models.GetProcessFlag()
	if err != nil {
		logging.Error(fmt.Sprintf("get process_flag err:%v", err))
		return
	}
	for _, proc := range procs {
		p, err := models.GetProcDetail(proc.ProcID)
		if err != nil {
			logging.Error(fmt.Sprintf("get process detail [id:%v] err:%v", proc.ID, err))
		}
		tcmprJson := dingtalk.ProcessMseesageToDingding(p, proc.Czr)
		asyncsendResponse := dingtalk.EappMessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendResponse != nil && asyncsendResponse.ErrCode == 0 {
			if err := models.UpdateProcessFlag(proc.ID); err != nil {
				logging.Error(fmt.Sprintf("%v update process_flag err:%v", proc.ID, err))
			}
		}
	}
}

//遍历一遍发送标志为0的信息，通知钉钉发送提报事项补充描述工作通知
func ProcessBcmsMessageDingding() {
	procs, err := models.GetProcessBcmsFlag()
	if err != nil {
		logging.Error(fmt.Sprintf("get process_flag err:%v", err))
		return
	}
	for _, proc := range procs {
		p, err := models.GetProcDetail(proc.ProcID)
		if err != nil {
			logging.Error(fmt.Sprintf("get process detail [id:%v] err:%v", proc.ProcID, err))
		}
		tcmprJson := dingtalk.ProcessBcmsMseesageToDingding(p)
		asyncsendResponse := dingtalk.EappMessageCorpconversationAsyncsend(tcmprJson)
		if asyncsendResponse != nil && asyncsendResponse.ErrCode == 0 {
			if err := models.UpdateProcessBcmsFlag(proc.ID); err != nil {
				logging.Error(fmt.Sprintf("%v update process_flag err:%v", proc.ID, err))
			}
		}
	}
}

//同步信息
func Sync() {
	logging.Error(fmt.Sprintf("DepartmentUserSync start..."))
	var (
		t  = time.Now().Format("2006-01-02") + " 00:00:00"
		sn = 30 //goroutine数量
		wt = 20 //发送递归请求的次数
	)
	count, e := dingtalk.OrgUserCount(wt)
	if e != nil {
		log.Println(e)
	}
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 90)
		DepartmentUserSync(wt, sn)
		if userNum, err := models.CountUserSyncNum(t); err == nil && count-userNum < 5 {
			break
		}
	}
	logging.Error(fmt.Sprintf("DepartmentUserSync success"))
}

//同步一次部门和人员信息
func DepartmentUserSync(wt, syncNum int) {
	var (
		cntChan = make(chan int)
		wg      = &sync.WaitGroup{}
	)
	deptIds, err := dingtalk.SubDepartmentList(wt)
	if err != nil {
		//logging.Error(fmt.Sprintf("%v", err))
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
