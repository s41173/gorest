package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-rest/config"
	"go-rest/models"
	"go-rest/utils"

	// "go-rest/utils"

	"github.com/gin-gonic/gin"
)

// var ctx = context.Background()

func Register(c *gin.Context) {
	type RegistrationForm struct {
		CChapter  int    `form:"cchapter" binding:"required"`
		TName     string `form:"tname" binding:"required"`
		TAddress  string `form:"taddress" binding:"required"`
		TZip      string `form:"tzip"`
		TPhone1   string `form:"tphone1" binding:"required"`
		TEmail    string `form:"temail" binding:"required,email"`
		CCity     string `form:"ccity" binding:"required"`
		TPassword string `form:"tpassword" binding:"required,min=6"`
		TDOB      string `form:"tdob" binding:"required"`
		TNik      string `form:"tnik" binding:"required"`
		TCarType  string `form:"tcartype" binding:"required"`
		TPoliceNo string `form:"tpoliceno" binding:"required"`
	}

	// Bind data dari form
	var form RegistrationForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Inisiasi model
	customer := models.Customer{}
	chapter := models.Chapter{}
	city := models.City{}

	// Validasi city
	if !city.CheckCity(models.DB, form.CCity) {
		c.JSON(400, gin.H{"error": "City Not Found"})
		return
	}

	// Validasi chapter
	if !chapter.CheckChapter(models.DB, form.CChapter) {
		c.JSON(400, gin.H{"error": "Chapter Not Found"})
		return
	}

	// Validasi nomor telepon
	if customer.CheckUserPhone(models.DB, form.TPhone1) {
		c.JSON(400, gin.H{"error": "Phone registered"})
		return
	}

	// Validasi email
	if customer.CheckUser(models.DB, form.TEmail) {
		c.JSON(400, gin.H{"error": "Email registered"})
		return
	}

	// Hash password sebelum menyimpan
	hashedPassword, err := utils.HashPassword(form.TPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}
	joined := time.Now()
	// Membuat objek customer baru
	newCustomer := models.Customer{
		ClubID:    form.CChapter,
		FirstName: form.TName,
		Address:   form.TAddress,
		Zip:       form.TZip,
		Phone1:    form.TPhone1,
		Email:     form.TEmail,
		City:      form.CCity,
		Password:  hashedPassword,
		Dob:       &form.TDOB,
		NIK:       form.TNik,
		CarType:   form.TCarType,
		PoliceNo:  form.TPoliceNo,
		Type:      "member",
		Created:   &joined,
	}

	// Menyimpan data ke database
	if err := models.DB.Create(&newCustomer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	latesID, err := customer.GetLatestID(models.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data terbaru"})
		return
	}
	log := models.Login{}
	logdata := ""
	device := ""
	err = log.AddLog(models.DB, int(latesID), &logdata, &device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{"status": "Registrasi Succesed"})

}

func ChangePassword(c *gin.Context) {

	var requestData struct {
		Username    string `json:"username" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	// Bind JSON data ke dalam requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	customer := models.Customer{}
	var res *models.Customer
	validUser := customer.CheckUser(models.DB, requestData.Username)
	validPhone := customer.CheckUserPhone(models.DB, requestData.Username)

	if !validUser && !validPhone {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found..!"})
		return
	}

	if validUser {
		res = customer.GetByUsername(models.DB, requestData.Username)
	} else if validPhone {
		res = customer.GetByPhone(models.DB, requestData.Username)
	}

	if utils.VerifyPassword(requestData.NewPassword, res.Password) == true {
		c.JSON(403, gin.H{"error": "Can't use previous password...!"})
		return
	}

	err := models.UpdatePassword(models.DB, int(res.ID), requestData.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func Update_user(c *gin.Context) {

	var requestData struct {
		Name      string `json:"name" binding:"required"`
		Address   string `json:"address" binding:"required"`
		Zip       string `json:"zip" binding:"required"`
		City      string `json:"city" binding:"required"`
		Cartype   string `json:"cartype" binding:"required"`
		Vehicleno string `json:"vehicleno" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// validasi city
	var city models.City
	if !city.CheckCity(models.DB, requestData.City) {
		c.JSON(400, gin.H{"error": "City Not Found"})
		return
	}

	// ambil token
	tokenStr := utils.Token(c)
	claims, err := utils.DecodeToken(tokenStr)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// ambil user
	var customer models.Customer
	if err := models.DB.First(&customer, claims.UserID).Error; err != nil {
		c.JSON(401, gin.H{"error": "User Id Not Found"})
		return
	}

	// update data
	now := time.Now()

	customer.FirstName = requestData.Name
	customer.Address = requestData.Address
	customer.Zip = requestData.Zip
	customer.City = requestData.City
	customer.CarType = requestData.Cartype
	customer.PoliceNo = requestData.Vehicleno
	customer.Updated = &now

	if err := models.DB.Save(&customer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	// hapus session lama redis
	key := fmt.Sprintf("session:%d", customer.ID)
	config.RDB.Del(config.Ctx, key)

	// ambil chapter code
	var chapter models.Chapter
	chaptercode := chapter.GetChapterCode(models.DB, int16(customer.ClubID))

	// buat session baru
	sessionData := utils.CustomerSession{
		UserID:      customer.ID,
		Code:        customer.QuinosID,
		Email:       customer.Email,
		Name:        customer.FirstName,
		Phone:       customer.Phone1,
		Chapter:     customer.ClubID,
		ChapterCode: chaptercode,
		Token:       tokenStr,
		Device:      "",
		LoginAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(500, gin.H{"error": "Session encode error"})
		return
	}

	// simpan session redis
	if err := config.RDB.Set(
		config.Ctx,
		key,
		jsonData,
		12*time.Hour,
	).Err(); err != nil {

		c.JSON(500, gin.H{"error": "redis error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Update Success"})
}

func Update_image(c *gin.Context) {

	// ambil token
	tokenStr := utils.Token(c)
	claims, err := utils.DecodeToken(tokenStr)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// ambil user
	var customer models.Customer
	if err := models.DB.First(&customer, claims.UserID).Error; err != nil {
		c.JSON(401, gin.H{"error": "User Id Not Found"})
		return
	}

	// update data

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File Not Found",
		})
		return
	}

	// Max 1MB
	if file.Size > 1*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size max 1MB",
		})
		return
	}

	// Buka file untuk cek mime type
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read file",
		})
		return
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed read header file",
		})
		return
	}

	fileType := http.DetectContentType(buffer)

	// validasi mime type image
	if !strings.HasPrefix(fileType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image type required",
		})
		return
	}

	// ambil extension

	ext := filepath.Ext(file.Filename)
	fileName := utils.SplitSpace(strconv.Itoa(int(customer.ID))+"_"+customer.FirstName) + ext

	// lokasi simpan
	savePath := "./uploads/" + fileName

	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}

	now := time.Now()
	customer.Image = fileName
	customer.Updated = &now

	if err := models.DB.Save(&customer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	// hapus session lama redis
	key := fmt.Sprintf("session:%d", customer.ID)
	config.RDB.Del(config.Ctx, key)

	// ambil chapter code
	var chapter models.Chapter
	chaptercode := chapter.GetChapterCode(models.DB, int16(customer.ClubID))

	// buat session baru
	sessionData := utils.CustomerSession{
		UserID:      customer.ID,
		Code:        customer.QuinosID,
		Email:       customer.Email,
		Name:        customer.FirstName,
		Phone:       customer.Phone1,
		Chapter:     customer.ClubID,
		ChapterCode: chaptercode,
		Image:       utils.BaseURL(c) + "/uploads/" + fileName,
		Token:       tokenStr,
		Device:      "",
		LoginAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		c.JSON(500, gin.H{"error": "Session encode error"})
		return
	}

	// simpan session redis
	if err := config.RDB.Set(
		config.Ctx,
		key,
		jsonData,
		12*time.Hour,
	).Err(); err != nil {

		c.JSON(500, gin.H{"error": "redis error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Upload success",
		"path":    savePath,
		"type":    fileType,
		"ext":     ext,
	})
}

func Index_city(c *gin.Context) {

	city := models.City{}
	res, err := city.GetAllCities(models.DB)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"results": res})
	}
}
