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
	correctPassword := "j4ykiran1"
	wrongPassword := "salah123"
	nonexistentUser := "tidak_ada@example.com"

	// -------------------------------
	// 1️⃣ Username salah → harus "user not found"
	// fmt.Println("\n[TEST] Login with non-existent user")
	_, err := service.Login(nonexistentUser, correctPassword)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("Expected 'user not found', got: %v", err)
	}

	// -------------------------------
	// 2️⃣ Password salah → harus "invalid password"
	fmt.Println("\n[TEST] Login with wrong password")
	_, err = service.Login(existingUser, wrongPassword)
	if err == nil || err.Error() != "invalid password" {
		t.Errorf("Expected 'invalid password', got: %v", err)
	}

	// -------------------------------
	// 3️⃣ Login sukses → dapat JWT
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
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
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
