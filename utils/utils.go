// utils/password.go

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-rest/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// decode token
var jwtKey = []byte("merciku")

// Customer Session untuk redis
type CustomerSession struct {
	UserID      int64  `json:"user_id"`
	Code        string `json:"code"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Chapter     int    `json:"chapter"`
	ChapterCode string `json:"chapter_code"`
	Token       string `json:"token"`
	Device      string `json:"device"`
	LoginAt     string `json:"login_at"`
}

// Claims struct untuk memetakan data token
type Claims struct {
	UserID       int    `json:"userid"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Email        string `json:"username"`
	Phone        string `json:"phone"`
	Premium      int8   `json:"premium"`
	Chapter      int    `json:"chapter"`
	Chapter_code string `json:"chapter_code"`
	jwt.StandardClaims
}

// auth middleware

func Token(c *gin.Context) string {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header missing"})
		return ""
	}

	// Pisahkan Bearer dan token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(401, gin.H{"error": "Invalid Authorization header format"})
		return ""
	}

	tokenStr := parts[1]

	// Bersihkan karakter yang tidak perlu
	tokenStr = strings.Trim(tokenStr, `"`)
	tokenStr = strings.TrimSpace(tokenStr)

	return tokenStr
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := Token(c)
		// Panggil fungsi otentikasi
		if !Otentikasi(tokenStr) {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Token mismatchxx",
			})
			return
		}
		c.Next()
	}
}

// ini fungsi otentikasi
func Otentikasi(tokenString string) bool {

	claims, err := DecodeToken(tokenString)
	if err != nil {
		fmt.Println("Token decode error:", err)
		return false
	}

	// fmt.Println("Userid:", claims.UserID)

	key := fmt.Sprintf("session:%d", claims.UserID)

	val, err := config.RDB.Get(config.Ctx, key).Result()
	if err != nil {
		// fmt.Println("Redis error:", err)
		// fmt.Println("Redis key:", key)
		// fmt.Println("Redis error:", err)
		return false
	}

	// Unmarshal JSON
	var session CustomerSession
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		// fmt.Println("JSON unmarshal error:", err)
		return false
	}

	// Ambil token saja
	redisToken := session.Token

	// fmt.Println("Redis token:", val)

	if redisToken != tokenString {
		fmt.Println("Token mismatch")
		// fmt.Println("Token Redis : ", val)
		// fmt.Println("Token Header : ", tokenString)
		return false
	}

	// return true
	return true
}

// DecodeToken mendekode JWT token
func DecodeToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	// Parse token dengan secret key
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Validasi token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Cek apakah token sudah kedaluwarsa
	// if claims.ExpiresAt < time.Now().Unix() {
	// 	return nil, errors.New("token has expired")
	// }

	return claims, nil
}

// GetLocalTime mengembalikan waktu saat ini dalam zona waktu lokal
func GetLocalTime() time.Time {
	// Ganti dengan nama zona waktu yang sesuai
	loc, err := time.LoadLocation("Asia/Jakarta") // Misalnya untuk WIB
	if err != nil {
		panic(err) // Atau tangani kesalahan sesuai kebutuhan Anda
	}
	return time.Now().In(loc) // Kembalikan waktu saat ini dalam zona waktu lokal
}

// HashPassword menerima password dan mengembalikan hash-nya
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassword membandingkan password dengan hash yang disimpan
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	} else {
		return true // jika error sama dengan null maka true (tidak ada error)
	}
	// return err == nil
	// return false

	// fmt.Println("Password input:", password)
	// fmt.Println("Hash DB:", hash)
	// return false
}
