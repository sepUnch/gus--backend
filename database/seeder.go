package database

import (
	"log"

	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"gorm.io/gorm"
)

// SeedAchievementTypes mengisi data awal tipe achievement jika kosong
func SeedAchievementTypes(db *gorm.DB) {
	var count int64
	db.Model(&models.AchievementType{}).Count(&count)

	// Hanya isi jika tabel kosong
	if count == 0 {
		types := []models.AchievementType{
			{Name: "TRACK_COMPLETION"},
			{Name: "SERIES_COMPLETION"},
			{Name: "MANUAL_AWARD"},
		}

		if err := db.Create(&types).Error; err != nil {
			log.Println("Gagal seeding achievement types:", err)
		} else {
			log.Println("Berhasil seeding achievement types!")
		}
	}
}