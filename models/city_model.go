package models

import "gorm.io/gorm"

// Club struct
type City struct {
	ID       int16   `gorm:"primaryKey;autoIncrement" json:"id"`
	Id_prov  *string `gorm:"type:varchar(100);default:null" json:"id_prov"`
	Province string  `gorm:"type:text;not null" json:"province"`
	Type     *string `gorm:"type:text;default:null" json:"type"`
	Nama     *string `gorm:"type:text;default:null" json:"nama"`
	Zip      *string `gorm:"type:varchar(100);default:null" json:"zip"`
}

func (City) GetAllCities(db *gorm.DB) ([]City, error) {
	var cities []City

	err := db.Order("nama ASC").Find(&cities).Error
	if err != nil {
		return nil, err
	}

	return cities, nil
}

// versi bahasa manusia
func (c City) CheckCity(db *gorm.DB, nama string) bool {
	var count int64
	db.Model(&City{}).
		Where("nama = ?", nama).
		Count(&count)

	if count > 0 {
		return true
	} else {
		return false
	}
}

// TableName sets the name of the database table (Chapter == nama struct)
func (City) TableName() string {
	return "kabupaten_rj"
}
