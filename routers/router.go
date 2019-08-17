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

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		//上传文件
		apiv1.POST("/file/upload", api.UploadFile)
		//生成海报
		//apiv1.POST("/poster/generate", v1.GeneratePoster)
		//免登
		apiv1.POST("/login", dingtalk.Login)
		//鉴权
		apiv1.GET("/js_api_config", dingtalk.JsApiConfig)

		//注册事件回调
		apiv1.GET("/callback/reg", dingtalk.RegisterCallback)
		// 查询事件回调
		apiv1.GET("/callback/query", dingtalk.QueryCallback)
		// 更新事件回调
		apiv1.POST("/callback/update", dingtalk.UpdateCallback)
		// 删除事件回调
		apiv1.GET("/callback/delete", dingtalk.DeleteCallback)
		// 获取回调失败的结果
		apiv1.GET("/callback/failed", dingtalk.GetFailedCallbacks)
		//获取回调的结果
		apiv1.POST("/callback/detail", dingtalk.GetCallbacks)

		//同步一次部门用户信息
		apiv1.GET("/sync", dingtalk.DepartmentUserSync)
		//获取部门用户信息同步条数
		apiv1.GET("/syncnum", dingtalk.DepartmentUserSyncNum)

		//获取部门详情
		apiv1.GET("/department/detail", dingtalk.GetDepartmentByID)
		//获取部门列表
		apiv1.GET("/department/list", dingtalk.GetDepartmentByParentID)
		//获取部门用户列表
		apiv1.GET("/user/list", dingtalk.GetUserByDepartmentID)

		//发消息
		apiv1.POST("/msg/send", dingtalk.SendMsg)
		//发消息(内网)
		apiv1.POST("/msg/sendmobile", dingtalk.SendMsgMobile)
		//获取消息列表
		apiv1.GET("/msg/list", dingtalk.GetMsgs)
		//获取消息详情
		apiv1.GET("/msg/detail", dingtalk.GetMsgByID)
		//删除消息
		apiv1.GET("/msg/delete", dingtalk.DeleteMsg)
	}
	return r
}
