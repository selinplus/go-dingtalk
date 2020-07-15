package fsdj

import (
	"github.com/gin-gonic/gin"
	"github.com/selinplus/go-dingtalk/pkg/util"
	"html/template"
)

func Home(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}
