package main

// import (
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/base64"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// )

// func signRequest(method, url, headers, body, secretKey string) string {
// 	// 规范化请求
// 	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s", method, url, headers, body)

// 	// 生成HMAC签名
// 	mac := hmac.New(sha256.New, []byte(secretKey))
// 	mac.Write([]byte(canonicalRequest))
// 	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

// 	return signature
// }

// func main() {
// 	accessKey := "YOUR_ACCESS_KEY"
// 	secretKey := "YOUR_SECRET_KEY"
// 	method := "GET"
// 	url := "http://localhost:8080/resource"
// 	headers := "Host:localhost:8080"
// 	body := ""

// 	// 生成签名
// 	signature := signRequest(method, url, headers, body, secretKey)

// 	// 发送请求
// 	authHeader := fmt.Sprintf("API %s:%s", accessKey, signature)
// 	req, err := http.NewRequest(method, url, strings.NewReader(body))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return
// 	}
// 	req.Header.Set("Authorization", authHeader)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response:", err)
// 		return
// 	}

// 	fmt.Println("Response:", string(responseBody))
// }
