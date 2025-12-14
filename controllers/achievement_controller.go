package controllers

import (
	"net/http"
	"time" // Pastikan paket time diimpor

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

// CreateAchievementType creates a new achievement type
func CreateAchievementType(c *gin.Context) {
	var input dto.CreateAchievementTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	achievType := models.AchievementType{Name: input.Name}
	if err := database.DB.Create(&achievType).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create achievement type", err.Error())
		return
	}
	utils.APIResponse(c, http.StatusCreated, "Achievement type created", achievType)
}

// GetAchievementTypes fetches all achievement types
func GetAchievementTypes(c *gin.Context) {
	var types []models.AchievementType
	database.DB.Find(&types)
	utils.APIResponse(c, http.StatusOK, "Achievement types fetched", types)
}

// CreateAchievement creates a new achievement
func CreateAchievement(c *gin.Context) {
	var input dto.CreateAchievementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	achiev := models.Achievement{
		Name:              input.Name,
		Description:       input.Description,
		IconURL:           input.IconURL,
		AchievementTypeID: input.AchievementTypeID,
	}

	if err := database.DB.Create(&achiev).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create achievement", err.Error())
		return
	}

	database.DB.Preload("Type").First(&achiev, achiev.ID)

	utils.APIResponse(c, http.StatusCreated, "Achievement created", achiev)
}

// GetAchievements fetches all achievements
func GetAchievements(c *gin.Context) {
	var achievs []models.Achievement
	database.DB.Preload("Type").Find(&achievs)
	utils.APIResponse(c, http.StatusOK, "Achievements fetched", achievs)
}

// UpdateAchievement updates an existing achievement
func UpdateAchievement(c *gin.Context) {
	id := c.Param("id")
	var achiev models.Achievement
	if err := database.DB.First(&achiev, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Achievement not found", nil)
		return
	}

	var input dto.UpdateAchievementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := database.DB.Model(&achiev).Updates(input).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to update achievement", err.Error())
		return
	}

	database.DB.Preload("Type").First(&achiev, achiev.ID)
	utils.APIResponse(c, http.StatusOK, "Achievement updated", achiev)
}

// AwardAchievementToUser gives an achievement to a user
func AwardAchievementToUser(c *gin.Context) {
	var input dto.AwardAchievementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, input.UserID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	var achiev models.Achievement
	if err := database.DB.First(&achiev, input.AchievementID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Achievement not found", nil)
		return
	}

	userAchiev := models.UserAchievement{
		UserID:        input.UserID,
		AchievementID: input.AchievementID,
		EarnedAt:      time.Now(), // MENGGUNAKAN time.Now() YANG BENAR
	}

	// Menggunakan FirstOrCreate untuk menghindari duplikat jika itu perilaku yang diinginkan,
	// atau tetap gunakan Create dan tangani error constraint unik.
	// Berdasarkan snippet yang Anda berikan, FirstOrCreate digunakan.
	if err := database.DB.FirstOrCreate(&userAchiev, models.UserAchievement{UserID: input.UserID, AchievementID: input.AchievementID}).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to award achievement", err.Error())
		return
	}

	// Ambil ulang data untuk mendapatkan preload yang lengkap
	// Catatan: FirstOrCreate mungkin tidak mengembalikan objek yang di-preload sepenuhnya secara langsung jika ia menemukan objek yang sudah ada
	// Jadi kita query lagi untuk memastikan kita mengembalikan detail lengkap.
	database.DB.Preload("User").Preload("Achievement.Type").First(&userAchiev, "user_id = ? AND achievement_id = ?", userAchiev.UserID, userAchiev.AchievementID)
	
	utils.APIResponse(c, http.StatusCreated, "Achievement awarded", userAchiev)
}

// RevokeAchievementFromUser removes an achievement from a user
func RevokeAchievementFromUser(c *gin.Context) {
	var input dto.AwardAchievementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result := database.DB.Delete(&models.UserAchievement{}, "user_id = ? AND achievement_id = ?", input.UserID, input.AchievementID)
	if result.Error != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Error revoking achievement", result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		utils.APIResponse(c, http.StatusNotFound, "User does not have this achievement", nil)
		return
	}

	utils.APIResponse(c, http.StatusOK, "Achievement revoked", nil)
}

// GetAchievementCount counts all achievements
func GetAchievementCount(c *gin.Context) {
    var count int64

    // Menghitung total baris di tabel achievements
    if err := database.DB.Model(&models.Achievement{}).Count(&count).Error; err != nil {
        utils.APIResponse(c, http.StatusInternalServerError, "Failed to count achievements", err.Error())
        return
    }

    // Mengirimkan angka count ke frontend
    utils.APIResponse(c, http.StatusOK, "Achievement count fetched", gin.H{
        "count": count,
    })
}