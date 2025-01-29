package models

import (
	"time"

	config "go-htmlcsstoimage/configs"
)

type Image struct {
	ID        int       `gorm:"primaryKey;not null;autoIncrement"`
	UserID    int       `gorm:"not null"`
	ImageName string    `gorm:"type:uuid;unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	DeletedAt time.Time `gorm:"index"`
	Type      string    `gorm:"type:varchar(10)"`
}

type GenerateImageRequest struct {
	HTML            string  `json:"html" form:"html"`
	CSS             string  `json:"css" form:"css"`
	URL             string  `json:"url" form:"url"`
	GoogleFonts     string  `json:"google_fonts" form:"google_fonts"`
	Selector        string  `json:"selector" form:"selector"`
	MsDelay         int     `json:"ms_delay" form:"ms_delay"`
	DeviceScale     float64 `json:"device_scale" form:"device_scale"`
	RenderWhenReady bool    `json:"render_when_ready" form:"render_when_ready"`
	FullScreen      bool    `json:"full_screen" form:"full_screen"`
	ViewportWidth   int64   `json:"viewport_width" form:"viewport_width"`
	ViewportHeight  int64   `json:"viewport_height" form:"viewport_height"`
}

type GenerateImageResponse struct {
	URL string `json:"url"`
}

// SaveImageToDatabase menyimpan informasi gambar ke database
func SaveImageToDatabase(image Image) error {
	query := `
		INSERT INTO images (user_id, image_name, created_at, type)
		VALUES ($1, $2, $3, $4)
	`
	err := config.DB.Exec(query, image.UserID, image.ImageName, image.CreatedAt, image.Type).Error
	return err
}
