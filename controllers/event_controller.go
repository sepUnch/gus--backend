package controllers

import (
	"net/http"
	"strings"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	var input dto.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	adminID, _ := c.Get("userID")

	event := models.Event{
		Name:        input.Name,
		CreatedByID: adminID.(uint),
	}

	if err := database.DB.Create(&event).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to create event", err.Error())
		return
	}

	if err := database.DB.Preload("CreatedBy").First(&event, event.ID).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to fetch event after creation", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusCreated, "Event created successfully", event)
}

func JoinEvent(c *gin.Context) {
	var input dto.JoinEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	memberID, _ := c.Get("userID")

	var event models.Event
	if err := database.DB.Where("event_code = ?", input.EventCode).First(&event).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Event code not found", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, memberID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	if err := database.DB.Model(&event).Association("Members").Append(&user); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			utils.APIResponse(c, http.StatusOK, "You have already joined this event", nil)
			return
		}
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to join event", err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Successfully joined event", event)
}
