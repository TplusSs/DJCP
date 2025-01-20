package web

import (
	"fmt"
	"ZEDB/global"
	"ZEDB/run"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ZEDBDj 模拟定级
func ZEDBDj(c *gin.Context) {
	c.HTML(http.StatusOK, "dj.html", gin.H{"ImageURL": global.ImageURL})
}

// ZEDBHome ZEDBIndex 单主机首页
func ZEDBHome(c *gin.Context) {
	c.HTML(http.StatusOK, "ZEDBHome.html", gin.H{"Version": global.Version, "ImageURL": global.ImageURL})

}

// ZEDBIndex 单主机首页
func ZEDBIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"Version": global.Version, "ImageURL": global.ImageURL})
}

// ZEDBIndexFile 多主机首页
func ZEDBIndexFile(c *gin.Context) {
	c.HTML(http.StatusOK, "indexFile.html", gin.H{"Version": global.Version, "ImageURL": global.ImageURL})
}

// ZEDBSubmitFile 先获取上传的文件；判断格式是否为xlsx；转换为临时txt文件；通过share函数执行多主机模式
func ZEDBSubmitFile(c *gin.Context) {
	filename, err := c.FormFile("uploaded-file")
	if err != nil {
		ZEDBErrorhtml("error", "上传文件失败了哦！选择文件了吗？", c)
	}
	if filepath.Ext(filename.Filename) != ".xlsx" {
		ZEDBErrorhtml("error", "文件只允许上传xlsx格式哦！！", c)
	}
	tempfilenamexlsx := fmt.Sprintf("%v.xlsx", time.Now().Unix())
	tempfilenametxt := fmt.Sprintf("%v.txt", time.Now().Unix())
	tempfilenamezip := fmt.Sprintf("%v.zip", time.Now().Unix())
	//退出时删除临时文件
	defer func() {
		os.Remove(tempfilenamexlsx)
		os.Remove(tempfilenametxt)
		os.Remove(tempfilenamezip)
	}()
	//保存上传文件
	err = c.SaveUploadedFile(filename, tempfilenamexlsx)
	if err != nil {
		ZEDBErrorhtml("error", "上传xlsx文件保存失败！", c)
	}
	if CreateTmpTxt(tempfilenamexlsx, tempfilenametxt) {
		mode := c.PostForm("mode")
		var allserver []Service  //Service结构体存储所有主机记录
		var alliplist []string   //预期成功的主机
		var successlist []string //实际成功的主机

		filedata, _ := os.ReadFile(tempfilenametxt)
		for _, s := range strings.Split(string(filedata), "\n") {
			if strings.Count(s, "~~") != 4 {
				continue
			}
			namesplit := strings.Split(s, "~~")
			//增加到所有主机切片中
			allserver = append(allserver, Service{Name: namesplit[0], User: namesplit[2], Ip: namesplit[1], Port: namesplit[4], Time: time.Now().Format(time.DateTime), Type: mode, Status: Failed})
			//增加保存的文件路径名称到切片中
			apendname := filepath.Join(global.Succpath, mode, fmt.Sprintf("%s_%s.log", namesplit[0], namesplit[1]))
			if mode == "MySQL" || mode == "Redis" || mode == "pgsql" || mode == "Linux" || mode == "oracle" {
				apendname = strings.ReplaceAll(apendname, ".log", ".html")
			}
			//如果是网络设备：拼接目录时需要更改为Route
			if mode == "h3c" || mode == "huawei" {
				apendname = filepath.Join(global.Succpath, "Route", fmt.Sprintf("%s_%s.log", namesplit[0], namesplit[1]))
			}
			alliplist = append(alliplist, apendname)
			//删除同名主机记录
			os.Remove(apendname)
		}
		switch mode {
		case "h3c":
			run.Rourange(tempfilenametxt, "~~", run.Defroutecmd) //运行H3C多主机模式
		case "huawei":
			run.Rourange(tempfilenametxt, "~~", run.DefroutecmdHuawei) //运行huawei多主机模式，待测试
		default:
			run.Rangefile(tempfilenametxt, "~~", mode) //运行多主机模式
		}
		//如果文件文件则写入到成功主机列表中
		for _, v := range alliplist {
			if global.PathExists(v) {
				successlist = append(successlist, v)
			}
		}
		defer FileAppendJson(successlist, allserver)
		if len(successlist) == 0 {
			ZEDBErrorhtml("error", fmt.Sprintf("%d个主机全部执行失败了哦!", len(alliplist)), c)
			c.Abort()
			return
		}
		// 退出时如果sava=false，则删除文件
		defer func() {
			if !save {
				for _, s := range successlist {
					os.Remove(s)
				}
			}
		}()
		err := CreateZipFromFiles(successlist, tempfilenamezip)
		if err != nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			ZEDBErrorhtml("error", "打包成zip包失败了！", c)
			c.Abort()
			return
		}
		//返回压缩包
		sendFile(tempfilenamezip, c)
	}
}

