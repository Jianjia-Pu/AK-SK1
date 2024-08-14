package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type APIKey struct {
	gorm.Model
	AccessKey string `gorm:"primaryKey"`
	SecretKey string
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ExpiresAt time.Time
	Status    string `gorm:"default:active"`
}

type UsedRandomString struct {
	gorm.Model
	RandomString string    `gorm:"uniqueIndex"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./api_keys.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// 自动迁移表结构
	db.AutoMigrate(&APIKey{}, &UsedRandomString{})
}

func findSecretKey(accessKey string) (string, error) {
	var apiKey APIKey
	result := db.Where("access_key = ? AND status = 'active'", accessKey).First(&apiKey)
	if result.Error != nil {
		return "", result.Error
	}
	fmt.Println("Found Secret Key:", apiKey.SecretKey)
	return apiKey.SecretKey, nil
}

func verifySignature(randomString, signature, secretKey string) bool {
	// 生成HMAC签名
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(randomString))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func handler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "API" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	authParts := strings.SplitN(parts[1], ":", 3)
	if len(authParts) != 3 {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	accessKey := authParts[0]
	randomString := authParts[1]
	clientSignature := authParts[2]

	// 检查随机字符串是否已经被使用过
	var usedRandomString UsedRandomString
	result := db.Where("random_string = ?", randomString).First(&usedRandomString)
	if result.Error == nil {
		http.Error(w, "Random string already used", http.StatusUnauthorized)
		return
	}

	// 从数据库中获取 secretKey
	secretKey, err := findSecretKey(accessKey)
	if err != nil {
		http.Error(w, "Invalid access key", http.StatusUnauthorized)
		return
	}

	if !verifySignature(randomString, clientSignature, secretKey) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// 将随机字符串标记为已使用
	usedRandomString = UsedRandomString{RandomString: randomString}
	db.Create(&usedRandomString)

	// 验证通过，返回资源
	w.Write([]byte("Resource data"))
}

func main() {
	initDB()
	http.HandleFunc("/resource", handler)
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
