package controllers

import (
	"net/http"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/dto"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// ================= GET PROFILE (UNIVERSAL) =================
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// ================= UPDATE PROFILE (MEMBER ONLY) =================
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Email unik
	if req.Email != user.Email {
		var existing models.User
		if err := c.DB.Where("email = ? AND id != ?", req.Email, userID).
			First(&existing).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email already taken"})
			return
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := c.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// ================= ADMIN =================
func SetUserRole(c *gin.Context) {
	userID := c.Param("id")
	var input dto.UpdateUserRoleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	user.Role = input.Role
	if err := database.DB.Save(&user).Error; err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to update role", nil)
		return
	}

	utils.APIResponse(c, http.StatusOK, "Role updated", user)
}

func GetUserCount(c *gin.Context) {
	var count int64

	if err := database.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to count users",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users": count,
		},
	})
}
