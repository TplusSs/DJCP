package main

import (
    "encoding/json"
    "fmt"
    "time"
    "ZEDB/utils"
)

// ActivationCode 激活码结构体
type ActivationCode struct {
    FeatureCode     string    `json:"feature_code"`
    ActivationTime  time.Time `json:"activation_time"`
    ValidDays       int       `json:"valid_days"`
}

// GenerateActivationCode 根据特征码生成激活码
func GenerateActivationCode(encryptedFeatureCode string, validDays int) (string, error) {
    activationCode := ActivationCode{
        FeatureCode:     encryptedFeatureCode,
        ActivationTime:  time.Now(),
        ValidDays:       validDays,
    }
    activationCodeBytes, err := json.Marshal(activationCode)
    if err!= nil {
        return "", err
    }

    // 使用统一的加密函数
    encryptedActivationCode, err := utils.Encrypt(activationCodeBytes, []byte(utils.EncryptionKey))
    if err!= nil {
        return "", err
    }

    return encryptedActivationCode, nil
}

func main() {
    fmt.Print("请输入特征码: ")
    var encryptedFeatureCode string
    fmt.Scanln(&encryptedFeatureCode)

    fmt.Print("请输入有效期天数 (数字，1 代表 1 天): ")
    var validDays int
    fmt.Scanln(&validDays)

    activationCode, err := GenerateActivationCode(encryptedFeatureCode, validDays)
    if err!= nil {
        fmt.Println("生成激活码失败:", err)
        return
    }

    fmt.Println("加密后的激活码:", activationCode)
}