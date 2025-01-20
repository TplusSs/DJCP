package main

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "fmt"
    "os"
    "time"
    "ZEDB/utils"
    "ZEDB/cmd"
)

// 加密密钥，实际使用中应使用更安全的密钥管理方式
const encryptionKey = "1234567890123456"

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

func main() {
    // 生成特征码
    featureCode, err := utils.GenerateFeatureCode()
    if err != nil {
        fmt.Println("生成特征码失败:", err)
        return
    }
    fmt.Println("特征码:", featureCode)
    fmt.Println("请将特征码发送给开发者获取激活码。")
    fmt.Print("请输入激活码: ")
    var encryptedActivationCodeStr string
    fmt.Scanln(&encryptedActivationCodeStr)

    // 解密激活码
    decryptedActivationCode, err := decrypt(encryptedActivationCodeStr, []byte(encryptionKey))
    if err != nil {
        fmt.Println("解密激活码失败:", err)
        os.Exit(1)
    }

    if utils.ValidateActivationCode(featureCode, string(decryptedActivationCode)) {
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