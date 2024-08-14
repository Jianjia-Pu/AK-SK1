package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func signRequest(randomString, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(randomString))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println("Generated Signature:", signature)
	return signature
}

func main() {
	accessKey := "qwerty"
	secretKey := "qazwsx"
	method := "GET"
	url := "http://localhost:8080/resource"

	// 生成一个随机字符串
	randomString, err := generateRandomString(16)
	if err != nil {
		fmt.Println("Error generating random string:", err)
		return
	}

	// 生成签名
	signature := signRequest(randomString, secretKey)

	// 发送请求
	authHeader := fmt.Sprintf("API %s:%s:%s", accessKey, randomString, signature)
	fmt.Println("Authorization Header:", authHeader)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(responseBody))
}
