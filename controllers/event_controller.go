package controllers

import (
	"go-rest/config"
	"go-rest/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func xxxIndex_event(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	if err := config.DB.First(&event, id).Error; err != nil {

		// fmt.Println("Error wois saat mencari produk:", err)

		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"Events": event})
}

func Index_event(c *gin.Context) {
	// Mendeklarasikan request struct untuk menampung data JSON dari request
	var requestData struct {
		Status  string `json:"status"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		Chapter string `json:"chapter"`
	}

	// Bind JSON data ke dalam requestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mengambil data property untuk mendapatkan image URL
	var property models.Property
	data, err := property.GetProperty(config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan data property"})
		return
	}
	imageUrl := data["image_url"].(string)

	// Mendeklarasikan slice untuk menyimpan events
	var events []models.Event
	query := config.DB

	// Filter berdasarkan chapter dan status jika parameter tersedia
	if requestData.Chapter != "" {
		query = query.Where("clubid = ?", requestData.Chapter)
	}
	if requestData.Status != "" {
		query = query.Where("done = ?", requestData.Status)
	}

	// Mengatur limit dan offset
	if requestData.Limit > 0 {
		query = query.Limit(requestData.Limit)
	} else {
		query = query.Limit(10) // Default limit jika tidak ada yang ditentukan
	}
	if requestData.Offset > 0 {
		query = query.Offset(requestData.Offset)
	}

	// Mendapatkan data dari database
	if err := query.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan data", "details": err.Error()})
		return
	}

	// Menyusun hasil looping dari events
	var result []gin.H
	for _, event := range events {

		// chapter
		var chapter models.Chapter
		chaptercode := chapter.GetChapterCode(config.DB, int16(event.Clubid))

		combinedImageUrl := imageUrl
		if event.Image != nil {
			combinedImageUrl += *event.Image
		}

		donedesc := "AVAILABLE"
		if event.Done == 1 {
			donedesc = "DONE"
		}

		typdesc := "PRIVATE"
		if event.Type == 1 {
			typdesc = "COMBINATED"
		}

		// Menambahkan data event ke dalam hasil
		eventData := gin.H{
			"id":                  event.ID,
			"chapter":             chaptercode,
			"chapter_id":          event.Clubid,
			"code":                event.Code,
			"name":                event.Name,
			"dates":               event.Dates,
			"desc":                event.Desc,
			"image":               combinedImageUrl,
			"fee":                 event.Fee,
			"type":                event.Type,
			"type_desc":           typdesc,
			"minimum_participant": event.MinimumParticipant,
			"done":                event.Done,
			"done_desc":           donedesc,
			"point":               event.Point,
			"created":             event.Created,
			"updated":             event.Updated,
			"deleted":             event.Deleted,
		}
		// fmt.Printf("Event ID: %d, ClubID: %d, Name: %s\n", event.ID, event.ClubID, event.Name)
		result = append(result, eventData)
	}

	// Mengirimkan respons JSON
	c.JSON(http.StatusOK, gin.H{"results": result})
}

// fungsi get by id
func Get_event(c *gin.Context) {
	var event models.Event
	id := c.Param("id")

	if err := config.DB.First(&event, id).Error; err != nil {

		// fmt.Println("Error wois saat mencari produk:", err)

		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"product": event})
}
