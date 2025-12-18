package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Zain0205/gdgoc-subbmission-be-go/database"
	"github.com/Zain0205/gdgoc-subbmission-be-go/models"
	"github.com/Zain0205/gdgoc-subbmission-be-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func DownloadCertificate(c *gin.Context) {
	trackID := c.Param("id")
	userID, _ := c.Get("userID")

	// 1. Ambil Data
	var track models.Track
	if err := database.DB.First(&track, trackID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "Track not found", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.APIResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	// 2. Setup PDF (Landscape A4)
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.SetAutoPageBreak(false, 0)

	pdf.AddPage()

	// 3. Load Template
	imagePath := "./public/certificate-template.png"
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		utils.APIResponse(c, http.StatusInternalServerError, "Template missing", nil)
		return
	}
	pdf.Image(imagePath, 0, 0, 297, 210, false, "", 0, "")

	// ==========================================
	// TITIK KOORDINAT PRESISI (Disesuaikan dengan Lingkaran Merah Mas)
	// ==========================================

	// A. NAMA USER (Di Garis "Presented To")
	// X=95: Geser ke kanan supaya TIDAK menimpa tulisan "Presented to"
	// Y=102: Pas di atas garis tengah
	pdf.SetFont("Helvetica", "B", 26)
	pdf.SetTextColor(50, 50, 50)
	pdf.SetXY(95, 97)
	pdf.Cell(150, 10, strings.ToUpper(user.Name))

	// B. NAMA TRACK (Di Tengah Kosong)
	// Y=138: Di area kosong tengah
	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(66, 133, 244) // Biru Google

	pageWidth, _ := pdf.GetPageSize()
	trackWidth := pdf.GetStringWidth(track.TrackName)
	xTrack := (pageWidth - trackWidth) / 2
	pdf.SetXY(xTrack, 138)
	pdf.Cell(trackWidth, 10, track.TrackName)

	// C. TANGGAL (Di Garis "DATE" - Kiri Bawah)
	// X=40: Geser dikit dari margin kiri supaya pas di atas garis
	// Y=176: Pas menempel di garis bawah
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(80, 80, 80)
	dateText := time.Now().Format("January 02, 2006")

	pdf.SetXY(23, 174)
	pdf.Cell(40, 10, dateText)

	// D. LOKASI (Di Garis "LOCATION" - Tengah Bawah)
	// X=115: Di sebelah kanan tanggal
	// Y=176: Sejajar dengan tanggal
	pdf.SetXY(75, 174)
	pdf.Cell(80, 10, "Universitas Esa Unggul")

	// 4. Output
	fileName := fmt.Sprintf("Certificate-%s.pdf", trackID)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/pdf")

	err := pdf.Output(c.Writer)
	if err != nil {
		utils.APIResponse(c, http.StatusInternalServerError, "Failed to generate PDF", err.Error())
	}
}
