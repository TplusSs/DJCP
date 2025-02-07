package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "time"

    "ZEDB/utils"
    "ZEDB/cmd"
)

// ActivationInfo 保存激活信息的结构体
type ActivationInfo struct {
    EncryptedFeatureCode    string    `json:"encrypted_feature_code"`
    EncryptedActivationCode string    `json:"encrypted_activation_code"`
    ActivationTime          time.Time `json:"activation_time"`
    ValidDays               int       `json:"valid_days"`
}

// saveActivationInfo 保存激活信息到文件
func saveActivationInfo(info ActivationInfo) error {
    data, err := json.Marshal(info)
    if err!= nil {
        return err
    }
    return ioutil.WriteFile("activation_info.json", data, 0644)
}

// loadActivationInfo 从文件加载激活信息
func loadActivationInfo() (ActivationInfo, error) {
    var info ActivationInfo
    data, err := ioutil.ReadFile("activation_info.json")
    if err!= nil {
        return info, err
    }
    err = json.Unmarshal(data, &info)
    return info, err
}

// checkActivation 检查激活信息是否有效
func checkActivation(info ActivationInfo) bool {
    if utils.ValidateActivationCode(info.EncryptedFeatureCode, info.EncryptedActivationCode) {
        expirationTime := info.ActivationTime.AddDate(0, 0, info.ValidDays)
        if time.Now().Before(expirationTime) {
            return true
        }
    }
    return false
}

func main() {
    // 尝试加载激活信息
    info, err := loadActivationInfo()
    if err == nil && checkActivation(info) {
        fmt.Println("程序已激活，可以使用。")
        // 原有的程序逻辑
        start := time.Now()
        cmd.Execute()
        end := time.Now().Sub(start)
        fmt.Printf("[*] 任务结束,耗时: %s\n", end)
        return
    }

    // 生成特征码
    encryptedFeatureCode, err := utils.GenerateFeatureCode()
    if err!= nil {
        fmt.Println("生成特征码失败:", err)
        return
    }

    // 输出加密后的特征码
    fmt.Println("加密后的特征码:", encryptedFeatureCode)
    fmt.Println("请将特征码发送给开发者获取激活码。")

    // 检查之前是否有激活信息但已过期，若有则提示重新激活
    if err == nil {
        fmt.Println("激活码已过期，请重新激活。")
    }

    fmt.Print("请输入激活码: ")
    var encryptedActivationCodeStr string
    fmt.Scanln(&encryptedActivationCodeStr)

    // 验证激活码
    if utils.ValidateActivationCode(encryptedFeatureCode, encryptedActivationCodeStr) {
        fmt.Println("激活成功，可以使用程序。")
        // 保存激活信息
        activationInfo := ActivationInfo{
            EncryptedFeatureCode:    encryptedFeatureCode,
            EncryptedActivationCode: encryptedActivationCodeStr,
            ActivationTime:          time.Now(),
            ValidDays:               30, // 假设有效期为 30 天，可根据实际情况修改
        }
        err := saveActivationInfo(activationInfo)
        if err!= nil {
            fmt.Println("保存激活信息失败:", err)
        }

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