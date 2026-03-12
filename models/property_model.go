package models

import (
	"strings"

	"gorm.io/gorm"
)

type Property struct {
	ID                  int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                string  `gorm:"type:varchar(100);not null" json:"name"`
	Address             string  `gorm:"type:text;not null" json:"address"`
	Coordinate          string  `gorm:"type:varchar(50);not null" json:"coordinate"`
	Phone1              string  `gorm:"type:varchar(100);not null" json:"phone1"`
	Phone2              string  `gorm:"type:varchar(100);not null" json:"phone2"`
	Fax                 string  `gorm:"type:varchar(100);not null" json:"fax"`
	Email               string  `gorm:"type:varchar(100);not null" json:"email"`
	BillingEmail        string  `gorm:"type:varchar(100);not null" json:"billing_email"`
	TechnicalEmail      *string `gorm:"type:varchar(100);default:null" json:"technical_email"`
	CcEmail             string  `gorm:"type:varchar(100);not null" json:"cc_email"`
	Zip                 int64   `gorm:"type:int;not null" json:"zip"`
	City                string  `gorm:"type:char(30);not null" json:"city"`
	AccountName         string  `gorm:"type:varchar(100);not null" json:"account_name"`
	AccountNo           string  `gorm:"type:varchar(100);not null" json:"account_no"`
	Bank                string  `gorm:"type:text;not null" json:"bank"`
	SiteName            string  `gorm:"type:varchar(100);not null" json:"site_name"`
	Logo                *string `gorm:"type:text;default:null" json:"logo"`
	MetaDescription     string  `gorm:"type:text;not null" json:"meta_description"`
	MetaKeyword         string  `gorm:"type:text;not null" json:"meta_keyword"`
	Manager             *string `gorm:"type:varchar(100);default:null" json:"manager"`
	Accounting          *string `gorm:"type:varchar(100);default:null" json:"accounting"`
	EmailLink           *string `gorm:"type:text;default:null" json:"email_link"`
	UrlUpload           *string `gorm:"type:text;default:null" json:"url_upload"`
	ImageUrl            *string `gorm:"type:text;default:null" json:"image_url"`
	ApiImageUrl         *string `gorm:"type:text;default:null" json:"api_image_url"`
	NotifUrl            *string `gorm:"type:text;default:null" json:"notif_url"`
	NotifToken          *string `gorm:"type:text;default:null" json:"notif_token"`
	PosUrl              *string `gorm:"type:text;default:null" json:"pos_url"`
	PosToken            *string `gorm:"type:text;default:null" json:"pos_token"`
	InvoiceUrl          *string `gorm:"type:text;default:null" json:"invoice_url"`
	PgUrl               *string `gorm:"type:text;default:null" json:"pg_url"`
	PgToken             *string `gorm:"type:text;default:null" json:"pg_token"`
	ShipUrl             *string `gorm:"type:text;default:null" json:"ship_url"`
	ShipToken           *string `gorm:"type:text;default:null" json:"ship_token"`
	CourierIntegration  int8    `gorm:"type:smallint;not null;default:0" json:"courier_integration"`
	ShippingIntegration int8    `gorm:"type:smallint;not null;default:0" json:"shipping_integration"`
	ShippingVendor      string  `gorm:"type:varchar(30);not null" json:"shipping_vendor"`
	OperationalHours    *string `gorm:"type:varchar(5);default:null" json:"operational_hours"`
	DistanceLimit       int16   `gorm:"type:smallint;not null" json:"distance_limit"`
	PointCalculate      float64 `gorm:"type:decimal(9,0);not null" json:"point_calculate"`
	Created             *string `gorm:"type:datetime;null" json:"created"`
	Updated             *string `gorm:"type:datetime;default:null" json:"updated"`
	Deleted             *string `gorm:"type:datetime;default:null" json:"deleted"`
}

// GetProperty retrieves property data from the database
func (p *Property) GetProperty(db *gorm.DB) (map[string]interface{}, error) {
	var property Property

	// Retrieve data from database
	if err := db.First(&property).Error; err != nil {
		return nil, err
	}

	// Handle operational_hours
	var start, end string
	if property.OperationalHours != nil {
		phours := strings.Split(*property.OperationalHours, "-")
		if len(phours) == 2 {
			start = phours[0]
			end = phours[1]
		}
	}

	// Handle URL upload
	urlUpload := "./"
	if property.UrlUpload != nil {
		urlUpload = *property.UrlUpload
	}

	// Handle image URL
	imageUrl := "base_url/images/"
	if property.ImageUrl != nil {
		imageUrl = *property.ImageUrl
	}

	// Assemble the data in a map
	data := map[string]interface{}{
		"name":                 property.Name,
		"address":              property.Address,
		"coordinate":           property.Coordinate,
		"phone1":               property.Phone1,
		"phone2":               property.Phone2,
		"fax":                  property.Fax,
		"email":                property.Email,
		"email_link":           property.EmailLink,
		"billing_email":        property.BillingEmail,
		"technical_email":      property.TechnicalEmail,
		"cc_email":             property.CcEmail,
		"zip":                  property.Zip,
		"city":                 property.City,
		"account":              property.AccountName,
		"acc_no":               property.AccountNo,
		"bank":                 property.Bank,
		"manager":              property.Manager,
		"accounting":           property.Accounting,
		"site_name":            property.SiteName,
		"logo":                 property.Logo,
		"url_upload":           urlUpload,
		"image_url":            imageUrl,
		"meta_desc":            property.MetaDescription,
		"meta_key":             property.MetaKeyword,
		"notif_url":            property.NotifUrl,
		"notif_token":          property.NotifToken,
		"pos_url":              property.PosUrl,
		"pos_token":            property.PosToken,
		"invoice_url":          property.InvoiceUrl,
		"pg_url":               property.PgUrl,
		"pg_token":             property.PgToken,
		"ship_url":             property.ShipUrl,
		"ship_token":           property.ShipToken,
		"courier_integration":  property.CourierIntegration,
		"start":                start,
		"end":                  end,
		"distance_limit":       property.DistanceLimit,
		"shipping_integration": property.ShippingIntegration,
		"shipping_vendor":      strings.ToLower(property.ShippingVendor),
		"point_calculate":      property.PointCalculate,
	}

	return data, nil
}

// TableName sets the name of the database table
func (Property) TableName() string {
	return "property"
}
