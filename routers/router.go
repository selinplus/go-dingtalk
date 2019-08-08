package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/middleware/cors"
	"github.com/selinplus/go-dingtalk/pkg/export"
	"github.com/selinplus/go-dingtalk/pkg/qrcode"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"github.com/selinplus/go-dingtalk/routers/api"
	"github.com/selinplus/go-dingtalk/routers/api/v1/dingtalk"
	"net/http"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.CORSMiddleware())
	//port := strconv.Itoa(setting.ServerSetting.HttpPort + 1)

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	//r.POST("/file/upload", api.UploadImage)
	r.POST("/uploads/file", api.UploadImageByID)

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		//上传文件
		apiv1.POST("/file/upload", api.UploadFile)
		//生成海报
		//apiv1.POST("/poster/generate", v1.GeneratePoster)
		apiv1.POST("/login", dingtalk.Login)
		apiv1.GET("/js_api_config", dingtalk.JsApiConfig)
		//部门同步
		apiv1.GET("/syncinfo", dingtalk.DepartmentUserSync)
		//发消息
		apiv1.POST("/msg/send", dingtalk.SendMsg)
		//获取消息列表
		apiv1.GET("/msg/list", dingtalk.GetMsgs)
		//获取消息详情
		apiv1.GET("/msg/detail", dingtalk.GetMsgByID)
		//删除消息
		apiv1.GET("/msg/delete", dingtalk.DeleteMsg)
	}
	return r
}
