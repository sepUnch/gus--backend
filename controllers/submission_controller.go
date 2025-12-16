package controllers

import (
	"errors"
	"net/http"
	"time" // Jangan lupa import time

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateSubmission
func CreateSubmission(c *gin.Context) {
	var input dto.CreateSubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	memberID, _ := c.Get("userID")

	// --- LOGIC BARU MULAI DI SINI ---

	// 1. Ambil Data Series Dulu
	var series models.Series
	if err := database.DB.First(&series, input.SeriesID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Series not found", nil)
		return
	}

	// 2. Cek apakah Series dinonaktifkan Admin (Manual Close)
	if !series.IsActive {
		utils.APIResponse(c, http.StatusForbidden, "Submissions for this series have been closed by the admin.", nil)
		return
	}

	// 3. Cek Deadline (Realtime Check)
	// Jika waktu sekarang (Now) SUDAH MELEWATI (After) Deadline
	if time.Now().After(series.Deadline) {
		utils.APIResponse(c, http.StatusForbidden, "The deadline for this series has passed.", nil)
		return
	}

	var existingSubmission models.Submission
	// Cek apakah user ini (memberID) sudah punya submission di series ini (input.SeriesID)
	if err := database.DB.Where("user_id = ? AND series_id = ?", memberID, input.SeriesID).First(&existingSubmission).Error; err == nil {
		// Jika err == nil, artinya DATA DITEMUKAN -> Berarti Duplikat!
		utils.APIResponse(c, http.StatusConflict, "You have already submitted for this series.", nil)
		return
	}

	// --- LOGIC LAMA (Verifikasi Absen) ---

	var verification models.UserSeriesVerification
	err := database.DB.Where("series_id = ? AND user_id = ?", input.SeriesID, memberID).First(&verification).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.APIResponse(c, http.StatusForbidden, "You must verify your attendance for this series before submitting", nil)
		return
	} else if err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Error checking verification", err.Error())
		return
	}

	// --- PROSES SIMPAN SUBMISSION ---

	submission := models.Submission{
		SeriesID: input.SeriesID,
		FileURL:  input.FileURL,
		UserID:   memberID.(uint),
		Score:    0, // Default score
	}

	if err := database.DB.Create(&submission).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create submission", err.Error())
		return
	}

	database.DB.Preload("User").Preload("Series").First(&submission, submission.ID)
	utils.APIResponse(c, http.StatusCreated, "Submission created successfully", submission)
}

// ... Fungsi GetSubmissionsBySeries dan GetSubmissionByID biarkan tetap sama ...
// (Salin saja kode lama untuk fungsi Get di bawah ini)
func GetSubmissionsBySeries(c *gin.Context) {
	seriesID := c.Param("seriesId")

	var submissions []models.Submission
	err := database.DB.Preload("User").Where("series_id = ?", seriesID).Find(&submissions).Error
	if err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Submissions fetched successfully", submissions)
}

func GetSubmissionByID(c *gin.Context) {
	id := c.Param("id")
	var submission models.Submission

	if err := database.DB.Preload("User").Preload("Series").First(&submission, id).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Submission not found", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Submission fetched successfully", submission)
}
