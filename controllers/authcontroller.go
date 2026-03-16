package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-rest/config"
	"go-rest/models"
	"go-rest/utils"

	// "go-rest/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
)

// var ctx = context.Background()

func TestRedis(c *gin.Context) {

	// set value
	err := config.RDB.Set(config.Ctx, "test_key", "Hello Redis", 10*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get value
	val, err := config.RDB.Get(config.Ctx, "test_key").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"redis_value": val,
	})
}

func Verify(c *gin.Context) {

	var requestData struct {
		Phone string `json:"phone" binding:"required"`
		Otp   int    `json:"otp" binding:"required"`
	}
	// Bind JSON data ke dalam requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	customer := models.Customer{}
	validPhone := customer.CheckUserPhone(models.DB, requestData.Phone)
	if validPhone == false {
		c.JSON(404, gin.H{"error": "User Not Found"})
		return
	}

	res := customer.GetByPhone(models.DB, requestData.Phone)
	if res.Verified == 1 {
		c.JSON(403, gin.H{"error": "User Has Been Verified"})
		return
	}

	login := models.Login{}
	if login.CekOTP(models.DB, int(res.ID), requestData.Otp) == true {
		c.JSON(403, gin.H{"error": "OTP tidak sesuai"})
		return
	}

	result := customer.UpdateVerified(models.DB, int(res.ID))
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": "Verified Success"})
}

func ForgotPassword(c *gin.Context) {

	var requestData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Otp      int    `json:"otp" binding:"required"`
	}
	// Bind JSON data ke dalam requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// fmt.Println("Hasil request : ", requestData)

	customer := models.Customer{}
	validUser := customer.CheckUser(models.DB, requestData.Username)
	validPhone := customer.CheckUserPhone(models.DB, requestData.Username)

	if !validUser && !validPhone {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found..!"})
		return
	}

	var res *models.Customer

	if validUser == true {
		res = customer.GetByUsername(models.DB, requestData.Username)
	} else if validPhone == true {
		res = customer.GetByPhone(models.DB, requestData.Username)
	}

	// cek apakah verified || status == 1
	if res.Status == 0 {
		c.JSON(453, gin.H{"error": "User Not Active..!"})
		return
	}
	if res.Verified == 0 {
		c.JSON(452, gin.H{"error": "User Not Verified..!"})
		return
	}
	if utils.VerifyPassword(requestData.Password, res.Password) == true {
		c.JSON(403, gin.H{"error": "Can't use previous password...!"})
		return
	}
	// fmt.Println("Eksekusi : ")
	login := models.Login{}
	reslogin := login.CekRegCount(models.DB, int(res.ID))

	// fmt.Println("Res Login:", reslogin.Log)

	if reslogin.Log == nil || *reslogin.Log != strconv.Itoa(requestData.Otp) {
		c.JSON(403, gin.H{"error": "Invalid OTP"})
		return
	}

	// eksekusi set password
	err := models.UpdatePassword(models.DB, int(res.ID), requestData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func Logout(c *gin.Context) {

	tokenStr := utils.Token(c)
	claims, err := utils.DecodeToken(tokenStr)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
	}
	login := models.Login{}
	login.LogoutUser(models.DB, claims.UserID)

	// redis delete
	key := fmt.Sprintf("session:%d", claims.UserID)

	err = config.RDB.Del(config.Ctx, key).Err()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed logout",
		})
		return
	}

	// redis end

	c.JSON(200, gin.H{
		"result": "logout success",
	})

}

