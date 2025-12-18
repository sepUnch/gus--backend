package controllers

import (
	"net/http"
	"strconv"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database" // Sesuaikan import
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

// GetCommentsBySeries: Ambil semua komentar di satu series
func GetCommentsBySeries(c *gin.Context) {
	seriesID := c.Param("id") // Mengambil ID series dari URL
	var comments []models.Comment

	// Ambil komentar urut dari yang paling baru (descending)
	// Preload("User") penting agar kita dapat nama & avatar pengomentar
	if err := database.DB.Preload("User").Where("series_id = ?", seriesID).Order("created_at desc").Find(&comments).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch comments", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Comments fetched successfully", comments)
}

func CreateComment(c *gin.Context) {
	seriesIDStr := c.Param("id") // Ambil ID sebagai string

	// 1. Validasi User Login
	userID, exists := c.Get("userID")
	if !exists {
		utils.APIResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// 2. Validasi Input JSON
	var input dto.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	// 3. Konversi Series ID dari String ke Uint
	// Kita butuh konversi ini agar bisa masuk ke struct models.Comment
	seriesID, err := strconv.Atoi(seriesIDStr)
	if err != nil {
		utils.APIResponse(c, http.StatusBadRequest, "Invalid Series ID", nil)
		return
	}

	// 4. Buat Object Comment (Variable ini sekarang AKAN DIPAKAI di bawah)
	comment := models.Comment{
		UserID:   userID.(uint),
		SeriesID: uint(seriesID), // Hasil konversi
		Content:  input.Content,
	}

	// 5. Simpan menggunakan GORM (Create)
	// Ini akan otomatis mengisi CreatedAt, UpdatedAt, dan menghandle query
	if err := database.DB.Create(&comment).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to post comment", err.Error())
		return
	}

	// Preload user agar data balikan API langsung memuat nama & avatar (untuk update UI instant)
	database.DB.Preload("User").First(&comment, comment.ID)

	utils.APIResponse(c, http.StatusCreated, "Comment posted successfully", comment)
}
