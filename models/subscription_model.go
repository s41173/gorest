package models

import (
	"time"
)

// Club struct
type Subscription struct {
	ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Customer      int        `gorm:"column:customer" json:"customer"`
	P1            *time.Time `gorm:"column:p1" json:"p1"`
	P2            *time.Time `gorm:"column:p2" json:"p2"`
	P3            *time.Time `gorm:"column:p3" json:"p3"`
	P4            *time.Time `gorm:"column:p4" json:"p4"`
	P5            *time.Time `gorm:"column:p5" json:"p5"`
	P6            *time.Time `gorm:"column:p6" json:"p6"`
	P7            *time.Time `gorm:"column:p7" json:"p7"`
	P8            *time.Time `gorm:"column:p8" json:"p8"`
	P9            *time.Time `gorm:"column:p9" json:"p9"`
	P10           *time.Time `gorm:"column:p10" json:"p10"`
	P11           *time.Time `gorm:"column:p11" json:"p11"`
	P12           *time.Time `gorm:"column:p12" json:"p12"`
	RegMonth      string     `gorm:"column:reg_month;type:varchar(5);default:'1'" json:"reg_month"`
	FinancialYear string     `gorm:"column:financial_year;type:varchar(4)" json:"financial_year"`
	Created       *time.Time `gorm:"type:datetime;default:null" json:"created"` // Pointer to allow null
	Updated       *time.Time `gorm:"type:datetime;default:null" json:"updated"` // Pointer to allow null
	Deleted       *time.Time `gorm:"type:datetime;default:null" json:"deleted"` // Pointer to allow null
}

// TableName sets the name of the database table
func (Subscription) TableName() string {
	return "customer_subscription"
}

// func (subscription *Subscription) AddSubscription(db *gorm.DB, userID int, logData *string, device *string) error {
// 	if subscription.CekUser(db, userID) {
// 		// Jika user belum ada, lakukan insert
// 		joinedTime := time.Now()
// 		newLog := Login{
// 			UserID: userID,
// 			Log:    logData,
// 			Device: device,
// 			Joined: &joinedTime,
// 		}
// 		if err := db.Create(&newLog).Error; err != nil {
// 			return err
// 		}
// 	} else {
// 		// Jika user sudah ada, lakukan update
// 		if err := login.EditLog(db, userID, *logData, device); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
