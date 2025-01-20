package web

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"ZEDB/global"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"
)

var save bool

//go:embed favicon.ico
var faviconFS embed.FS

//go:embed template/*
var f embed.FS

//go:embed img/*
var imgf embed.FS

func Start(cmd *cobra.Command, args []string) {

	if !global.PathExists("cert/cert.pem") || !global.PathExists("cert/cert.key") {
		CreateCert()
	}

	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetString("port")
	save, _ = cmd.Flags().GetBool("save")
	global.ImageURL = fmt.Sprintf("https://%s:%s/bgi", ip, port)
	r := gin.Default()
	tmpl := template.Must(template.New("").ParseFS(f, "template/*"))
	r.SetHTMLTemplate(tmpl)
	r.Use(func(c *gin.Context) {
		c.Header("author", "TPlus")
		c.Next()
	})
	r.NoRoute(func(c *gin.Context) {
		ZEDBErrorhtml("404", "sorry~请求不存在哦!", c)
	})
	r.GET("/favicon.ico", faviconHandler) //路由图标
	r.GET("/bgi", bgiHandler)             // 背景图片
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/zedb/home")
	})

	ZEDB := r.Group("/zedb")
	{
		ZEDB.GET("/home", ZEDBHome)              //首页
		ZEDB.GET("/index", ZEDBIndex)            //单主机index
		ZEDB.GET("/indexfile", ZEDBIndexFile)    //多主机index
		ZEDB.GET("/modefile", ZEDBMondeFileGet)  //返回模板文件
		ZEDB.POST("/submit", ZEDBSubmit)         //提交单主机任务
		ZEDB.POST("/submitfile", ZEDBSubmitFile) //提交多主机任务
		ZEDB.GET("/history", ZEDBHistory)        //历史记录
		ZEDB.GET("/update", ZEDBUpdate)          //检查更新
		ZEDB.GET("/dj", ZEDBDj)                  //模拟定级首页
	}
	// Windows、Mac下在默认浏览器中打开网页
	go func() {
		if runtime.GOOS == "windows" {
			cmd := exec.Command("cmd", "/C", fmt.Sprintf("start https://%s:%s/zedb/home", ip, port))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error opening the browser:", err)
			}
		}
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("open", fmt.Sprintf("https://%s:%s/zedb/home", ip, port))
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error opening the browser:", err)
			}
		}
	}()
	// 启动gin
	r.RunTLS(ip+":"+port, "cert/cert.pem", "cert/key.pem")
}

func faviconHandler(c *gin.Context) {
	c.FileFromFS("favicon.ico", http.FS(faviconFS))
}

func bgiHandler(c *gin.Context) {
	c.FileFromFS("img/2.png", http.FS(imgf))
}
