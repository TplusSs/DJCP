package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "net"
    "strings"
    "time"
    "bytes"        // 添加 bytes 包的导入
)

// 加密密钥，实际使用中应使用更安全的密钥管理方式
const encryptionKey = "1234567890123456"

// GetMACAddress 获取本机 MAC 地址
func GetMACAddress() (string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
            addrs, err := iface.Addrs()
            if err != nil {
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
                if hw != "" {
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
    if err != nil {
        return "", err
    }
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    data := mac + currentTime
    hash := sha256.Sum256([]byte(data))
    featureCode := hex.EncodeToString(hash[:])

    // 加密特征码
    encryptedFeatureCode, err := encrypt([]byte(featureCode), []byte(encryptionKey))
    if err != nil {
        return "", err
    }
    return encryptedFeatureCode, nil
}

// ActivationCode 激活码结构体
type ActivationCode struct {
    FeatureCode     string    `json:"feature_code"`
    ActivationTime  time.Time `json:"activation_time"`
    ValidPeriod     string    `json:"valid_period"`
}

// ValidateActivationCode 验证激活码
func ValidateActivationCode(encryptedFeatureCode string, encryptedActivationCodeStr string) bool {
    // 解密特征码
    decryptedFeatureCode, err := decrypt(encryptedFeatureCode, []byte(encryptionKey))
    if err != nil {
        return false
    }

    // 解密激活码
    decryptedActivationCode, err := decrypt(encryptedActivationCodeStr, []byte(encryptionKey))
    if err != nil {
        return false
    }

    var activationCode ActivationCode
    err = json.Unmarshal([]byte(decryptedActivationCode), &activationCode)
    if err != nil {
        return false
    }
    if activationCode.FeatureCode != string(decryptedFeatureCode) {
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

// 加密函数
func encrypt(data []byte, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    data = PKCS7Padding(data, aes.BlockSize)
    ciphertext := make([]byte, aes.BlockSize+len(data))
    mode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
    mode.CryptBlocks(ciphertext[aes.BlockSize:], data)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 解密函数
func decrypt(ciphertext string, key []byte) ([]byte, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return nil, err
    }
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    mode := cipher.NewCBCDecrypter(block, key[:aes.BlockSize])
    mode.CryptBlocks(data[aes.BlockSize:], data[aes.BlockSize:])
    data = PKCS7UnPadding(data[aes.BlockSize:])
    return data, nil
}

// PKCS7 填充
func PKCS7Padding(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
}

// PKCS7 去填充
func PKCS7UnPadding(data []byte) []byte {
    length := len(data)
    unpadding := int(data[length-1])
    return data[:(length - unpadding)]
}