package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Customer struct
type Customer struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	QuinosID        string     `gorm:"type:varchar(100);default:null" json:"quinos_id"`
	ClubID          int        `gorm:"column:clubid;type:int(5);default:null" json:"clubid"`
	FirstName       string     `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName        string     `gorm:"type:varchar(100);default:null" json:"last_name"`
	Type            string     `gorm:"type:enum('member','veteran');not null" json:"type"`
	Address         string     `gorm:"type:text;not null" json:"address"`
	ShippingAddress string     `gorm:"type:text;not null" json:"shipping_address"`
	Phone1          string     `gorm:"type:varchar(20);not null" json:"phone1"`
	Phone2          string     `gorm:"type:varchar(25);default:null" json:"phone2"`
	Fax             string     `gorm:"type:varchar(25);not null" json:"fax"`
	Email           string     `gorm:"type:varchar(100);not null" json:"email"`
	Password        string     `gorm:"type:text;not null" json:"password"`
	Website         string     `gorm:"type:varchar(50);not null" json:"website"`
	State           string     `gorm:"type:varchar(100);not null" json:"state"`
	City            string     `gorm:"type:varchar(100);not null" json:"city"`
	Region          string     `gorm:"type:varchar(100);not null" json:"region"`
	Zip             string     `gorm:"type:varchar(10);not null" json:"zip"`
	Notes           string     `gorm:"type:text;not null" json:"notes"`
	Image           string     `gorm:"type:varchar(100);default:null" json:"image"`
	NPWP            string     `gorm:"type:varchar(100);default:null" json:"npwp"`
	Profession      string     `gorm:"type:varchar(100);default:null" json:"profession"`
	Organization    string     `gorm:"type:varchar(100);default:null" json:"organization"`
	MemberNo        string     `gorm:"type:varchar(100);default:null" json:"member_no"`
	Instagram       string     `gorm:"type:varchar(100);default:null" json:"instagram"`
	Joined          *string    `gorm:"type:datetime;default:null" json:"joined"`
	Premium         int8       `gorm:"type:smallint(1);not null" json:"premium"`
	Status          int8       `gorm:"type:smallint(1);not null" json:"status"`
	VoucherClaimed  int8       `gorm:"type:smallint(1);not null" json:"voucher_claimed"`
	MType           int8       `gorm:"column:mtype;type:smallint(3);not null" json:"mtype"`
	Dob             *string    `gorm:"type:date;default:null" json:"dob"` // Pointer to allow null
	NIK             string     `gorm:"type:text;default:null" json:"nik"`
	CarType         string     `gorm:"type:text;default:null" json:"car_type"`
	ChasisNo        string     `gorm:"type:text;default:null" json:"chasis_no"`
	EngineNo        string     `gorm:"type:text;default:null" json:"engine_no"`
	PoliceNo        string     `gorm:"type:varchar(50);default:null" json:"police_no"`
	CarImage        string     `gorm:"type:text;default:null" json:"car_image"`
	Expired         *string    `gorm:"type:datetime;default:null" json:"expired"` // Pointer to allow null
	Verified        int8       `gorm:"type:smallint(1);not null;default:0" json:"verified"`
	BillingNum      int8       `gorm:"type:smallint(2);not null;default:1" json:"billing_num"`
	Created         *time.Time `gorm:"type:datetime;default:null" json:"created"` // Pointer to allow null
	Updated         *time.Time `gorm:"type:datetime;default:null" json:"updated"` // Pointer to allow null
	Deleted         *time.Time `gorm:"type:datetime;default:null" json:"deleted"` // Pointer to allow null
}

func (Customer) TableName() string {
	return "customer"
}

func (c *Customer) CheckUser(db *gorm.DB, username string) bool {
	var count int64
	db.Model(&Customer{}).
		Where("email = ?", username).
		// Where("status = ?", 1).
		Where("deleted IS NULL").
		Count(&count)

	return count > 0
}

func (c *Customer) CheckUserPhone(db *gorm.DB, phone string) bool {
	var count int64
	db.Model(&Customer{}).
		Where("phone1 = ?", phone).
		// Where("status = ?", 1).
		Where("deleted IS NULL").
		Count(&count)

	return count > 0
}

// GetByUsername mencari data customer berdasarkan email dengan kondisi status dan deleted
func (customer *Customer) GetByUsername(db *gorm.DB, username string) *Customer {
	var result Customer
	// Mencari data customer dengan kondisi tertentu
	err := db.Model(&Customer{}).
		Where("email = ?", username).
		// Where("status = ?", 1).
		Where("deleted IS NULL").
		First(&result).Error

	// Jika tidak ditemukan, kembalikan nil
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return &result
}

func (customer *Customer) GetByPhone(db *gorm.DB, phone string) *Customer {
	var result Customer
	// Mencari data customer dengan kondisi tertentu
	err := db.Model(&Customer{}).
		Where("phone1 = ?", phone).
		// Where("status = ?", 1).
		Where("deleted IS NULL").
		First(&result).Error

	// Jika tidak ditemukan, kembalikan nil
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return &result
}

// UpdatePassword untuk memperbarui password customer
func UpdatePassword(db *gorm.DB, userID int, newPassword string) error {
	// Mengenkripsi password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password di database
	err = db.Model(&Customer{}).Where("id = ?", userID).Update("password", hashedPassword).Error
	return err
}

func (l *Customer) GetLatestID(db *gorm.DB) (int, error) {
	var id int
	err := db.Model(&Customer{}).
		Select("id").
		Order("id DESC").
		Limit(1).
		Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (customer *Customer) CekCustomer(db *gorm.DB, userID int) bool {
	var count int64
	db.Model(&Customer{}).
		Where("id = ?", userID).
		Where("verified = 1").
		Or("Status = 1").
		Count(&count)
		// tambahkan OR
	return count == 0
}

func (customer *Customer) UpdateVerified(db *gorm.DB, userID int) error {

	// Update password di database
	err := db.Model(&Customer{}).Where("id = ?", userID).Update("verified", true).Error
	return err
}
