package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

// CreateSeries membuat series baru
func CreateSeries(c *gin.Context) {
	var input dto.CreateSeriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	autoCode := utils.GenerateRandomString(6)

	fmt.Println("DEBUG: CREATING SERIES WITH CODE:", autoCode)

	series := models.Series{
		SeriesName:       input.SeriesName,
		Description:      input.Description,
		TrackID:          input.TrackID,
		Deadline:         input.Deadline,
		IsCompetition:    false,
		OrderIndex:       0,
		VerificationCode: autoCode,
		IsActive:         true, // Explicitly set active on create
	}

	if err := database.DB.Create(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create series", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusCreated, "Series created successfully", series)
}

func ToggleSeriesStatus(c *gin.Context) {
	id := c.Param("id")
	var series models.Series

	// Cari series
	if err := database.DB.First(&series, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Series not found", nil)
		return
	}

	// Balik statusnya (True -> False, False -> True)
	series.IsActive = !series.IsActive

	// Simpan perubahan
	if err := database.DB.Save(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to update series status", err.Error())
		return
	}

	statusMsg := "closed"
	if series.IsActive {
		statusMsg = "opened"
	}

	utils.APIResponse(c, http.StatusOK, "Series successfully "+statusMsg, series)
}

// -------------------------------------------------------------

func GetSeriesByID(c *gin.Context) {
	id := c.Param("id")
	var series models.Series

	if err := database.DB.First(&series, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Series not found", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Series fetched successfully", series)
}

func SetSeriesVerificationCode(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		VerificationCode string `json:"verification_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var series models.Series
	if err := database.DB.First(&series, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Series not found", nil)
		return
	}

	series.VerificationCode = input.VerificationCode
	if err := database.DB.Save(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to update code", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Verification code updated", series)
}

func VerifySeriesCode(c *gin.Context) {
	seriesID := c.Param("id")

	var input struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, _ := c.Get("userID")

	var series models.Series
	if err := database.DB.Where("id = ? AND verification_code = ?", seriesID, input.Code).First(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusBadRequest, "Invalid verification code or series not found", nil)
		return
	}

	var existingVerif models.UserSeriesVerification
	if err := database.DB.Where("user_id = ? AND series_id = ?", userID, seriesID).First(&existingVerif).Error; err == nil {
		utils.APIResponse(c, http.StatusOK, "You have already verified this series", nil)
		return
	}

	verification := models.UserSeriesVerification{
		UserID:     userID.(uint),
		SeriesID:   series.ID,
		VerifiedAt: time.Now(),
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to verify attendance", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Attendance verified successfully", nil)
}

func UpdateSeries(c *gin.Context) {
	id := c.Param("id")
	var series models.Series

	// 1. Cari Series dulu
	if err := database.DB.First(&series, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Series not found", nil)
		return
	}

	// 2. Bind Input (Kita pakai struct inline atau DTO yang sama dgn create)
	var input struct {
		SeriesName  string    `json:"series_name"`
		Description string    `json:"description"`
		TrackID     uint      `json:"track_id"`
		Deadline    time.Time `json:"deadline"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 3. Update Field
	series.SeriesName = input.SeriesName
	series.Description = input.Description
	series.TrackID = input.TrackID
	series.Deadline = input.Deadline

	// 4. Save
	if err := database.DB.Save(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to update series", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Series updated successfully", series)
}

// [BARU] Delete Series (Soft Delete)
func DeleteSeries(c *gin.Context) {
	id := c.Param("id")
	var series models.Series

	// 1. Cek keberadaan data
	if err := database.DB.First(&series, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Series not found", nil)
		return
	}

	// 2. Hapus (GORM akan lakukan Soft Delete karena ada deleted_at di struct)
	if err := database.DB.Delete(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to delete series", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Series deleted successfully", nil)
}
