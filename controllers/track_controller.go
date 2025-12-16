package controllers

import (
    "net/http"

    "github.com/Zain0205/gdgoc-subbmission-be-go/database"
    "github.com/Zain0205/gdgoc-subbmission-be-go/dto"
    "github.com/Zain0205/gdgoc-subbmission-be-go/models"
    "github.com/Zain0205/gdgoc-subbmission-be-go/utils"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm" // Pastikan import ini ada untuk fungsi sorting Preload
)

// CreateTrack
func CreateTrack(c *gin.Context) {
    var input dto.CreateTrackInput
    if err := c.ShouldBindJSON(&input); err != nil {
        utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
        return
    }

    adminID, _ := c.Get("userID")

    track := models.Track{
        TrackName:   input.TrackName,
        Description: input.Description,
        CreatedByID: adminID.(uint),
        TrackType:   input.TrackType,
    }

    // Default value jika kosong
    if track.TrackType == "" {
        track.TrackType = "STUDY_JAM"
    }

    if err := database.DB.Create(&track).Error; err != nil {
        utils.APIResponse(c, http.StatusInternalServerError, "Failed to create track", err.Error())
        return
    }

    // Load data admin pembuatnya untuk response
    database.DB.Preload("CreatedBy").First(&track, track.ID)
    
    utils.APIResponse(c, http.StatusCreated, "Track created successfully", track)
}

// GetAllTracks (Untuk List Halaman Depan)
func GetAllTracks(c *gin.Context) {
    var tracks []models.Track
    
    // Preload Series agar kita bisa menghitung jumlahnya di Frontend (track.series.length)
    // Preload CreatedBy agar tau siapa yang buat
    if err := database.DB.Preload("Series").Preload("CreatedBy").Find(&tracks).Error; err != nil {
        utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch tracks", err.Error())
        return
    }
    
    utils.APIResponse(c, http.StatusOK, "Tracks fetched successfully", tracks)
}

// GetTrackWithSeries (Untuk Halaman Detail Track)
func GetTrackWithSeries(c *gin.Context) {
    trackID := c.Param("id")
    var track models.Track

    // LOGIC PENTING:
    // Kita melakukan Preload "Series" tapi dengan kondisi ORDER BY.
    // Agar series muncul urut berdasarkan order_index, lalu ID.
    if err := database.DB.
        Preload("Series", func(db *gorm.DB) *gorm.DB {
            return db.Order("order_index asc, id asc")
        }).
        Preload("CreatedBy").
        First(&track, trackID).Error; err != nil {
        
        utils.APIResponse(c, http.StatusNotFound, "Track not found", err.Error())
        return
    }

    utils.APIResponse(c, http.StatusOK, "Track fetched successfully", track)
}