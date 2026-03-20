package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"go-rest/config"
	"go-rest/models"
	"go-rest/utils"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/exp/rand"
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
func (s *AuthService) Otp(username string) (int, int, error) { // return otp, userid, err

	otp := 200
	userid := 0
	customer := models.Customer{}
	validUser := customer.CheckUser(config.DB, username)
	validPhone := customer.CheckUserPhone(config.DB, username)

	if !validUser && !validPhone {
		return 400, userid, fmt.Errorf("User Not Found")
	}

	var res *models.Customer

	if validUser == true {
		res = customer.GetByUsername(config.DB, username)
	} else if validPhone == true {
		res = customer.GetByPhone(config.DB, username)
	}

	userid = int(res.ID)
	login := models.Login{}
	// get reqcount
	reslogin := login.CekRegCount(config.DB, userid)
	if reslogin.ReqCount < 3 {
		rand.Seed(uint64(time.Now().UnixNano()))
		// Menghasilkan angka acak antara 1000 dan 9999
		otp = rand.Intn(9000) + 1000
		err := login.SetOTP(config.DB, userid, otp)
		if err != nil {
			return 500, userid, fmt.Errorf("Failed to set otp..!")
		}
	} else {
		return 403, userid, fmt.Errorf("Maximum OTP Request..!")
	}

	return otp, userid, nil
}

func (s *AuthService) Forgot(username, password string, otp int) (bool, error) {

	customer := models.Customer{}
	validUser := customer.CheckUser(config.DB, username)
	validPhone := customer.CheckUserPhone(config.DB, username)

	if !validUser && !validPhone {
		return false, fmt.Errorf("User Not Found")
	}

	var res *models.Customer

	if validUser == true {
		res = customer.GetByUsername(config.DB, username)
	} else if validPhone == true {
		res = customer.GetByPhone(config.DB, username)
	}

	// cek apakah verified || status == 1
	if res.Status == 0 {
		return false, fmt.Errorf("User Not Active")
	}
	if res.Verified == 0 {
		return false, fmt.Errorf("User Not Verified")
	}
	if utils.VerifyPassword(password, res.Password) == true {
		return false, fmt.Errorf("Can't use previous password...!")
	}

	login := models.Login{}
	reslogin := login.CekRegCount(config.DB, int(res.ID))

	// fmt.Println("Res Login:", reslogin.Log)

	if reslogin.Log == nil || *reslogin.Log != strconv.Itoa(otp) {
		return false, fmt.Errorf("Invalid OTP.!")
	}

	// eksekusi set password
	err := models.UpdatePassword(config.DB, int(res.ID), password)
	if err != nil {
		return false, fmt.Errorf("Failed to update password..!")
	}

	// clear log
	device := ""
	err = login.EditLog(config.DB, int(res.ID), device, &device)
	if err != nil {
		return false, fmt.Errorf("Failed to update log..!")
	}

	return true, nil
}
