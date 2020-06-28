package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret      string
	PrefixUrl      string
	AppPrefixUrl   string
	LoginUrl       string
	TokenTimeout   int64
	DingtalkMsgUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type Dingtalk struct {
	CorpID       string
	OapiHost     string
	CallBackHost string
	Token        string
	AesKey       string
}

var DingtalkSetting = &Dingtalk{}

type MsgApp struct {
	AgentID   string
	AppKey    string
	AppSecret string
	Domain    string
}

var MsgAppSetting = &MsgApp{}

type EApp struct {
	AgentID   string
	AppKey    string
	AppSecret string
	Domain    string
}

var EAppSetting = &EApp{}

type YdksApp struct {
	AgentID   string
	AppKey    string
	AppSecret string
	Domain    string

	YdksSavePath string
	FlowLimit    int
}

var YdksAppSetting = &YdksApp{}

type FsdjEapp struct {
	AgentID   string
	AppKey    string
	AppSecret string
	Domain    string

	FsdjSavePath string
}

var FsdjEappSetting = &FsdjEapp{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("dingtalk", DingtalkSetting)
	mapTo("msgapp", MsgAppSetting)
	mapTo("eapp", EAppSetting)
	mapTo("ydks", YdksAppSetting)
	mapTo("fsdj", FsdjEappSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
