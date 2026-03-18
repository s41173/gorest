package services

import (
	"fmt"
	"os"
	"time"

	"go-rest/config"
	"go-rest/models"
	"go-rest/utils"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(username, password string) (string, error) {
	var user models.Customer

	err := config.DB.
		Where("(email = ? OR phone1 = ?) AND status = 1 AND deleted IS NULL", username, username).
		Take(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("db error: %v", err)
	}

	if user.Verified == 0 {
		return "", fmt.Errorf("user not verified")
	}

	if user.Password == "" || !utils.VerifyPassword(password, user.Password) {
		return "", fmt.Errorf("invalid password")
	}

	// JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}

	secretkey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", fmt.Errorf("failed generate token")
	}

	// Base URL fallback
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Chapter
	chapter := models.Chapter{}
	chaptercode := chapter.GetChapterCode(config.DB, int16(user.ClubID))

	// Session Redis
	session := utils.CustomerSession{
		UserID:      user.ID,
		Code:        user.QuinosID,
		Email:       user.Email,
		Name:        user.FirstName,
		Phone:       user.Phone1,
		Chapter:     user.ClubID,
		Image:       baseURL + "/uploads/" + user.Image,
		ChapterCode: chaptercode,
		Device:      "",
		Token:       tokenString,
		LoginAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := SetSession(session, 12*time.Hour); err != nil {
		return "", fmt.Errorf("redis error: %v", err)
	}

	// Log login
	logModel := models.Login{}
	device := ""
	if err := logModel.AddLog(config.DB, int(user.ID), &tokenString, &device); err != nil {
		return "", err
	}

	return tokenString, nil
}
