package routers

import (
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

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.CORSMiddleware())
	//port := strconv.Itoa(setting.ServerSetting.HttpPort + 1)

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/upload", api.UploadImage)
	r.POST("/uploads/file", api.UploadFile)

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{

		//生成海报
		//apiv1.POST("/poster/generate", v1.GeneratePoster)
		apiv1.POST("/login", dingtalk.Login)
	}
	return r
}
