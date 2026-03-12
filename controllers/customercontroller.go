package controllers

import (
	"net/http"
	"time"

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
