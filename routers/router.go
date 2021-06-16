package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/middleware/cors"
	"github.com/selinplus/go-dingtalk/middleware/eappm"
	"github.com/selinplus/go-dingtalk/middleware/h5m"
	"github.com/selinplus/go-dingtalk/middleware/jwt"
	"github.com/selinplus/go-dingtalk/middleware/ydksjwt"
	"github.com/selinplus/go-dingtalk/pkg/export"
	"github.com/selinplus/go-dingtalk/pkg/fsdjsrv"
	"github.com/selinplus/go-dingtalk/pkg/qrcode"
	"github.com/selinplus/go-dingtalk/pkg/upload"
	"github.com/selinplus/go-dingtalk/pkg/ydksrv"
	"github.com/selinplus/go-dingtalk/routers/api"
	"github.com/selinplus/go-dingtalk/routers/api/fsdj"
	"github.com/selinplus/go-dingtalk/routers/api/v1/dev"
	"github.com/selinplus/go-dingtalk/routers/api/v1/dingtalk"
	"github.com/selinplus/go-dingtalk/routers/api/v1/h5"
	"github.com/selinplus/go-dingtalk/routers/api/ydks"
	"net/http"
)

//InitRouter initialize routing information
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
	//福山品牌党建文件下载
	r.StaticFS("/fsdj/dj_image", http.Dir(fsdjsrv.GetFsdjEappImageFullPath()))

	//调用智税平台登陆
	r.POST("/login", api.Login)
	//代理转发智税平台登陆/获取路由表
	r.POST("/r_login", api.Rlogin)
	r.GET("/r_route", api.GetRoutes)
	//福山品牌党建
	r.GET("/fsdj", fsdj.Home)
	r.Static("/css", "runtime/static/css")
	r.Static("/js", "runtime/static/js")
	r.Static("/img", "runtime/static/img")
	r.Static("/logo.png", "runtime/static/logo.png")
	r.Static("/favicon.ico", "runtime/static/favicon.ico")

	//上传文件
	r.POST("/file/upload", api.UploadFile)
	//清理文件
	r.GET("/file/cleanup", api.CleanUpFile)
	//接收值班通知推送消息
	r.POST("/onduty", api.OnDuty)

	//======================== 内网 ==========================//
	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取部门用户信息同步条数
		apiv1.GET("/syncnum", dingtalk.DepartmentUserSyncNum)
		//获取部门详情
		apiv1.GET("/department/detail", dingtalk.GetDepartmentByID)
		//获取用户部门详情(内网)
		apiv1.GET("/department/mobile", dingtalk.GetDepartmentByUserMobile)
		//获取部门列表
		apiv1.GET("/department/list", dingtalk.GetDepartmentByParentID)
		//获取部门列表(不含外部部门)
		apiv1.GET("/department/innerlist", dingtalk.GetDepartmentByParentIDWithNoOuter)
		//获取部门用户列表
		apiv1.GET("/user/list", dingtalk.GetUserByDepartmentID)
		//模糊查询用户
		apiv1.GET("/user/mc", dingtalk.GetUserByMc)

		//======================== 消息助手 ==========================//
		//获取最近联系人列表
		apiv1.GET("/msg/userlist", h5.GetRecentContacter)
		//发消息(内网)
		apiv1.POST("/msg/sendmobile", h5.SendMsgMobile)
		//获取消息列表
		apiv1.GET("/msg/list", h5.GetMsgs)
		//获取消息详情
		apiv1.GET("/msg/detail", h5.GetMsgByID)
		//删除消息
		apiv1.GET("/msg/delete", h5.DeleteMsg)

		//增加通讯录组
		apiv1.POST("/msg/addbook", h5.AddAddressbook)
		//删除通讯录组
		apiv1.GET("/msg/delbook", h5.DeleteAddressbook)
		//修改通讯录组名称
		apiv1.POST("/msg/updbook", h5.UpdateAddressbook)
		//获取通讯组列表
		apiv1.GET("/msg/getbooks", h5.GetAddressbooks)
		//增加联系人
		apiv1.POST("/msg/addcontacter", h5.AddContacter)
		//删除联系人
		apiv1.POST("/msg/delcontacter", h5.DeleteContacter)
		//获取通讯组联系人列表
		apiv1.GET("/msg/getcontacters", h5.GetContacters)

		//======================== 记事本 ==========================//
		//新建记事本
		apiv1.POST("/note/add", h5.AddNote)
		//删除记事本
		apiv1.GET("/note/delete", h5.DeleteNote)
		//修改记事本内容
		apiv1.POST("/note/update", h5.UpdateNote)
		//获取记事本列表
		apiv1.GET("/note/list", h5.GetNoteList)
		//查询记事本详情
		apiv1.GET("/note/detail", h5.GetNoteDetail)

		//======================== 钉钉网盘 ==========================//
		//获取当前文件夹文件列表
		apiv1.GET("/netdisk/list", h5.GetFileListByDir)
		//上传网盘文件
		apiv1.POST("/netdisk/upload", h5.AddNetdiskFile)
		//修改网盘文件&从回收站恢复
		apiv1.POST("/netdisk/update", h5.UpdateNetdiskFile)
		//移动到回收站
		apiv1.GET("/netdisk/trash", h5.MoveToTrash)
		//删除文件
		apiv1.GET("/netdisk/delete", h5.DeleteNetdiskFile)

		//获取用户网盘文件夹列表
		apiv1.GET("/netdisk/tree", h5.GetNetdiskDirTree)
		//新建网盘文件夹
		apiv1.POST("/netdisk/mkdir", h5.AddNetdiskDir)
		//修改文件夹
		apiv1.POST("/netdisk/updatedir", h5.UpdateNetdiskDir)
		//删除文件夹
		apiv1.GET("/netdisk/deldir", h5.DeleteNetdiskDir)

		//======================== 设备管理 ==========================//
		//增加设备管理机构
		apiv1.POST("/dev/adddept", dev.AddDevdept)
		//修改设备管理机构
		apiv1.POST("/dev/upddept", dev.UpdateDevdept)
		//导出人员&机构代码表
		apiv1.GET("/dev/exp_info", dev.ExportDevdepUserInfo)
		//获取设备管理机构信息
		apiv1.GET("/dev/deptinfo", dev.GetDept)
		//获取设备管理机构上级有管理员的机构信息
		apiv1.GET("/dev/ptdeptinfo", dev.GetParentDept)
		//获取设备管理机构列表(树结构)
		apiv1.GET("/dev/tree", dev.GetDevdeptTree)
		//获取设备管理机构列表(循环遍历)
		apiv1.GET("/dev/deptlist", dev.GetDevdeptBySjjgdm)
		//获取设备管理机构列表(bz:0-管理员不可选;1-管理员可选)
		apiv1.GET("/dev/deptglylist", dev.GetDevdeptGlyList)
		//获取设备管理机构列表
		apiv1.GET("/dev/deptbgrlist", dev.GetDevdeptBgrList)
		//删除设备管理机构
		apiv1.GET("/dev/deldept", dev.DeleteDevdept)
		//删除当前机构管理员&保管人
		apiv1.GET("/dev/delgly", dev.DelDevdeptGly)
		//获取当前机构管理员信息
		apiv1.GET("/dev/deptgly", dev.GetDevdeptGly)
		//获取当前机构保管人信息
		apiv1.GET("/dev/deptbgr", dev.GetDevdeptBgr)
		//获取当前用户为机构管理员的所有机构列表
		apiv1.GET("/dev/gly", dev.GetDevGly)

		//增加设备使用人员
		apiv1.POST("/dev/adduser", dev.AddDevuser)
		//修改设备使用人员
		apiv1.POST("/dev/upduser", dev.UpdateDevuser)
		//获取设备使用人员列表
		apiv1.GET("/dev/userlist", dev.GetDevuserList)
		//删除设备使用人员
		apiv1.GET("/dev/deluser", dev.DeleteDevuser)
		//登录后获取人员身份信息
		apiv1.GET("/dev/login", dev.LoginInfo)

		//单项录入
		apiv1.POST("/dev/add", dev.AddDevinfo)
		//批量导入
		apiv1.POST("/dev/imp", dev.ImpDevinfos)
		//更新设备存放位置
		apiv1.POST("/dev/update_cfwz", dev.UpdateDevinfoCfwz)
		//更新设备管理机构、使用人、所属机构、所属位置
		apiv1.POST("/dev/update_admin", dev.UpdateDevinfoByAdmin)
		//更新设备信息
		apiv1.POST("/dev/update", dev.UpdateDevinfo)
		//更新设备二维码打印次数
		apiv1.POST("/dev/update_qrcode", dev.UpdateDevinfoPnum)
		//删除设备信息
		apiv1.GET("/dev/del", dev.DelDevinfo)
		//查询设备详情
		apiv1.GET("/dev/detail", dev.GetDevinfoByID)
		//获取设备列表(inner多条件查询设备)
		apiv1.GET("/dev/list", dev.GetDevinfos)
		//获取设备列表(管理员端,多条件查询设备)
		apiv1.GET("/dev/listgly", dev.GetDevinfosGly)
		//导出设备清册
		apiv1.GET("/dev/list_export", dev.ExportDevInfosGly)
		//获取当前操作人所有流水记录
		apiv1.GET("/devmod/lslist", dev.GetDevMods)
		//根据流水号查询记录
		apiv1.GET("/devmod/lsdetail", dev.GetDevModetails)
		//设备下发
		apiv1.POST("/dev/issued", dev.Issued)
		//设备机构变更申请
		apiv1.POST("/dev/change_jgks", dev.ChangeJgks)
		//管理员处理设备机构变更申请&交回申请
		apiv1.POST("/dev/jgbg_jhrk", dev.GlyChangeJgks)
		//设备分配(管理员入库)&借出&收回&上交&交回申请
		apiv1.POST("/dev/allocate", dev.Allocate)
		//获取设备列表(管理员查询||eapp使用人查询)
		apiv1.GET("/dev/listbybz", dev.GetDevinfosByUser)
		//设备流水记录查询
		apiv1.GET("/devmod/list", dev.GetDevModList)

		//新增盘点任务
		apiv1.POST("/dev/cktask", dev.GetDevCkTask)
		//获取盘点任务列表
		apiv1.GET("/dev/cktasks", dev.GetDevCkTasks)
		//获取盘点任务清册明细
		apiv1.GET("/dev/ckdetail", dev.GetDevCkDetail)
		//导出盘点任务清册明细
		apiv1.GET("/dev/ckdetail_export", dev.ExportDevCkDetail)

		//根据id获取待办&已办详情
		apiv1.GET("/dev/tododetail", dev.GetDevTodosOrDonesByTodoid)
		//获取待办列表(交回设备)
		apiv1.GET("/dev/todolist", dev.GetDevTodosOrDones)
		//获取已办列表(交回设备)
		apiv1.GET("/dev/donelist", dev.GetDevTodosOrDones)
		//获取待办列表(上交设备)
		apiv1.GET("/dev/uptodolist", dev.GetUpDevTodosOrDones)
		//获取已办列表(上交设备)
		apiv1.GET("/dev/updonelist", dev.GetUpDevTodosOrDones)

		//查询设备状态代码
		apiv1.GET("/dev/state", dev.GetDevstate)
		//查询设备类型代码(树结构)
		apiv1.GET("/dev/type", dev.GetDevtypeTree)
		//查询设备类型代码
		//apiv1.GET("/dev/type", dev.GetDevtype)
		//查询操作类型代码
		apiv1.GET("/dev/op", dev.GetDevOp)
		//查询操作属性代码
		apiv1.GET("/dev/prop", dev.GetDevProp)

		//======================== 福山品牌党建 ==========================//
		//模糊查询福山区用户
		apiv1.GET("/fsdj/user/mc", fsdj.GetFsdjUserByMc)

		//增加学习小组
		apiv1.POST("/fsdj/group/add", fsdj.AddGroup)
		//修改学习小组
		apiv1.POST("/fsdj/group/upd", fsdj.UpdGroup)
		//获取学习小组信息
		apiv1.GET("/fsdj/group/info", fsdj.GetGroup)
		//获取学习小组列表(树结构)
		apiv1.GET("/fsdj/group/tree", fsdj.GetGroupTree)
		//删除学习小组
		apiv1.GET("/fsdj/group/del", fsdj.DelGroup)

		//设置学习小组管理员
		apiv1.POST("/fsdj/gly/add", fsdj.UpdGroup)
		//删除学习小组管理员
		apiv1.GET("/fsdj/gly/del", fsdj.DelGroupGly)
		//获取学习小组管理员信息
		apiv1.GET("/fsdj/gly/group", fsdj.GetGroupGly)

		//增加学习小组成员
		apiv1.POST("/fsdj/member/add", fsdj.AddGroupMember)
		//移动学习小组成员
		apiv1.POST("/fsdj/member/move", fsdj.UpdGroupMember)
		//获取学习小组成员列表
		apiv1.GET("/fsdj/member/list", fsdj.GetGroupMembers)
		//删除学习小组成员
		apiv1.GET("/fsdj/member/del", fsdj.DelGroupMember)

		//上传文件
		apiv1.POST("/fsdj/upload", fsdj.FsdjUploadFile)

		//党建图文发布
		apiv1.POST("/fsdj/topic/post", fsdj.PostStudyTopic)
		//党建图文修改
		apiv1.POST("/fsdj/topic/edit", fsdj.UpdStudyTopic)
		//党建图文审核发布
		apiv1.POST("/fsdj/topic/approve", fsdj.UpdStudyTopic)
		//党建图文撤销
		apiv1.POST("/fsdj/topic/cancel", fsdj.UpdStudyTopic)
		//党建图文分享
		apiv1.POST("/fsdj/topic/share", fsdj.UpdStudyTopic)
		//获取图文详情
		apiv1.GET("/fsdj/topic/detail/:id", fsdj.GetStudyTopic)
		//获取图文列表
		apiv1.GET("/fsdj/topics", fsdj.GetStudyTopics)
		//党建图文删除
		apiv1.GET("/fsdj/topic/del/:id", fsdj.DelStudyTopic)

		//党建活动发布
		apiv1.POST("/fsdj/act/post", fsdj.PostStudyAct)
		//党建活动修改
		apiv1.POST("/fsdj/act/edit", fsdj.UpdStudyAct)
		//党建活动审核发布
		apiv1.POST("/fsdj/act/approve", fsdj.UpdStudyAct)
		//党建活动撤销
		apiv1.POST("/fsdj/act/cancel", fsdj.UpdStudyAct)
		//党建活动分享
		apiv1.POST("/fsdj/act/share", fsdj.UpdStudyAct)
		//获取党建活动详情
		apiv1.GET("/fsdj/act/detail/:id", fsdj.GetStudyAct)
		//获取党建活动列表
		apiv1.GET("/fsdj/acts/list", fsdj.GetStudyActs)
		//党建活动删除
		apiv1.GET("/fsdj/act/del/:id", fsdj.DelStudyAct)
		//获取党建活动参与待信息审核列表
		apiv1.GET("/fsdj/act/join_approves", fsdj.GetApproveStudyActs)
		//党建活动参与信息审核
		apiv1.GET("/fsdj/act/join_approve/:id", fsdj.ApproveStudyAct)
		//查询党建活动参与人员列表
		apiv1.GET("/fsdj/act/users", fsdj.GetStudyActUsers)
		//党建活动参与人员统计
		apiv1.GET("/fsdj/acts/count/:id", fsdj.CountStudyAct)

		//党员风采审核发布&驳回
		apiv1.POST("/fsdj/hlt/approve", fsdj.UpdStudyHlt)
		//党员风采撤销
		apiv1.POST("/fsdj/hlt/cancel", fsdj.UpdStudyHlt)
		//党员风采查看
		apiv1.GET("/fsdj/hlt/detail/:id", fsdj.GetStudyHlt)
		//获取手机号获取党员风采列表
		apiv1.GET("/fsdj/hlts/mylist", fsdj.GetStudyHltsByUserid)
		//获取活动id获取党员风采列表
		apiv1.GET("/fsdj/hlts/list", fsdj.GetStudyHlts)

		//查看党员签到情况
		apiv1.GET("/fsdj/singins", fsdj.GetSigninsByUserid)
		//查看某天全局签到情况
		apiv1.GET("/fsdj/singins/all", fsdj.GetSigninsByQdrq)
	}
	//======================== 外网H5——消息助手&&记事本&&小网盘 ==========================//
	apiv2 := r.Group("/api/v2")
	apiv2.Use(h5m.JWT())
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
		//查询事件回调
		apiv2.GET("/callback/query", dingtalk.QueryCallback)
		//更新事件回调
		apiv2.POST("/callback/update", dingtalk.UpdateCallback)
		//删除事件回调
		apiv2.GET("/callback/delete", dingtalk.DeleteCallback)
		//获取回调失败的结果
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

		//获取最近联系人列表
		apiv2.GET("/msg/userlist", h5.GetRecentContacter)
		//发消息
		apiv2.POST("/msg/send", h5.SendMsg)
		//获取消息列表
		apiv2.GET("/msg/list", h5.GetMsgs)
		//获取消息详情
		apiv2.GET("/msg/detail", h5.GetMsgByID)
		//删除消息
		apiv2.GET("/msg/delete", h5.DeleteMsg)
		//获取通讯组列表
		apiv2.GET("/msg/getbooks", h5.GetAddressbooks)
		//获取通讯组联系人列表
		apiv2.GET("/msg/getcontacters", h5.GetContacters)

		//新建记事本
		apiv2.POST("/note/add", h5.AddNote)
		//删除记事本
		apiv2.GET("/note/delete", h5.DeleteNote)
		//修改记事本内容
		apiv2.POST("/note/update", h5.UpdateNote)
		//获取记事本列表
		apiv2.GET("/note/list", h5.GetNoteList)
		//查询记事本详情
		apiv2.GET("/note/detail", h5.GetNoteDetail)

		//获取当前文件夹文件列表
		apiv2.GET("/netdisk/list", h5.GetFileListByDir)
		//上传网盘文件
		apiv2.POST("/netdisk/upload", h5.AddNetdiskFile)
		//修改网盘文件&从回收站恢复
		apiv2.POST("/netdisk/update", h5.UpdateNetdiskFile)
		//移动到回收站
		apiv2.GET("/netdisk/trash", h5.MoveToTrash)
		//删除文件
		apiv2.GET("/netdisk/delete", h5.DeleteNetdiskFile)

		//获取用户网盘文件夹列表
		apiv2.GET("/netdisk/tree", h5.GetNetdiskDirTree)
		//新建网盘文件夹
		apiv2.POST("/netdisk/mkdir", h5.AddNetdiskDir)
		//修改网盘文件夹
		apiv2.POST("/netdisk/updatedir", h5.UpdateNetdiskDir)
		//删除网盘文件夹
		apiv2.GET("/netdisk/deldir", h5.DeleteNetdiskDir)
	}
	//======================== 外网Eapp——设备管理&&事件提报 ==========================//
	apiv3 := r.Group("/api/v3")
	apiv3.Use(eappm.JWT())
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

		//登录后获取人员身份信息
		apiv3.GET("/dev/login", dev.LoginInfo)
		//获取设备管理机构及人员列表(循环遍历)
		apiv3.GET("/dev/epptree", dev.GetDevdeptEppTree)
		//查询设备详情
		apiv3.GET("/dev/detail", dev.GetDevinfoByID)
		//获取交回设备待入库列表
		apiv3.GET("/dev/tobestored", dev.GetDevinfosToBeStored)
		//更新设备存放位置
		apiv3.POST("/dev/update_cfwz", dev.UpdateDevinfoCfwz)
		//设备流水记录查询
		apiv3.GET("/devmod/list", dev.GetDevModList)
		//设备机构变更申请
		apiv3.POST("/dev/change_jgks", dev.ChangeJgks)
		//管理员处理设备机构变更申请&交回申请
		apiv3.POST("/dev/jgbg_jhrk", dev.GlyChangeJgks)
		//设备分配(管理员入库)&借出&收回&交回申请
		apiv3.POST("/dev/allocate", dev.Allocate)
		//获取设备列表(inner多条件查询设备)
		apiv3.GET("/dev/list", dev.GetDevinfos)
		//获取当前用户为机构管理员的所有机构列表
		apiv3.GET("/dev/gly", dev.GetDevGly)
		//获取设备管理机构列表(bz:0-管理员不可选;1-管理员可选)
		apiv3.GET("/dev/deptglylist", dev.GetDevdeptGlyList)
		//获取设备管理机构列表
		apiv3.GET("/dev/deptbgrlist", dev.GetDevdeptBgrList)
		//设备下发
		apiv3.POST("/dev/issued", dev.Issued)
		//获取设备列表(管理员查询||eapp使用人查询)
		apiv3.GET("/dev/listbybz", dev.GetDevinfosByUser)
		//根据id获取待办&已办详情
		apiv3.GET("/dev/tododetail", dev.GetDevTodosOrDonesByTodoid)
		//获取待办列表(交回设备)
		apiv3.GET("/dev/todolist", dev.GetDevTodosOrDones)
		//获取已办列表(交回设备)
		apiv3.GET("/dev/donelist", dev.GetDevTodosOrDones)
		//获取盘点任务列表
		apiv3.GET("/dev/cktasks", dev.GetDevCkTasks)
		//获取盘点任务清册明细
		apiv3.GET("/dev/ckdetail", dev.GetDevCkDetail)
		//设备盘点拍照上传
		apiv3.POST("/dev/check_img", dev.DevCheckImg)
		//设备盘点
		apiv3.GET("/dev/check", dev.GetDevCheck)

		//提报事项保存&&提交
		apiv3.POST("/proc/add", h5.AddProc)
		//(未)提交事项修改&&提交
		apiv3.POST("/proc/update", h5.UpdateProc)
		//作废&&删除提报事项
		apiv3.GET("/proc/delete", h5.DeleteProc)
		//查询提报事项详情
		apiv3.GET("/proc/detail", h5.GetProcDetail)
		//获取待办列表
		apiv3.GET("/proc/todolist", h5.GetProcTodoList)
		//获取已办列表(全部)
		apiv3.GET("/proc/donelist", h5.GetProcDoneList)
		//获取已办列表(已办结)
		apiv3.GET("/proc/donelistend", h5.GetProcDoneListEnd)
		//获取已办列表(未办结)
		apiv3.GET("/proc/donelistdoing", h5.GetProcDoneListDoing)
		//事件处理(退回&&通过)
		apiv3.POST("/proc/deal", h5.DealProc)
		//事件处理流水记录查询
		apiv3.GET("/proc/list", h5.GetProcModList)
		//发起补充描述
		apiv3.POST("/proc/bcms", h5.ProcBcms)
		//补充描述提交
		apiv3.POST("/proc/commitbcms", h5.UpdateProcBcms)

		//获取手工提报人员列表
		apiv3.GET("/proc/custlist", dev.GetProcCustomizeList)
		//获取下一节点操作人
		apiv3.GET("/proc/czr", dev.GetProcCzr)
		//查询提报类型代码
		apiv3.GET("/proc/type", dev.GetProcType)
	}

	//======================== 外网Eapp——福山品牌党建 ==========================//
	apifsdj := r.Group("/api/fsdj")
	apifsdj.Use(eappm.JWT())
	{
		//上传文件
		apifsdj.POST("/upload", fsdj.FsdjUploadFile)
		//文件下载
		apifsdj.StaticFS("/dj_image", http.Dir(fsdjsrv.GetFsdjEappImageFullPath()))

		//免登
		apifsdj.POST("/login", fsdj.Login)

		//获取学习小组列表(树结构)
		apifsdj.GET("/group/tree", fsdj.GetGroupTree)
		//获取学习小组成员列表
		apifsdj.GET("/member/list", fsdj.GetGroupMembers)

		//获取图文详情
		apifsdj.GET("/topic/detail/:id", fsdj.GetStudyTopic)
		//获取图文列表
		apifsdj.GET("/topics", fsdj.GetStudyTopics)

		//获取党建活动详情
		apifsdj.GET("/act/detail/:id", fsdj.GetStudyAct)
		//获取党建活动列表
		apifsdj.GET("/acts/list", fsdj.GetStudyActs)
		//参与党建活动
		apifsdj.GET("/act/join/:id", fsdj.JoinStudyAct)
		//获取我的党建活动列表
		apifsdj.GET("/acts/mylist", fsdj.GetStudyActsMine)
		//党建活动参与人员统计
		apifsdj.GET("/acts/count/:id", fsdj.CountStudyAct)

		//党员风采发布
		apifsdj.POST("/hlt/post", fsdj.PostStudyHlt)
		//党员风采修改
		apifsdj.POST("/hlt/edit", fsdj.UpdStudyHlt)
		//党员学习笔记类风采自我推选
		apifsdj.POST("/hlt/star_note", fsdj.UpdStudyHlt)
		//党员风采删除
		apifsdj.GET("/hlt/del/:id", fsdj.DelStudyHlt)
		//点赞党员风采
		apifsdj.GET("/hlt/star/:id", fsdj.AddStudyHltStar)
		//取消点赞党员风采
		apifsdj.GET("/hlt/star_cancel/:id", fsdj.CancelStudyHltStar)
		//党员风采查看
		apifsdj.GET("/hlt/detail/:id", fsdj.GetStudyHlt)
		//根据userid获取党员风采列表
		apifsdj.GET("/hlts/mylist", fsdj.GetStudyHltsByUserid)
		//根据活动id获取党员风采列表
		apifsdj.GET("/hlts/list", fsdj.GetStudyHlts)

		//签到
		apifsdj.GET("/signin", fsdj.StudySignin)
		//查看党员签到情况
		apifsdj.GET("/signins", fsdj.GetSigninsByUserid)
	}

	//======================== 外网H5——烟台税图_以地控税&&风险直推 ==========================//
	apiydks := r.Group("/api/ydks")
	apiydks.Use(ydksjwt.Check())
	{
		//内网数据文件下载路径
		apiydks.StaticFS("/inner/file", http.Dir(ydksrv.GetYdksFullPath()))
		//内网生成数据文件
		apiydks.GET("/inner/datafile", ydks.GenDataFile)
		//内网发送待办任务
		apiydks.POST("/inner/workrecord", ydks.Workrecord)
		//内网更新待办任务
		apiydks.POST("/inner/updworkrecord", ydks.UpdWorkrecord)
		//内网获取待办任务推送及更新情况
		apiydks.GET("/inner/workrecords", ydks.GetWorkrecords)
		//内网上传文件
		apiydks.POST("/inner/upload", ydks.YdksUploadFile)
		//内网获取外网业务数据
		apiydks.GET("/inner/outer_data", ydks.GetOuterData)

		//获取部门列表
		apiydks.GET("/inner/depts", dingtalk.GetDepartmentByParentID)
		//获取部门用户列表
		apiydks.GET("/inner/users", dingtalk.GetUserByDepartmentID)
		//模糊查询用户
		apiydks.GET("/inner/user/mc", dingtalk.GetUserByMc)

		//外网文件下载路径
		apiydks.StaticFS("/outer/file", http.Dir(ydksrv.GetYdksFullPath()))
		//外网接收业务数据
		apiydks.POST("/outer/recv", ydks.Recv)
		//外网获取已推送待办任务
		apiydks.GET("/outer/workrecords", ydks.GetWorkrecords)
	}
	return r
}
