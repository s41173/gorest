// utils/password.go

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-rest/config"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// decode token
// var jwtKey = []byte("merciku")

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
	Image       string `json:"image"`
	LoginAt     string `json:"login_at"`
}

// Claims struct untuk memetakan data token
type Claims struct {
	UserID       int64  `json:"user_id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Email        string `json:"username"`
	Phone        string `json:"phone"`
	Premium      int8   `json:"premium"`
	Chapter      int    `json:"chapter"`
	Chapter_code string `json:"chapter_code"`
	jwt.StandardClaims
}

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

// auth middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := Token(c)
		// Panggil fungsi otentikasi
		if !Otentikasi(tokenStr) {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Token mismatch",
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

	// fmt.Println("Otentikasi->Userid:", claims.UserID)

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
		fmt.Println("Token mismatch dari otentikasi")
		// fmt.Println("Token Redis : ", val)
		// fmt.Println("Token Header : ", tokenString)
		return false
	}

	return true
}

// DecodeToken mendekode JWT token
func DecodeToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	secret := os.Getenv("JWT_SECRET")
	// secret := "dodol"
	if secret == "" {
		return nil, errors.New("JWT_SECRET is empty")
	}

	// fmt.Println("Secret Key : ", secret)
	// fmt.Println("Token : ", tokenStr)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {

		// 🔒 Validasi algoritma (hindari attack)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
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

// HashPassword menerima password dan mengembalikan hash-nya
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type AppError struct {
	Code    int
	Message string
}

// Error implements error.
func (a *AppError) Error() string {
	panic("unimplemented")
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
