package controllers

import (
	"net/http"
	"sort"
	"time"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

// DTO Khusus untuk Activity agar JSON-nya rapi
type ActivityDTO struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`       // Nama Track atau Series
	Type        string    `json:"type"`        // "track" atau "series"
	Action      string    `json:"action"`      // "New track added" atau "Updated"
	Time        time.Time `json:"time"`
}

// GetDashboardData (Statistik + Recent Activity)
func GetDashboardData(c *gin.Context) {
	// 1. Ambil Statistik Angka (Counts)
	var totalTracks, totalSeries, totalUsers int64

	// Count Data dari DB
	database.DB.Model(&models.Track{}).Count(&totalTracks)
	database.DB.Model(&models.Series{}).Count(&totalSeries)
	database.DB.Model(&models.User{}).Count(&totalUsers)

	// 2. Ambil Recent Activity (Gabungan Tracks & Series)
	var activities []ActivityDTO

	// A. Ambil 5 Track Terakhir (Updated/Created)
	var recentTracks []models.Track
	database.DB.Order("updated_at desc").Limit(5).Find(&recentTracks)

	for _, t := range recentTracks {
		action := "Updated"
		// Logic sederhana: jika selisih created & updated < 10 detik, anggap baru dibuat
		if t.CreatedAt.Equal(t.UpdatedAt) || t.UpdatedAt.Sub(t.CreatedAt).Seconds() < 10 {
			action = "New track added"
		} else {
			action = "Track updated"
		}

		activities = append(activities, ActivityDTO{
			ID:     t.ID,
			Title:  t.TrackName,
			Type:   "track",
			Action: action,
			Time:   t.UpdatedAt,
		})
	}

	// B. Ambil 5 Series Terakhir
	var recentSeries []models.Series
	database.DB.Order("updated_at desc").Limit(5).Find(&recentSeries)

	for _, s := range recentSeries {
		action := "Updated"
		if s.CreatedAt.Equal(s.UpdatedAt) || s.UpdatedAt.Sub(s.CreatedAt).Seconds() < 10 {
			action = "New series added"
		} else {
			action = "Series updated"
		}

		activities = append(activities, ActivityDTO{
			ID:     s.ID,
			Title:  s.SeriesName,
			Type:   "series",
			Action: action,
			Time:   s.UpdatedAt,
		})
	}

	// 3. Sorting Gabungan (Yang paling baru di atas)
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Time.After(activities[j].Time)
	})

	// 4. Ambil 5 Teratas saja agar tidak kepanjangan
	if len(activities) > 5 {
		activities = activities[:5]
	}

	// 5. Response JSON
	utils.APIResponse(c, http.StatusOK, "Dashboard data fetched", gin.H{
		"counts": gin.H{
			"tracks": totalTracks,
			"series": totalSeries,
			"users":  totalUsers,
			"active": totalUsers, // Sementara disamakan, nanti bisa difilter by 'is_active'
		},
		"activities": activities,
	})
}