func Decode(c *gin.Context) {

	tokenStr := utils.Token(c)
	claims, err := utils.DecodeToken(tokenStr)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	key := fmt.Sprintf("session:%d", claims.UserID)

	val, err := config.RDB.Get(config.Ctx, key).Result()
	if err == redis.Nil {
		c.JSON(401, gin.H{"error": "session expired"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var session map[string]interface{}

	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed decode redis json"})
		return
	}

	c.JSON(200, gin.H{
		"result": session,
	})
}

func Login(c *gin.Context) {
	var requestData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// Bind JSON data ke dalam requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// fmt.Println("Hasil request : ", requestData)

	customer := models.Customer{}
	chapter := models.Chapter{}
	validUser := customer.CheckUser(models.DB, requestData.Username)
	validPhone := customer.CheckUserPhone(models.DB, requestData.Username)

	if !validUser && !validPhone {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found..!"})
		return
	} else {
		var res *models.Customer

		if validUser == true {
			res = customer.GetByUsername(models.DB, requestData.Username)
		} else if validPhone == true {
			res = customer.GetByPhone(models.DB, requestData.Username)
		}

		// cek apakah verified || status == 1
		if res.Status == 0 {
			c.JSON(453, gin.H{"error": "User Not Active..!"})
			return
		} else if res.Verified == 0 {
			c.JSON(452, gin.H{"error": "User Not Verified..!"})
		} else {
			if res != nil && utils.VerifyPassword(requestData.Password, res.Password) { // Pastikan Anda mengganti *res.Password sesuai dengan field yang benar
				// User valid dan password cocok

				chaptercode := chapter.GetChapterCode(models.DB, int16(res.ClubID))
				// fmt.Println("Chapter Code:", chaptercode)

				// add jwt
				// Membuat token JWT
				claims := struct {
					UserID int64 `json:"userid"`
					jwt.RegisteredClaims
				}{
					UserID: res.ID,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // Contoh: 2 jam dari sekarang
					},
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString([]byte("merciku")) // Ganti "vinkoo" dengan secret key yang lebih aman
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
					return
				}

				// tambah log
				log := models.Login{}
				device := ""
				err = log.AddLog(models.DB, int(res.ID), &tokenString, &device)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// session untuk redis
				sessionData := utils.CustomerSession{
					UserID:      res.ID,
					Code:        res.QuinosID,
					Email:       res.Email,
					Name:        res.FirstName,
					Phone:       res.Phone1,
					Chapter:     res.ClubID,
					Image:       utils.BaseURL(c) + "/uploads/" + res.Image,
					ChapterCode: chaptercode,
					Token:       tokenString,
					Device:      device,
					LoginAt:     time.Now().Format("2006-01-02 15:04:05"),
				}

				jsonData, err := json.Marshal(sessionData)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Session encode error"})
					return
				}

				// simpan redis
				key := fmt.Sprintf("session:%d", res.ID)

				err = config.RDB.Set(
					config.Ctx,
					key,
					jsonData,
					12*time.Hour).Err()

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "redis error",
					})
					return
				}

				// end redis

				// Mengirimkan token sebagai respon
				c.JSON(200, gin.H{
					"token": tokenString,
					"type":  "Bearer",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
			}
		}
	}

}

func RegOTP(c *gin.Context) {
	var requestData struct {
		Username string `json:"username"`
	}
	// Bind JSON data ke dalam requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// fmt.Println("Hasil request : ", requestData)

	customer := models.Customer{}
	validUser := customer.CheckUser(models.DB, requestData.Username)
	validPhone := customer.CheckUserPhone(models.DB, requestData.Username)

	if !validUser && !validPhone {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found..!"})
		return
	} else {
		var res *models.Customer

		if validUser == true {
			res = customer.GetByUsername(models.DB, requestData.Username)
		} else if validPhone == true {
			res = customer.GetByPhone(models.DB, requestData.Username)
		}

		// cek apakah verified || status == 1
		// if res.Status == 0 {
		// 	c.JSON(453, gin.H{"error": "User Not Active..!"})
		// 	return
		// }
		// else if res.Verified == 0 {
		// 	c.JSON(452, gin.H{"error": "User Not Verified..!"})
		// }

		login := models.Login{}
		// get reqcount
		reslogin := login.CekRegCount(models.DB, int(res.ID))
		// fmt.Println("Log data : ", reslogin.ReqCount)
		if reslogin.ReqCount < 3 {
			rand.Seed(uint64(time.Now().UnixNano()))
			// Menghasilkan angka acak antara 1000 dan 9999
			logid := rand.Intn(9000) + 1000
			err := login.SetOTP(models.DB, int(res.ID), logid)
			if err == nil {
				// kirim notif
				userIDStr := strconv.Itoa(int(res.ID))
				currentTime := time.Now().Format("2006-01-02 15:04:05")                                  // Format to desired layout
				logMessage := "OTP Reg : " + currentTime + " Kode Pin OTP Anda : " + strconv.Itoa(logid) // Convert logid to string
				// fmt.Println(logMessage)
				notif := utils.SendNotif(c, "7", userIDStr, logMessage, logMessage, "login")
				if notif == false {
					c.JSON(403, gin.H{"error": "Failure to send otp..!"})
				} else {
					c.JSON(200, gin.H{"error": "OTP has been sent..!"})
				}
			}
		} else {
			c.JSON(403, gin.H{"error": "Maximum OTP Request..!"})
		}

	}
}
