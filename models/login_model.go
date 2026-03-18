package models

import (
	"time"

	"gorm.io/gorm"
)

// Club struct
type Login struct {
	ID         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int        `gorm:"column:userid" json:"userid"`
	Log        *string    `gorm:"type:text;default:null" json:"log"`
	Device     *string    `gorm:"type:text;default:null" json:"device"`
	Joined     *time.Time `gorm:"type:datetime;default:null" json:"joined"`
	ReqCount   int16      `gorm:"type:smallint;not null" json:"req_count"`
	ReqCreated *time.Time `gorm:"type:datetime;default:null" json:"req_created"`
}

// TableName sets the name of the database table
func (Login) TableName() string {
	return "customer_login_status"
}

func (login *Login) AddLog(db *gorm.DB, userID int, logData *string, device *string) error {
	if login.CekUser(db, userID) {
		// Jika user belum ada, lakukan insert
		joinedTime := time.Now().UTC()
		newLog := Login{
			UserID: userID,
			Log:    logData,
			Device: device,
			Joined: &joinedTime,
		}
		if err := db.Create(&newLog).Error; err != nil {
			return err
		}
	} else {
		// Jika user sudah ada, lakukan update
		if err := login.EditLog(db, userID, *logData, device); err != nil {
			return err
		}
	}
	return nil
}

func (login *Login) CekUser(db *gorm.DB, userID int) bool {
	var count int64
	db.Model(&Login{}).Where("userid = ?", userID).Count(&count)
	return count == 0
}

func (login *Login) CekUserToken(db *gorm.DB, userID int, log string) bool {
	var count int64
	db.Model(&Login{}).
		Where("userid = ?", userID).
		Where("log = ?", log).
		Count(&count)
	return count > 0
}

// function untuk cek req count

func (login *Login) CekRegCount(db *gorm.DB, userID int) *Login {
	var result Login
	// today := time.Now().Format("2006-01-02")
	err := db.Where("userid = ?", userID).Limit(1).
		// Where("DATE(req_created) = ?", today).
		First(&result).Error

	// Jika tidak ditemukan, kembalikan nil
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return &result
}

func (login *Login) EditLog(db *gorm.DB, userID int, logData string, device *string) error {
	updateData := map[string]interface{}{
		"log":    logData,
		"device": device,
		"joined": time.Now().Format("2006-01-02 15:04:05"),
	}

	// fmt.Println("Waktu : ", updateData)

	// Melakukan update berdasarkan userID
	if err := db.Model(&Login{}).Where("userid = ?", userID).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

func (login *Login) LogoutUser(db *gorm.DB, userID int64) error {
	updateData := map[string]interface{}{
		"log":       nil,
		"device":    nil,
		"joined":    gorm.Expr("NULL"),
		"req_count": 0,
	}

	// fmt.Println("Waktu : ", updateData)

	// Melakukan update berdasarkan userID
	if err := db.Model(&Login{}).Where("userid = ?", userID).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

// SetOTP mengatur OTP dan memperbarui `req_count` di database
func (login *Login) SetOTP(db *gorm.DB, userID int, otp int) error {
	// Menyiapkan data untuk diupdate
	updateData := map[string]interface{}{
		"log":         otp,
		"req_count":   gorm.Expr("IFNULL(req_count, 0) + 1"), // Increment req_count
		"req_created": time.Now(),
	}

	// Update data pada tabel berdasarkan `userid`
	err := db.Model(&Login{}).
		Where("userid = ?", userID).
		Updates(updateData).Error

	return err
}

func (login *Login) CekOTP(db *gorm.DB, userID int, Otp int) bool {
	var count int64
	db.Model(&Login{}).
		Where("userid = ?", userID).
		Where("log = ?", Otp).
		Count(&count)
	return count == 0
}
