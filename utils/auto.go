package utils

import (
    "crypto/sha256"
    "encoding/json"
    "net"
    "strings"
    "time"
    "encoding/hex" // 导入 hex 包
)

// GetMACAddress 获取本机 MAC 地址
func GetMACAddress() (string, error) {
    interfaces, err := net.Interfaces()
    if err!= nil {
        return "", err
    }
    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp!= 0 && iface.Flags&net.FlagLoopback == 0 {
            addrs, err := iface.Addrs()
            if err!= nil {
                continue
            }
            for _, addr := range addrs {
                var ip net.IP
                switch v := addr.(type) {
                case *net.IPNet:
                    ip = v.IP
                case *net.IPAddr:
                    ip = v.IP
                }
                if ip == nil || ip.IsLoopback() {
                    continue
                }
                hw := iface.HardwareAddr.String()
                if hw!= "" {
                    return strings.ReplaceAll(hw, ":", ""), nil
                }
            }
        }
    }
    return "", nil
}

// GenerateFeatureCode 生成特征码
func GenerateFeatureCode() (string, error) {
    mac, err := GetMACAddress()
    if err!= nil {
        return "", err
    }
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    data := mac + currentTime
    hash := sha256.Sum256([]byte(data))
    // 使用 hex 包的 EncodeToString 函数
    return hex.EncodeToString(hash[:]), nil 
}

// ActivationCode 激活码结构体
type ActivationCode struct {
    FeatureCode   string    `json:"feature_code"`
    ActivationTime time.Time `json:"activation_time"`
    ValidPeriod   string    `json:"valid_period"`
}

// ValidateActivationCode 验证激活码
func ValidateActivationCode(featureCode string, activationCodeStr string) bool {
    var activationCode ActivationCode
    err := json.Unmarshal([]byte(activationCodeStr), &activationCode)
    if err!= nil {
        return false
    }
    if activationCode.FeatureCode!= featureCode {
        return false
    }
    now := time.Now()
    var expirationTime time.Time
    switch activationCode.ValidPeriod {
    case "week":
        expirationTime = activationCode.ActivationTime.AddDate(0, 0, 7)
    case "month":
        expirationTime = activationCode.ActivationTime.AddDate(0, 1, 0)
    case "year":
        expirationTime = activationCode.ActivationTime.AddDate(1, 0, 0)
    default:
        return false
    }
    return now.Before(expirationTime)
}