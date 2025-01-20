package main

import (
    "fmt"
    "os"
    "time"

    "ZEDB/utils"
    "ZEDB/cmd"
)

func main() {
    // 生成特征码
    encryptedFeatureCode, err := utils.GenerateFeatureCode()
    if err != nil {
        fmt.Println("生成特征码失败:", err)
        return
    }
    fmt.Println("加密后的特征码:", encryptedFeatureCode)
    fmt.Println("请将特征码发送给开发者获取激活码。")
    fmt.Print("请输入激活码: ")
    var encryptedActivationCodeStr string
    fmt.Scanln(&encryptedActivationCodeStr)

    // 验证激活码
    if utils.ValidateActivationCode(encryptedFeatureCode, encryptedActivationCodeStr) {
        fmt.Println("激活成功，可以使用程序。")
        // 原有的程序逻辑
        start := time.Now()
        cmd.Execute()
        end := time.Now().Sub(start)
        fmt.Printf("[*] 任务结束,耗时: %s\n", end)
    } else {
        fmt.Println("激活码无效或已过期。")
        // 处理激活失败的情况，如退出程序或限制功能
        os.Exit(1)
    }
}