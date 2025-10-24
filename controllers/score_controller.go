package controllers

import (
	"net/http"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

func CreateScore(c *gin.Context) {
	var input dto.CreateScoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	adminID, _ := c.Get("userID")

	var submission models.Submission
	if err := database.DB.First(&submission, input.SubmissionID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Submission not found", nil)
		return
	}

	var score models.Score
	err := database.DB.Where("submission_id = ?", input.SubmissionID).First(&score).Error

	if err == nil {
		score.Value = input.Value
		score.Feedback = input.Feedback
		score.AdminID = adminID.(uint)
		database.DB.Save(&score)

		if err := database.DB.Preload("Admin").First(&score, score.ID).Error; err != nil {
			utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch score after update", err.Error())
			return
		}

		utils.APIResponse(c, http.StatusOK, "Score updated successfully", score)
		return
	}

	score = models.Score{
		Value:        input.Value,
		Feedback:     input.Feedback,
		AdminID:      adminID.(uint),
		SubmissionID: input.SubmissionID,
	}

	if err := database.DB.Create(&score).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create score", err.Error())
		return
	}

	if err := database.DB.Preload("Admin").First(&score, score.ID).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch score after creation", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusCreated, "Score created successfully", score)
}
