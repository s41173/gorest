package tests

import (
	"fmt"
	"os"
	"testing"

	"go-rest/config"
	"go-rest/services"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Load .env dari root
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// Connect DB
	config.ConnectDatabase()
	if config.DB == nil {
		panic("Database pointer is nil")
	}

	fmt.Println("Database connected successfully")

	// Connect Redis
	config.InitRedis()

	// Jalankan test
	code := m.Run()

	os.Exit(code) // 🔥 WAJIB
}

func TestLoginScenarios(t *testing.T) {
	service := services.NewAuthService()

	// 🔹 Ganti dengan user yang benar-benar ada di DB
	existingUser := "082277014410"
	correctPassword := "newpassword"

	// -------------------------------
	// Login sukses → dapat JWT
	fmt.Println("\n[TEST] Login with correct credentials")
	token, err := service.Login(existingUser, correctPassword)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Fatal("Expected JWT token, got empty string")
	}

	fmt.Println("Generated token:", token)

	// Validasi JWT
	secretkey := os.Getenv("JWT_SECRET")
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatal("Invalid JWT token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed convert claims")
	}

	fmt.Println("JWT claims:", claims)
}

func TestOtpScenarios(t *testing.T) {
	service := services.NewAuthService()

	username := "082277014410"

	otp, userid, err := service.Otp(username)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if otp == 0 {
		t.Errorf("expected otp, got 0")
	}

	if userid == 0 {
		t.Errorf("expected valid user id")
	}
}

func TestForgotPassScenarios(t *testing.T) {
	service := services.NewAuthService()

	username := "082277014410"
	password := "j4ykiran"

	// 🔹 generate OTP dulu
	otp, _, err := service.Otp(username)
	if err != nil {
		t.Fatalf("failed generate otp: %v", err)
	}

	// 🔹 pakai OTP yang valid
	res, err := service.Forgot(username, password, otp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res {
		t.Errorf("expected success = true")
	}
}
