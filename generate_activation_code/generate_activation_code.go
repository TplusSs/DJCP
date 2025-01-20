// FILEPATH: DJCP/generate_activation_code.go
package main

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/tjfoc/gmsm/sm4"
)

// ActivationCode 激活码结构体
type ActivationCode struct {
    FeatureCode   string    `json:"feature_code"`
    ActivationTime time.Time `json:"activation_time"`
    ValidPeriod   string    `json:"valid_period"`
}

// GenerateActivationCode 根据特征码生成激活码
func GenerateActivationCode(featureCode string, validPeriod string) (string, error) {
    activationCode := ActivationCode{
        FeatureCode:   featureCode,
        ActivationTime: time.Now(),
        ValidPeriod:   validPeriod,
    }
    activationCodeBytes, err := json.Marshal(activationCode)
    if err != nil {
        return "", err
    }

    // SM4 加密
    key := []byte("1234567890abcdef") // 16 字节的密钥
    sm4Cipher, err := sm4.NewSm4Cipher(key)
    if err != nil {
        return "", err
    }
    encrypted := make([]byte, len(activationCodeBytes))
    sm4Cipher.Encrypt(encrypted, activationCodeBytes)

    return fmt.Sprintf("%x", encrypted), nil
}

func main() {
    fmt.Print("请输入特征码: ")
    var featureCode string
    fmt.Scanln(&featureCode)

    fmt.Print("请输入有效期 (week/month/year): ")
    var validPeriod string
    fmt.Scanln(&validPeriod)

    activationCode, err := GenerateActivationCode(featureCode, validPeriod)
    if err != nil {
        fmt.Println("生成激活码失败:", err)
        return
    }

    fmt.Println("加密后的激活码:", activationCode)
}