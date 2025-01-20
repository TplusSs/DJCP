package run

func pgsqlhtml() string {
	return `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PostgreSQL安全策略核查</title>
	<link rel="icon" href="https://s1.ax1x.com/2023/07/19/pC7B5sx.jpg" sizes="16x16">
    <style>

        body {
            display: grid;
            grid-template-columns: 1fr 200px;
            gap: 10px;
            font-family: Arial, sans-serif;
            position: relative;
        }

        table {
    		border-collapse: collapse;
   		 	margin-bottom: 20px;
    		width: 100%;
    		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    		table-layout: fixed;
    		word-wrap: break-word;
        }

        th,
        td {
            border: 1px solid #ddd;
            padding: 15px;
            text-align: left;
        }

        th {
            background-color: #007BFF;
            color: white;
            font-weight: bold;
        }

        tr:nth-child(even) {
            background-color: #f9f9f9;
        }

        tr:hover {
            background-color: #e6f2ff;
        }

        .watermark {
            font-size: 36px;
            color: rgba(128, 128, 128, 0.2);
            position: absolute;
            z-index: -1;
            transform: rotate(-30deg);
        }

        #toc {
            position: fixed;
            top: 20px;
            right: 30px;
            padding-left: 10px;
            background-color: #f8f9fa;
            padding: 10px;
            border: 1px solid #dee2e6;
            border-radius: 5px;
            box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
        }

        #toc ul {
            list-style-type: none;
            padding: 0;
        }

        #toc a {
            text-decoration: none;
            color: #333;
            display: block;
        }

        #toc a:hover {
            color: #007BFF;
        }
		.watermark {
			font-size: 36px;
			color: rgba(128, 128, 128, 0.2);
			position: absolute;
			z-index: 1000;
			transform: rotate(-30deg);
    	}
        pre {
            background-color: #f8f9fa;
            border: 1px solid #dee2e6;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
            font-family: "Courier New", Courier, monospace;
			white-space: pre-wrap;
			word-break: break-word;
        }
		.permissions {
			width: 350px;
			white-space: nowrap;
			overflow-x: auto;
		}
    </style>

<body>

    <div id="content">
        <center><h1>替换名称_PostgreSQL安全策略核查</h1></center>

        <h2 id="version">版本信息</h2>
        <table>
            <thead>
                <tr>
                    <th>版本信息</th>
                </tr>
            </thead>
            <tbody>
                版本信息详细信息
            </tbody>
        </table>

        <h2 id="userinfo">用户信息</h2>
        <table>
            <thead>
                <tr>
                    <th>角色</th>
                    <th>ID</th>
                    <th>超级管理员</th>
                    <th>过期时间</th>
                    <th>口令</th>
                    <th>可登录</th>
                    <th>可创建数据库</th>
                    <th>可创建角色</th>
                    <th>可继承</th>
                </tr>
            </thead>
            <tbody>
                	user详细信息
            </tbody>
        </table>

        <h2 id="roleinfo">角色信息</h2>
        <table>
            <thead>
                <tr>
                    <th>角色</th>
                    <th>隶属于</th>
                </tr>
            </thead>
            <tbody>
                	role详细信息
            </tbody>
        </table>

        <h2 id="libraries">已安装插件</h2>
        <table>
            <thead>
                <tr>
                    <th>插件名称</th>
                </tr>
            </thead>
            <tbody>
                	插件详细信息
            </tbody>
        </table>

        <h2 id="ipaddr">远程监听网卡</h2>
        <table>
            <thead>
                <tr>
                    <th>监听网卡地址</th>
                </tr>
            </thead>
            <tbody>
                	网卡地址详细信息
            </tbody>
        </table>

        <h2 id="TLS">最低TLS版本</h2>
        <table>
            <thead>
                <tr>
                    <th>版本</th>
                </tr>
            </thead>
            <tbody>
                	最低支持的TLS版本
            </tbody>
        </table>

        <h2 id="log">日志信息</h2>
        <table>
            <thead>
                <tr>
                    <th>错误日志状态</th>
                    <th>错误日志记录级别</th>
                    <th>日志文件存储目录</th>
                    <th>日志文件命令格式</th>
                    <th>查询语句开启状态</th>
                    <th>用户登录记录开启状态</th>
                    <th>用户登出记录开启状态</th>
                    <th>日志内容记录字段</th>
                    <th>日志发送类型</th>
                </tr>
            </thead>
            <tbody>
                	log信息
            </tbody>
        </table>

    </div>

    <div id="toc">
        <h3>目录</h3>
        <ul>
            <li><a href="#version">版本信息</a></li>
            <li><a href="#userinfo">用户信息</a></li>
            <li><a href="#roleinfo">角色信息</a></li>
            <li><a href="#libraries">插件信息</a></li>
            <li><a href="#ipaddr">监听网卡</a></li>
            <li><a href="#TLS">TLS最低版本</a></li>
            <li><a href="#log">日志信息</a></li>
        </ul>
    </div>
	<div id="watermark"></div>
</body>

</html>
<script>
    const watermarkNum = 30 // 生成水印数量
    build()

    function build(){
        for(var i = 0; i < watermarkNum; i++){
            addWatermark(i);
        }
    }

    function addWatermark(i){
        var watermark = document.getElementById("watermark");
        const top = i
        const left = random();
        const  html = '<div class="watermark" style="top: '+(top/watermarkNum)*100+'%; left: '+left+'%;">杭州中尔网络科技有限公司</div>'
        watermark.insertAdjacentHTML('afterend',html);
    }

    function random(){
       return Math.floor(Math.random() * 70) ;
    }
</script>
`
}
