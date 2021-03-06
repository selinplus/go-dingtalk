package ydksrv

import (
	"encoding/json"
	"fmt"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/file"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-gin-web/pkg/logging"
	"log"
	"os"
	"strings"
	"time"
)

// 遍历 ydks 消息，通知钉钉发送待办任务
func Ydksworkrecord() {
	if models.GetWorkrecordSendCnt() > setting.YdksAppSetting.FlowLimit {
		logging.Info("****以地控税****待办任务推送已超过流量限制!!!")
		logging.Info("****以地控税****待办任务推送已超过流量限制!!!")
		logging.Info("****以地控税****待办任务推送已超过流量限制!!!")
		return
	}
	records, err := models.GetWorkrecordFlag()
	if err != nil {
		logging.Error(fmt.Sprintf("get GetWorkrecord record err:%v", err))
		return
	}
	for _, record := range records {
		if record.Lb == "updRecord" { //更新待办任务状态
			asyncsendResponse, err := dingtalk.YdksWorkrecordUpdate(record.UserID, record.RecordID)
			if err != nil {
				log.Printf("%v update Workrecord err:%v", record.ID, err)
				continue
			}
			log.Println(asyncsendResponse)
			if asyncsendResponse != nil && asyncsendResponse.ErrCode == 0 {
				upd := map[string]interface{}{
					"flag_notice": "2",
					"result":      asyncsendResponse.Result,
					"tsrq":        time.Now().Format("2006-01-02 15:04:05"),
				}
				err := models.UpdateWorkrecordFlag(record.ID, upd)
				if err != nil {
					log.Printf("%v update Workrecord Flag err:%v", record.ID, err)
				}
			}
		} else { //发起待办任务
			asyncsendResponse, err := dingtalk.YdksWorkrecordAdd(record.Req)
			if err != nil {
				log.Printf("%v add Workrecord err:%v", record.ID, err)
				continue
			}
			log.Println(asyncsendResponse)
			if asyncsendResponse != nil {
				var upd = map[string]interface{}{}
				if asyncsendResponse.ErrCode == 0 {
					upd = map[string]interface{}{
						"flag_notice": "2",
						"record_id":   asyncsendResponse.RecordId,
						"tsrq":        time.Now().Format("2006-01-02 15:04:05"),
					}
				}
				if asyncsendResponse.ErrCode == 854001 { //重复任务
					upd = map[string]interface{}{
						"flag_notice": "2",
						"record_id":   asyncsendResponse.ErrMsg,
						"tsrq":        time.Now().Format("2006-01-02 15:04:05"),
					}
				}
				err := models.UpdateWorkrecordFlag(record.ID, upd)
				if err != nil {
					log.Printf("%v update Workrecord Flag err:%v", record.ID, err)
				}
			}
		}
	}
}

type WriteData struct {
	Data   []string `json:"data"`
	ErrMsg string   `json:"err_msg"`
}

// 将前一天 ydks 数据写入文件
func WriteIntoFile(date string) {
	lbs := []string{"td", "cjfc", "czfc"}
	for _, lb := range lbs {
		fileName := GetYdksFullPath() + "ydks-" + lb + "data-" + date + ".json"
		dstFile, err := os.Create(fileName)
		if err != nil {
			log.Println("Create data file err:", err)
			return
		}
		defer dstFile.Close()

		var datas []string
		var errMsg string
		list, err := models.GetYdksdata(date, lb)
		if err != nil {
			errMsg = err.Error()
		} else {
			for _, ytst := range list {
				datas = append(datas, ytst.Data)
			}
		}
		writeData := WriteData{
			Data:   datas,
			ErrMsg: errMsg,
		}
		jstr, err := json.Marshal(&writeData)
		if err != nil {
			logging.Error("ytst json.marshal failed,err:", err)
			return
		}
		dstFile.WriteString(string(jstr))
	}

}

// GetYdksFullPath get full save path
func GetYdksFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.YdksAppSetting.YdksSavePath
}

// 清理超过30天的以地控税文件
func CleanUpYdksFiles() {
	dirpath := GetYdksFullPath()
	files, err := file.FindFilesOlderThanDate(dirpath, 30)
	errNotExist := "open : The system cannot find the file specified."
	if err != nil && err.Error() != errNotExist {
		log.Println("CleanUp Ydks Files err:", err)
		return
	}
	for _, fileInfo := range files {
		if !strings.Contains(fileInfo.Name(), "ydks") {
			continue
		}
		err = os.Remove(dirpath + fileInfo.Name())
		if err != nil {
			log.Println("CleanUp Ydks Files err:", err)
		}
	}
}
