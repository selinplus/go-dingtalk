package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/middleware/cors"
	"github.com/selinplus/go-dingtalk/middleware/jwt"
	"github.com/selinplus/go-dingtalk/middleware/ot"
	"github.com/selinplus/go-dingtalk/middleware/sec"
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

	//上传文件
	r.POST("/file/upload", api.UploadFile)
	//清理文件
	r.GET("/file/cleanup", api.CleanUpFile)

	//内网
	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取部门用户信息同步条数
		apiv1.GET("/syncnum", dingtalk.DepartmentUserSyncNum)
		//获取部门详情
		apiv1.GET("/department/detail", dingtalk.GetDepartmentByID)
		//获取用户部门详情（内网）
		apiv1.GET("/department/mobile", dingtalk.GetDepartmentByUserMobile)
		//获取部门列表
		apiv1.GET("/department/list", dingtalk.GetDepartmentByParentID)
		//获取部门列表(不含外部部门)
		apiv1.GET("/department/innerlist", dingtalk.GetDepartmentByParentIDWithNoOuter)
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

		//新建记事本
		apiv1.POST("/note/add", dingtalk.AddNote)
		//删除记事本
		apiv1.GET("/note/delete", dingtalk.DeleteNote)
		//修改记事本内容
		apiv1.POST("/note/update", dingtalk.UpdateNote)
		//获取记事本列表
		apiv1.GET("/note/list", dingtalk.GetNoteList)
		//查询记事本详情
		apiv1.GET("/note/detail", dingtalk.GetNoteDetail)

		//获取当前文件夹文件列表
		apiv1.GET("/netdisk/list", dingtalk.GetFileListByDir)
		//上传网盘文件
		apiv1.POST("/netdisk/upload", dingtalk.AddNetdiskFile)
		//修改网盘文件&从回收站恢复
		apiv1.POST("/netdisk/update", dingtalk.UpdateNetdiskFile)
		//移动到回收站
		apiv1.GET("/netdisk/trash", dingtalk.MoveToTrash)
		//删除文件
		apiv1.GET("/netdisk/delete", dingtalk.DeleteNetdiskFile)

		//获取用户网盘文件夹列表
		apiv1.GET("/netdisk/tree", dingtalk.GetNetdiskDirTree)
		//新建网盘文件夹
		apiv1.POST("/netdisk/mkdir", dingtalk.AddNetdiskDir)
		//修改文件夹
		apiv1.POST("/netdisk/updatedir", dingtalk.UpdateNetdiskDir)
		//删除文件夹
		apiv1.GET("/netdisk/deldir", dingtalk.DeleteNetdiskDir)

		//增加设备管理机构
		apiv1.POST("/dev/adddept", dingtalk.AddDevdept)
		//修改设备管理机构
		apiv1.POST("/dev/upddept", dingtalk.UpdateDevdept)
		//获取设备管理机构列表(树结构)
		apiv1.GET("/dev/tree", dingtalk.GetDevdeptTree)
		//获取设备管理机构列表(循环遍历)
		apiv1.GET("/dev/deptlist", dingtalk.GetDevdeptBySjjgdm)
		//获取设备管理机构列表(bz:0-管理员不可选;1-管理员可选)
		apiv1.GET("/dev/deptglylist", dingtalk.GetDevdeptGlyList)
		//删除设备管理机构
		apiv1.GET("/dev/deldept", dingtalk.DeleteDevdept)
		//获取当前机构管理员信息
		apiv1.GET("/dev/deptgly", dingtalk.GetDevdeptGly)
		//获取当前用户为机构管理员的所有机构列表
		apiv1.GET("/dev/gly", dingtalk.GetDevGly)
		//增加设备使用人员
		apiv1.POST("/dev/adduser", dingtalk.AddDevuser)
		//修改设备使用人员
		apiv1.POST("/dev/upduser", dingtalk.UpdateDevuser)
		//获取设备使用人员列表
		apiv1.GET("/dev/userlist", dingtalk.GetDevuserList)
		//删除设备使用人员
		apiv1.GET("/dev/deluser", dingtalk.DeleteDevuser)

		//单项录入
		apiv1.POST("/dev/add", dingtalk.AddDevinfo)
		//批量导入
		apiv1.POST("/dev/imp", dingtalk.ImpDevinfos)
		//更新设备信息
		apiv1.POST("/dev/update", dingtalk.UpdateDevinfo)
		//查询设备详情
		apiv1.GET("/dev/detail", dingtalk.GetDevinfoByID)
		//获取设备列表(多条件查询)
		apiv1.GET("/dev/list", dingtalk.GetDevinfos)
		//获取当前操作人所有流水记录
		apiv1.GET("/devmod/lslist", dingtalk.GetDevMods)
		//根据流水号查询记录
		apiv1.GET("/devmod/lsdetail", dingtalk.GetDevModetails)

		/*
			//单项录入
			apiv1.POST("/dev/add", dingtalk.AddDevice)
			//批量导入
			apiv1.POST("/dev/imp", dingtalk.ImpDevices)
			//更新设备信息
			apiv1.POST("/dev/update", dingtalk.UpdateDevice)
			//查询设备详情
			apiv1.GET("/dev/detail", dingtalk.GetDeviceByID)
			//获取设备列表
			apiv1.GET("/dev/list", dingtalk.GetDevices)
		*/

		//查询设备信息及当前使用状态详情
		apiv1.GET("/dev/detailmod", dingtalk.GetDeviceModByDevID)
		//流转登记
		apiv1.POST("/devmod/add", dingtalk.AddDeviceMod)
		//设备流水记录查询
		apiv1.GET("/devmod/list", dingtalk.GetDevModList)

		//查询设备状态代码
		apiv1.GET("/dev/state", dingtalk.GetDevstate)
		//查询设备类型代码
		apiv1.GET("/dev/type", dingtalk.GetDevtype)
		//查询操作类型代码
		apiv1.GET("/dev/op", dingtalk.GetDevOp)
		//查询操作属性代码
		apiv1.GET("/dev/prop", dingtalk.GetDevProp)
	}
	//外网----消息助手&&记事本&&小网盘
	apiv2 := r.Group("/api/v2")
	apiv2.Use(sec.Sec())
	{
		//上传文件
		apiv2.POST("/file/upload", api.UploadFile)
		//文件下载
		apiv2.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

		//免登
		apiv2.POST("/login", dingtalk.Login)
		//鉴权
		apiv2.GET("/js_api_config", dingtalk.JsApiConfig)

		//注册事件回调
		apiv2.GET("/callback/reg", dingtalk.RegisterCallback)
		// 查询事件回调
		apiv2.GET("/callback/query", dingtalk.QueryCallback)
		// 更新事件回调
		apiv2.POST("/callback/update", dingtalk.UpdateCallback)
		// 删除事件回调
		apiv2.GET("/callback/delete", dingtalk.DeleteCallback)
		// 获取回调失败的结果
		apiv2.GET("/callback/failed", dingtalk.GetFailedCallbacks)
		//获取回调的结果
		apiv2.POST("/callback/detail", dingtalk.GetCallbacks)

		//同步一次部门用户信息
		apiv2.GET("/sync", dingtalk.DepartmentUserSync)
		//获取当日部门用户信息同步条数
		apiv2.GET("/syncnum", dingtalk.DepartmentUserSyncNum)
		//获取企业员工人数
		apiv2.GET("/usercount", dingtalk.OrgUserCount)

		//获取多部门详情时，排除outer属性
		apiv2.GET("/department/outer", dingtalk.GetDepartmentByIDs)
		//获取部门详情
		apiv2.GET("/department/detail", dingtalk.GetDepartmentByID)
		//获取部门列表
		apiv2.GET("/department/list", dingtalk.GetDepartmentByParentID)
		//获取部门用户列表
		apiv2.GET("/user/list", dingtalk.GetUserByDepartmentID)

		//发消息
		apiv2.POST("/msg/send", dingtalk.SendMsg)
		//获取消息列表
		apiv2.GET("/msg/list", dingtalk.GetMsgs)
		//获取消息详情
		apiv2.GET("/msg/detail", dingtalk.GetMsgByID)
		//删除消息
		apiv2.GET("/msg/delete", dingtalk.DeleteMsg)

		//新建记事本
		apiv2.POST("/note/add", dingtalk.AddNote)
		//删除记事本
		apiv2.GET("/note/delete", dingtalk.DeleteNote)
		//修改记事本内容
		apiv2.POST("/note/update", dingtalk.UpdateNote)
		//获取记事本列表
		apiv2.GET("/note/list", dingtalk.GetNoteList)
		//查询记事本详情
		apiv2.GET("/note/detail", dingtalk.GetNoteDetail)

		//获取当前文件夹文件列表
		apiv2.GET("/netdisk/list", dingtalk.GetFileListByDir)
		//上传网盘文件
		apiv2.POST("/netdisk/upload", dingtalk.AddNetdiskFile)
		//修改网盘文件&从回收站恢复
		apiv2.POST("/netdisk/update", dingtalk.UpdateNetdiskFile)
		//移动到回收站
		apiv2.GET("/netdisk/trash", dingtalk.MoveToTrash)
		//删除文件
		apiv2.GET("/netdisk/delete", dingtalk.DeleteNetdiskFile)

		//获取用户网盘文件夹列表
		apiv2.GET("/netdisk/tree", dingtalk.GetNetdiskDirTree)
		//新建网盘文件夹
		apiv2.POST("/netdisk/mkdir", dingtalk.AddNetdiskDir)
		//修改网盘文件夹
		apiv2.POST("/netdisk/updatedir", dingtalk.UpdateNetdiskDir)
		//删除网盘文件夹
		apiv2.GET("/netdisk/deldir", dingtalk.DeleteNetdiskDir)
	}
	//外网----设备管理&&事件提报
	apiv3 := r.Group("/api/v3")
	apiv3.Use(ot.OT())
	{
		//上传文件
		apiv3.POST("/file/upload", api.UploadFile)
		//文件下载
		apiv3.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

		//免登
		apiv3.POST("/login", dingtalk.Login)
		//鉴权
		apiv3.GET("/js_api_config", dingtalk.JsApiConfig)

		//获取多部门详情时，排除outer属性
		apiv3.GET("/department/outer", dingtalk.GetDepartmentByIDs)
		//获取部门详情
		apiv3.GET("/department/detail", dingtalk.GetDepartmentByID)
		//获取部门列表
		apiv3.GET("/department/list", dingtalk.GetDepartmentByParentID)
		//获取部门用户列表
		apiv3.GET("/user/list", dingtalk.GetUserByDepartmentID)

		//获取设备管理机构列表
		apiv3.GET("/dev/tree", dingtalk.GetDevdeptTree)
		//获取当前用户设备列表
		apiv3.GET("/dev/listbyuser", dingtalk.GetDevicesByUser)
		//查询设备信息及当前使用状态详情
		apiv3.GET("/dev/detailmod", dingtalk.GetDeviceModByDevID)
		//设备流水记录查询
		apiv3.GET("/devmod/list", dingtalk.GetDevModList)

		//提报事项保存&&提交
		apiv3.POST("/proc/add", dingtalk.AddProc)
		//(未)提交事项修改&&提交
		apiv3.POST("/proc/update", dingtalk.UpdateProc)
		//作废&&删除提报事项
		apiv3.GET("/proc/delete", dingtalk.DeleteProc)
		//查询提报事项详情
		apiv3.GET("/proc/detail", dingtalk.GetProcDetail)
		//获取待办列表
		apiv3.GET("/proc/todolist", dingtalk.GetProcTodoList)
		//获取已办列表(全部)
		apiv3.GET("/proc/donelist", dingtalk.GetProcDoneList)
		//获取已办列表(已办结)
		apiv3.GET("/proc/donelistend", dingtalk.GetProcDoneListEnd)
		//获取已办列表(未办结)
		apiv3.GET("/proc/donelistdoing", dingtalk.GetProcDoneListDoing)
		//事件处理(退回&&通过)
		apiv3.POST("/proc/deal", dingtalk.DealProc)
		//事件处理流水记录查询
		apiv3.GET("/proc/list", dingtalk.GetProcModList)
		//发起补充描述
		apiv3.POST("/proc/bcms", dingtalk.ProcBcms)
		//补充描述提交
		apiv3.POST("/proc/commitbcms", dingtalk.UpdateProcBcms)

		//获取手工提报人员列表
		apiv3.GET("/proc/custlist", dingtalk.GetProcCustomizeList)
		//获取下一节点操作人
		apiv3.GET("/proc/czr", dingtalk.GetProcCzr)
		//查询提报类型代码
		apiv3.GET("/proc/type", dingtalk.GetProcType)
	}
	return r
}