// ZEDBSubmit 单次提交任务
func ZEDBSubmit(c *gin.Context) {
	name, ip, user, passwd, port, mode, down := c.PostForm("name"), c.PostForm("ip"), c.PostForm("user"), c.PostForm("password"), c.PostForm("port"), c.PostForm("run_mode"), c.PostForm("down")
	savefilename := fmt.Sprintf("%s_%s.log", name, ip)                //保存的文件夹名：名称_ip.log
	successfile := filepath.Join(global.Succpath, mode, savefilename) //保存的完整路径
	//如果是网络设备：拼接目录时需要更改为Route
	if mode == "h3c" || mode == "huawei" {
		successfile = filepath.Join(global.Succpath, "Route", savefilename) //网络设备模式下的完整路径
	}
	if mode == "MySQL" || mode == "Redis" || mode == "pgsql" || mode == "Linux" || mode == "oracle" {
		successfile = strings.ReplaceAll(successfile, ".log", ".html")
	}

	if global.PathExists(successfile) {
		WriteJSONToHistory(Service{name, ip, user, port, mode, time.Now().Format(time.DateTime), Failed})
		ZEDBErrorhtml("error", "保存的文件中有重名文件，更换一个吧客官~", c)
		return
	}

	switch mode {
	case "h3c":
		for _, cmd := range run.Defroutecmd {
			run.Routessh(successfile, ip, user, passwd, port, cmd)
		}
	case "huawei":
		for _, cmd := range run.DefroutecmdHuawei {
			run.Routessh(successfile, ip, user, passwd, port, cmd)
		}
	default: //其他模式统一函数传参
		run.Onlyonerun(fmt.Sprintf("%s~~~%s~~~%s~~~%s~~~%s", name, ip, user, passwd, port), "~~~", mode)
	}
	if global.PathExists(successfile) {
		//如果不保存文件，文件返回后删除
		defer func() {
			if !save {
				os.Remove(successfile)
			}
		}()
		//down下载文件,preview预览文件
		if down == "down" {
			c.Header("Content-Description", "File Transfer")
			c.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf(fmt.Sprintf("%s_%s(%s).log", name, ip, mode)))
			if mode == "MySQL" || mode == "Redis" || mode == "pgsql" || mode == "Linux" || mode == "oracle" {
				c.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf(fmt.Sprintf("%s_%s(%s).html", name, ip, mode)))
			}
			c.Header("Content-Type", "application/octet-stream")
		}
		//返回文件
		WriteJSONToHistory(Service{name, ip, user, port, mode, time.Now().Format(time.DateTime), Success})

		c.File(successfile)
	} else {
		WriteJSONToHistory(Service{name, ip, user, port, mode, time.Now().Format(time.DateTime), Failed})
		ZEDBErrorhtml("error", "失败了哦客官~", c)
	}
}

// ZEDBMondeFileGet 返回模板文件
func ZEDBMondeFileGet(c *gin.Context) {
	//如果本地没有模板文件则生成一个
	if !global.PathExists(global.XlsxTemplateName) && !CreateTemplateXlsx() {
		ZEDBErrorhtml("error", "模板文件生成失败!", c)
	}
	// 返回模板文件
	sendFile(global.XlsxTemplateName, c)
}

// ZEDBErrorhtml 返回提示页面
func ZEDBErrorhtml(status, errbody string, c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", gin.H{"Status": status, "Message": errbody, "ImageURL": global.ImageURL})
}

// sendFile 发送文件
func sendFile(name string, c *gin.Context) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Header("Content-Type", "application/octet-stream")
	c.File(name)
}

// ZEDBUpdate 检查更新
func ZEDBUpdate(c *gin.Context) {
	release, err := global.CheckForUpdate()
	if err != nil {
		ZEDBErrorhtml("error", "获取最新版本失败,网络不好吧亲～", c)
		c.Abort()
		return
	}
	if release.TagName == global.Version {
		ZEDBErrorhtml("success", "非常好！当前是最新版本哦~", c)
		c.Abort()
		return
	}
	ZEDBErrorhtml("update", fmt.Sprintf("<a href='https://github.com/selinuxG/ZEDB-cli/releases' target='_blank'>当前版本为:%s,最新版本为:%s,点击此处进行更新！</a>", global.Version, release.TagName), c)

}

// ZEDBHistory 历史记录
func ZEDBHistory(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	allserver, err := parseJSONFile()
	if err != nil {
		ZEDBErrorhtml("nil", "不存在历史记录哦~", c)
		c.Abort()
		return
	}
	type ServerData struct {
		Id     int
		Name   string
		IP     string
		User   string
		Port   string
		Mode   string
		Time   string
		Status string
	}
	var dataSlice []ServerData

	// 倒序循环遍历切片
	for i := len(allserver) - 1; i >= 0; i-- {
		id := len(allserver) - i
		server := allserver[i]

		// 创建并初始化 ServerData 结构体
		data := ServerData{
			Id:     id,
			Name:   server.Name,
			IP:     server.Ip,
			User:   server.User,
			Port:   server.Port,
			Mode:   server.Type,
			Time:   server.Time,
			Status: server.Status,
		}
		// 将数据添加到 dataSlice
		dataSlice = append(dataSlice, data)
	}
	c.HTML(http.StatusOK, "ZEDBHistoryIndex.html", gin.H{"Data": dataSlice, "ImageURL": global.ImageURL})
}

// FileAppendJson 将成功主机对比allserver主机，写入到json文件中
// success = 采集完成目录\mode\name_ip.log,
// 留存了个bug，正常使用不会触发，所以不修复了没意义。
func FileAppendJson(success []string, allserver []Service) {
	for _, srv := range allserver {
		for _, succip := range success {
			if strings.Count(succip, srv.Ip) > 0 {
				srv.Status = Success
				break
			}
		}
		WriteJSONToHistory(srv)
	}
}
