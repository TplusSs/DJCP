// FILEPATH: DJCP/generate_activation_code.go
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
    ValidPeriod     string    `json:"valid_period"`
}

// GenerateActivationCode 根据特征码生成激活码
func GenerateActivationCode(encryptedFeatureCode string, validPeriod string) (string, error) {
    activationCode := ActivationCode{
        FeatureCode:     encryptedFeatureCode,
        ActivationTime:  time.Now(),
        ValidPeriod:     validPeriod,
    }
    activationCodeBytes, err := json.Marshal(activationCode)
    if err != nil {
        return "", err
    }

    // 使用统一的加密函数
    encryptedActivationCode, err := utils.encrypt(activationCodeBytes, []byte(utils.encryptionKey))
    if err != nil {
        return "", err
    }

    return encryptedActivationCode, nil
}

func main() {
    fmt.Print("请输入特征码: ")
    var encryptedFeatureCode string
    fmt.Scanln(&encryptedFeatureCode)

    fmt.Print("请输入有效期 (week/month/year): ")
    var validPeriod string
    fmt.Scanln(&validPeriod)

    activationCode, err := GenerateActivationCode(encryptedFeatureCode, validPeriod)
    if err != nil {
        fmt.Println("生成激活码失败:", err)
        return
    }

    fmt.Println("加密后的激活码:", activationCode)
}