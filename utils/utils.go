// utils/password.go

package utils

import (
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetLocalTime mengembalikan waktu saat ini dalam zona waktu lokal
func GetLocalTime() time.Time {
	// Ganti dengan nama zona waktu yang sesuai
	loc, err := time.LoadLocation("Asia/Jakarta") // Misalnya untuk WIB
	if err != nil {
		panic(err) // Atau tangani kesalahan sesuai kebutuhan Anda
	}
	return time.Now().In(loc) // Kembalikan waktu saat ini dalam zona waktu lokal
}

func SplitSpace(s string) string {

	// lowercase
	s = strings.ToLower(s)

	// replace non alphanumeric dengan "-"
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	// trim "-"
	s = strings.Trim(s, "-")

	return s
}

func BaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}
