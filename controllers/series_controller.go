package controllers

import (
	"net/http"
	"time"
	"fmt"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils" // Import package utils
	"github.com/gin-gonic/gin"
)

// CreateSeries membuat series baru dengan kode verifikasi otomatis
func CreateSeries(c *gin.Context) {
	var input dto.CreateSeriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 1. Generate Kode Unik Otomatis (6 Karakter)
	// Pastikan file utils/random.go sudah Anda buat sebelumnya
	autoCode := utils.GenerateRandomString(6)

	fmt.Println("===========================================")
    fmt.Println("DEBUG: SEDANG MEMBUAT SERIES BARU")
    fmt.Println("DEBUG: KODE YANG DIGENERATE ADALAH:", autoCode)
    fmt.Println("===========================================")

	series := models.Series{
		SeriesName:       input.SeriesName,
		Description:      input.Description,
		TrackID:          input.TrackID,
		Deadline:         input.Deadline,
		IsCompetition:    false, // Default value, bisa disesuaikan jika ada input
		OrderIndex:       0,     // Default value
		
		// 2. Masukkan kode hasil generate ke database
		VerificationCode: autoCode, 
	}

	if err := database.DB.Create(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create series", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusCreated, "Series created successfully", series)
}

// GetSeriesByID mengambil detail series berdasarkan ID
func GetSeriesByID(c *gin.Context) {
	id := c.Param("id")
	var series models.Series

	if err := database.DB.First(&series, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Series not found", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Series fetched successfully", series)
}

// SetSeriesVerificationCode (Opsional) jika Admin ingin mengubah kode secara manual nanti
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

// VerifySeriesCode memverifikasi kode yang dikirim user
func VerifySeriesCode(c *gin.Context) {
	// Ambil ID Series dari URL
	seriesID := c.Param("id")
	
	// Ambil input kode dari Body JSON
	var input struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userID, _ := c.Get("userID")

	// Cek apakah Series ada dan kodenya cocok
	var series models.Series
	if err := database.DB.Where("id = ? AND verification_code = ?", seriesID, input.Code).First(&series).Error; err != nil {
		utils.APIResponse(c, http.StatusBadRequest, "Invalid verification code or series not found", nil)
		return
	}

	// Cek apakah user sudah pernah verifikasi sebelumnya (agar tidak double)
	var existingVerif models.UserSeriesVerification
	if err := database.DB.Where("user_id = ? AND series_id = ?", userID, seriesID).First(&existingVerif).Error; err == nil {
		utils.APIResponse(c, http.StatusOK, "You have already verified this series", nil)
		return
	}

	// Simpan riwayat verifikasi user
	verification := models.UserSeriesVerification{
		UserID:     userID.(uint),
		SeriesID:   series.ID,
		VerifiedAt: time.Now(), // Helper waktu atau time.Now()
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to verify attendance", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Attendance verified successfully", nil)
}