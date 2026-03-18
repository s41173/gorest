package controllers

import (
	"net/http"

	"go-rest/config"
	"go-rest/models"

	"github.com/gin-gonic/gin"
)

func Index_chapter(c *gin.Context) {

	var chapter []models.Chapter

	config.DB.Find(&chapter)
	c.JSON(http.StatusOK, gin.H{"chapter": chapter})
}
