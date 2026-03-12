package models

import (
	"time"

	"gorm.io/gorm"
)

// Club struct
type Chapter struct {
	ID             int16      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code           *string    `gorm:"type:varchar(100);default:null" json:"code"`
	Name           string     `gorm:"type:text;not null" json:"name"`
	Desc           *string    `gorm:"type:text;default:null" json:"desc"`
	Instagram      *string    `gorm:"type:text;default:null" json:"instagram"`
	Logo           *string    `gorm:"type:varchar(100);default:null" json:"logo"`
	Bank1          *string    `gorm:"type:text;default:null" json:"bank1"`
	Bank2          *string    `gorm:"type:text;default:null" json:"bank2"`
	Bank3          *string    `gorm:"type:text;default:null" json:"bank3"`
	City           *string    `gorm:"type:text;default:null" json:"city"`
	Email          string     `gorm:"type:text;not null" json:"email"`
	Phone          string     `gorm:"type:varchar(50);not null" json:"phone"`
	Address        string     `gorm:"type:text;not null" json:"address"`
	RegistFee      float64    `gorm:"type:decimal(9,2);not null" json:"regist_fee"`
	MonthlyFee     float64    `gorm:"type:decimal(9,2);not null" json:"monthly_fee"`
	Discount       float64    `gorm:"type:decimal(7,3);not null" json:"discount"`
	Chief          *string    `gorm:"type:text;default:null" json:"chief"`
	Vice           *string    `gorm:"type:text;default:null" json:"vice"`
	Treasurer      *string    `gorm:"type:text;default:null" json:"treasurer"`
	BillingContact string     `gorm:"type:varchar(50);not null" json:"billing_contact"`
	Parent         int8       `gorm:"type:smallint(1);not null;default:0" json:"parent"`
	Created        *time.Time `gorm:"type:datetime;default:null" json:"created"` // Pointer to allow null
	Updated        *time.Time `gorm:"type:datetime;default:null" json:"updated"` // Pointer to allow null
	Deleted        *time.Time `gorm:"type:datetime;default:null" json:"deleted"` // Pointer to allow null
}

// versi bahasa manusia
func (chapter Chapter) GetChapterCode(db *gorm.DB, id int16) string {
	var result Chapter

	err := db.First(&result, id).Error
	if err != nil {
		return ""
	}

	if result.Code != nil {
		return *result.Code
	}

	return ""
}

// versi bahasa manusia
func (c Chapter) CheckChapter(db *gorm.DB, chapterid int) bool {
	var count int64
	db.Model(&Chapter{}).
		Where("id = ?", chapterid).
		Where("deleted IS NULL").
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

// TableName sets the name of the database table (Chapter == nama struct)
func (Chapter) TableName() string {
	return "club"
}
