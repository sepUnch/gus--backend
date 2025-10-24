package controllers

import (
	"net/http"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

func CreateSubmission(c *gin.Context) {
	var input dto.CreateSubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	memberID, _ := c.Get("userID")

	var count int64
	database.DB.Table("event_members").
		Where("event_id = ? AND user_id = ?", input.EventID, memberID).
		Count(&count)

	if count == 0 {
		utils.APIResponse(c, http.StatusForbidden, "You must join the event before submitting", nil)
		return
	}

	submission := models.Submission{
		Description: input.Description,
		FileURL:     input.FileURL,
		EventID:     input.EventID,
		UserID:      memberID.(uint),
	}

	if err := database.DB.Create(&submission).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create submission", err.Error())
		return
	}

	if err := database.DB.Preload("User").Preload("Event").First(&submission, submission.ID).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch submission after creation", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusCreated, "Submission created successfully", submission)
}

func GetSubmissionsByEvent(c *gin.Context) {
	eventID := c.Param("eventId")

	var submissions []models.Submission
	err := database.DB.Preload("User").Preload("Score").Where("event_id = ?", eventID).Find(&submissions).Error
	if err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Submissions fetched successfully", submissions)
}
