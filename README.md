# Go Gin Web dingtalk

## How to run

* run on dmz:
```
go build

```
* run on app:
```
change main.go
func init() {
    ...
    //cron.Setup
}
func main(){
    ...
    //go func() {
    	//	time.Sleep(time.Second * 10)
    	//	dingtalk.RegCallbackInit()
    	//}()
    ...
}

go build

```

### Required

- Mysql


### Ready

Create a **dingtalk database** 

### Conf

You should modify `conf/app.ini`

```
[app]
PageSize = 10
JwtSecret = 
PrefixUrl = 
#internet app
AppPrefixUrl = 
TokenTimeout = 30
#MseesageToDingding Url
DingtalkMsgUrl = 

[server]
#debug or release
RunMode = debug // change to release
...

[database]
Type = mysql
User = root
Password =
Host = 
Name = 
TablePrefix = 

...
```

### Run
```
$ cd $GOPATH/src/go-dingtalk

$ go run main.go 
```

Project information and existing API

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/...              --> github.com/selinplus/go-dingtalk/routers/api/v1.... (4 handlers)
[GIN-debug] POST   /api/v1/...              --> github.com/selinplus/go-dingtalk/routers/api/v1.... (4 handlers)

Listening port is 4449
```

## Features

- RESTful API
- Gorm
- logging
- Jwt-go
- Gin
- App configurable
- Cron