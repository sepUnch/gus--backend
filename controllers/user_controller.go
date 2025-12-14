package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"os"

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
	// FIX 1: Ganti "user_id" menjadi "userID" agar sesuai dengan Middleware
	userID, exists := ctx.Get("userID") 
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
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
			"avatar": user.Avatar, // Pastikan model User punya field Avatar
		},
	})
}

// ================= UPDATE PROFILE (MEMBER ONLY) =================
func (c *UserController) UpdateProfile(ctx *gin.Context) {
    userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
        return
    }

    // --- DEBUG PRINT (Cek di terminal docker nanti) ---
    fmt.Println("--- Menerima Request Update Profile ---")

    // Cari User
    var user models.User
    if err := c.DB.First(&user, userID).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
        return
    }

    // 1. Ambil Form Data (JANGAN PAKAI ShouldBindJSON)
    name := ctx.PostForm("name")
    email := ctx.PostForm("email")

    if name != "" { user.Name = name }
    if email != "" { user.Email = email }

    // 2. Handle Avatar
    file, err := ctx.FormFile("avatar")
    if err == nil {
        fmt.Println("File ditemukan:", file.Filename) // <--- Debug Print

        // Path Folder
        uploadDir := "public/uploads/avatars"
        
        // Pastikan folder ada
        if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
             os.MkdirAll(uploadDir, 0777) // Gunakan 0777 di dalam container
        }

        // Generate nama unik
        ext := filepath.Ext(file.Filename)
        filename := fmt.Sprintf("%d_%d%s", userID.(uint), time.Now().Unix(), ext)
        savePath := filepath.Join(uploadDir, filename)

        // Simpan
        if err := ctx.SaveUploadedFile(file, savePath); err != nil {
            fmt.Println("Gagal simpan file:", err.Error()) // <--- Debug Print Error
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image"})
            return
        }

        // Update DB
        user.Avatar = "uploads/avatars/" + filename
    } else {
        fmt.Println("Tidak ada file avatar yg dikirim/Error:", err.Error()) // <--- Debug Print Error
    }

    // Simpan ke DB
    c.DB.Save(&user)

    // Return Response
    token, _ := utils.GenerateToken(user.ID, user.Role)
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Profile updated successfully",
        "data": gin.H{
            "token": token,
            "user": gin.H{
                "id":     user.ID,
                "name":   user.Name,
                "email":  user.Email,
                "role":   user.Role,
                "avatar": user.Avatar,
            },
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
