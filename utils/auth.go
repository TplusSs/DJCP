package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    //"fmt"
    "net"
    "strings"
    "time"
    "bytes"
)

// EncryptionKey 加密密钥，实际使用中应使用更安全的密钥管理方式
const EncryptionKey = "1234567890123456"

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
    encryptedFeatureCode, err := Encrypt([]byte(featureCode), []byte(EncryptionKey))
    if err != nil {
        return "", err
    }
    //fmt.Printf("加密后的特征码: %s\n", encryptedFeatureCode)
    //fmt.Println("请将特征码发送给开发者获取激活码")
    return encryptedFeatureCode, nil
}

// ActivationCode 激活码结构体
type ActivationCode struct {
    FeatureCode     string    `json:"feature_code"`
    ActivationTime  time.Time `json:"activation_time"`
    ValidDays       int       `json:"valid_days"`
}

// ValidateActivationCode 验证激活码
func ValidateActivationCode(encryptedFeatureCode string, encryptedActivationCodeStr string) bool {
    // 解密特征码
    decryptedFeatureCode, err := Decrypt(encryptedFeatureCode, []byte(EncryptionKey))
    if err != nil {
        return false
    }

    // 解密激活码
    decryptedActivationCode, err := Decrypt(encryptedActivationCodeStr, []byte(EncryptionKey))
    if err != nil {
        return false
    }

    var activationCode ActivationCode
    err = json.Unmarshal([]byte(decryptedActivationCode), &activationCode)
    if err != nil {
        return false
    }

    // 新增: 检查 FeatureCode 是否是加密后的形式
    decryptedFeatureCodeFromActivationCode, err := Decrypt(activationCode.FeatureCode, []byte(EncryptionKey))
    if err != nil {
        return false
    }

    if string(decryptedFeatureCode) != string(decryptedFeatureCodeFromActivationCode) {
        return false
    }
    now := time.Now()
    expirationTime := activationCode.ActivationTime.AddDate(0, 0, activationCode.ValidDays)
    if now.After(expirationTime) {
        return false
    }
    return true
}

// Encrypt 加密函数
func Encrypt(data []byte, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    data = PKCS7Padding(data, aes.BlockSize)
    ciphertext := make([]byte, aes.BlockSize+len(data))
    mode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
    mode.CryptBlocks(ciphertext[aes.BlockSize:], data)
    encryptedStr := base64.StdEncoding.EncodeToString(ciphertext)
    return encryptedStr, nil
}

// Decrypt 解密函数
func Decrypt(ciphertext string, key []byte) ([]byte, error) {
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

// PKCS7Padding 填充函数
func PKCS7Padding(data []byte, blockSize int) []byte {
    padding := blockSize - len(data)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padtext...)
}

// PKCS7UnPadding 去填充函数
func PKCS7UnPadding(data []byte) []byte {
    length := len(data)
    unpadding := int(data[length-1])
    return data[:(length - unpadding)]
}