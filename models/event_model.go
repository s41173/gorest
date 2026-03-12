package models

import (
	"gorm.io/gorm"
)

type Event struct {
	ID                 int16   `gorm:"primaryKey;autoIncrement" json:"id"` // ID, smallint(3), auto_increment
	Clubid             int32   `gorm:"type:int(5);not null" json:"clubid"`
	Code               *string `gorm:"type:varchar(100);default:null" json:"code"`           // Code, varchar(100), nullable
	Name               string  `gorm:"type:text;not null" json:"name"`                       // Name, text, NOT NULL
	Dates              string  `gorm:"type:datetime;not null" json:"dates"`                  // Dates, datetime, NOT NULL
	Desc               *string `gorm:"type:text;default:null" json:"desc"`                   // Desc, text, nullable
	Image              *string `gorm:"type:varchar(100);default:null" json:"image"`          // Image, varchar(100), nullable
	Fee                float64 `gorm:"type:decimal(9,2);not null" json:"fee"`                // Fee, decimal(9,2), NOT NULL
	Type               int8    `gorm:"type:smallint(1);not null;default:0" json:"type"`      // Type, smallint(1), NOT NULL, default 0
	MinimumParticipant int16   `gorm:"type:smallint(4);not null" json:"minimum_participant"` // MinimumParticipant, smallint(4), NOT NULL
	Done               int8    `gorm:"type:smallint(1);not null;default:0" json:"done"`      // Done, smallint(1), NOT NULL, default 0
	Point              int16   `gorm:"type:smallint(3);not null;default:0" json:"point"`     // Point, smallint(3), NOT NULL, default 0
	Created            *string `gorm:"type:datetime;default:null" json:"created"`            // Created, datetime, nullable
	Updated            *string `gorm:"type:datetime;default:null" json:"updated"`            // Updated, datetime, nullable
	Deleted            *string `gorm:"type:datetime;default:null" json:"deleted"`            // Deleted, datetime, nullable
}

// TableName mengatur nama tabel untuk model Event
func (Event) TableName() string {
	return "event"
}

// GetLast retrieves a filtered list of events or the count of such events based on the parameters
func (e *Event) GetLast(db *gorm.DB, chapter *int16, status *int8, limit int, offset *int, count bool) (interface{}, error) {
	query := db.Model(&Event{}).Where("done = ?", 0).Order("done ASC").Order("dates ASC")

	// Filter by chapter if provided
	if chapter != nil {
		query = query.Where("clubid = ?", *chapter)
	}

	// Filter by status if provided
	if status != nil {
		query = query.Where("done = ?", *status)
	}

	// Return the count if count flag is true
	if count {
		var total int64
		if err := query.Count(&total).Error; err != nil {
			return nil, err
		}
		return total, nil
	}

	// Apply limit and offset if count flag is false
	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Retrieve events based on the query
	var events []Event
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
