package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/models"
	"github.com/selinplus/go-dingtalk/pkg/cron"
	"github.com/selinplus/go-dingtalk/pkg/dingtalk"
	"github.com/selinplus/go-dingtalk/pkg/gredis"
	"github.com/selinplus/go-dingtalk/pkg/logging"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"github.com/selinplus/go-dingtalk/routers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	util.Setup()
	cron.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/selinplus/go-dingtalk
// @license.name MIT
// @license.url https://github.com/selinplus/go-dingtalk/blob/master/LICENSE
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	if len(os.Args) == 2 {
		models.InitDb()
		log.Println("*******database init over*****")
		log.Println("*******please rerun the program*****")
		return
	}
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dir)
	log.Printf("[info] start http server listening %s", endPoint)

	go func() {
		time.Sleep(time.Second * 10)
		dingtalk.RegCallbackInit()
	}()

	err = server.ListenAndServe()
	if err != nil {
		log.Printf("init listen server fail:%v", err)
	}

	// If you want Graceful Restart, you need a Unix system and download github.com/fvbock/endless
	//endless.DefaultReadTimeOut = readTimeout
	//endless.DefaultWriteTimeOut = writeTimeout
	//endless.DefaultMaxHeaderBytes = maxHeaderBytes
	//server := endless.NewServer(endPoint, routersInit)
	//server.BeforeBegin = func(add string) {
	//	log.Printf("Actual pid is %d", syscall.Getpid())
	//}
	//
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Printf("Server err: %v", err)
	//}
}
