package controllers

import (
	"net/http"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context) {
	eventID := c.Param("eventId")

	var results []dto.LeaderboardResult

	err := database.DB.Table("users").
		Select("users.id as user_id, users.name, AVG(scores.value) as average_score").
		Joins("JOIN submissions ON users.id = submissions.user_id").
		Joins("JOIN scores ON submissions.id = scores.submission_id").
		Where("submissions.event_id = ?", eventID).
		Group("users.id, users.name").
		Order("average_score DESC").
		Scan(&results).Error

	if err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to generate leaderboard", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Leaderboard fetched successfully", results)
}
