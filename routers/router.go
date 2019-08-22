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
	"github.com/selinplus/go-dingtalk/routers/api/v2"
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

	//内网
	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		//上传文件
		apiv1.POST("/file/upload", api.UploadFile)

		//获取部门用户信息同步条数
		apiv1.GET("/syncnum", dingtalk.DepartmentUserSyncNum)
		//获取部门详情
		apiv1.GET("/department/detail", dingtalk.GetDepartmentByID)
		//获取用户部门详情（内网）
		apiv1.GET("/department/mobile", dingtalk.GetDepartmentByUserMobile)
		//获取部门列表
		apiv1.GET("/department/list", dingtalk.GetDepartmentByParentID)
		//获取部门用户列表
		apiv1.GET("/user/list", dingtalk.GetUserByDepartmentID)

		//发消息(内网)
		apiv1.POST("/msg/sendmobile", dingtalk.SendMsgMobile)
		//获取消息列表
		apiv1.GET("/msg/list", dingtalk.GetMsgs)
		//获取消息详情
		apiv1.GET("/msg/detail", dingtalk.GetMsgByID)
		//删除消息
		apiv1.GET("/msg/delete", dingtalk.DeleteMsg)

		//单项录入
		apiv1.POST("/dev/add", dingtalk.AddDevice)
		//批量导入
		apiv1.POST("/dev/imp", dingtalk.ImpDevices)
		//流转登记
		apiv1.POST("/dev/mod", dingtalk.AddDeviceMod)

		//查询设备状态代码
		apiv1.GET("/dev/state", dingtalk.GetDevstate)
		//查询设备类型代码
		apiv1.GET("/dev/type", dingtalk.GetDevtype)
		//查询操作类型代码
		apiv1.GET("/dev/op", dingtalk.GetDevOp)
	}
	//外网
	apiv2 := r.Group("/api/v2")
	{
		//上传文件
		apiv2.POST("/file/upload", api.UploadFile)
		//免登
		apiv2.POST("/login", v2.Login)
		//鉴权
		apiv2.GET("/js_api_config", v2.JsApiConfig)

		//注册事件回调
		apiv2.GET("/callback/reg", v2.RegisterCallback)
		// 查询事件回调
		apiv2.GET("/callback/query", v2.QueryCallback)
		// 更新事件回调
		apiv2.POST("/callback/update", v2.UpdateCallback)
		// 删除事件回调
		apiv2.GET("/callback/delete", v2.DeleteCallback)
		// 获取回调失败的结果
		apiv2.GET("/callback/failed", v2.GetFailedCallbacks)
		//获取回调的结果
		apiv2.POST("/callback/detail", v2.GetCallbacks)

		//同步一次部门用户信息
		apiv2.GET("/sync", v2.DepartmentUserSync)
		//获取部门用户信息同步条数
		apiv2.GET("/syncnum", v2.DepartmentUserSyncNum)

		//获取部门详情
		apiv2.GET("/department/detail", v2.GetDepartmentByID)
		//获取部门列表
		apiv2.GET("/department/list", v2.GetDepartmentByParentID)
		//获取部门用户列表
		apiv2.GET("/user/list", v2.GetUserByDepartmentID)

		//发消息
		apiv2.POST("/msg/send", v2.SendMsg)
		//获取消息列表
		apiv2.GET("/msg/list", v2.GetMsgs)
		//获取消息详情
		apiv2.GET("/msg/detail", v2.GetMsgByID)
		//删除消息
		apiv2.GET("/msg/delete", v2.DeleteMsg)

		//单项录入
		apiv2.POST("/dev/add", v2.AddDevice)
		//批量导入
		apiv2.POST("/dev/imp", v2.ImpDevices)
		//流转登记
		apiv2.POST("/dev/mod", v2.AddDeviceMod)

		//查询设备状态代码
		apiv2.GET("/dev/state", v2.GetDevstate)
		//查询设备类型代码
		apiv2.GET("/dev/type", v2.GetDevtype)
		//查询操作类型代码
		apiv2.GET("/dev/op", v2.GetDevOp)
	}
	return r
}
