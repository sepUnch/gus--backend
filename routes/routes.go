package routes

import (
	"github.com/Zain0205/gdgoc-subbmission-be-go/controllers"
	"github.com/Zain0205/gdgoc-subbmission-be-go/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	userController := controllers.UserController{DB: db}

	api := r.Group("/api")
	{
		// ================= AUTH =================
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// ================= UNIVERSAL (LOGIN REQUIRED) =================
		api.GET("/me", middleware.AuthMiddleware(), userController.GetProfile)

		// ================= PUBLIC =================
		public := api.Group("/")
		public.Use(middleware.AuthMiddleware())
		{
			public.GET("/tracks", controllers.GetAllTracks)
			public.GET("/tracks/:id", controllers.GetTrackWithSeries)
			public.GET("/leaderboard/track/:trackId", controllers.GetLeaderboardByTrack)
		}

		// ================= MEMBER =================
		member := api.Group("/member")
		member.Use(
			middleware.AuthMiddleware(),
			middleware.RoleMiddleware("member"),
		)
		{
			member.GET("/me", userController.GetProfile)
			member.PUT("/me", userController.UpdateProfile)
			member.POST("/submissions", controllers.CreateSubmission)
			member.POST("/series/:id/verify", controllers.VerifySeriesCode)
			member.GET("/me/achievements", controllers.GetMyAchievements)
		}

		// ================= ADMIN =================
		admin := api.Group("/admin")
		admin.Use(
			middleware.AuthMiddleware(),
			middleware.RoleMiddleware("admin"),
		)
		{
			admin.GET("/dashboard", controllers.GetDashboardData) 
			admin.GET("/stats", controllers.GetUserCount) // Route lama (opsional dihapus jika sudah diganti dashboard)

			admin.GET("/users", controllers.GetAllUsers)       // <--- List Users
			admin.DELETE("/users/:id", controllers.DeleteUser) // <--- Delete User

			admin.POST("/tracks", controllers.CreateTrack)
			admin.GET("/tracks/:id", controllers.GetTrackWithSeries)
			admin.PUT("/tracks/:id", controllers.UpdateTrack)    // <--- Tambahkan ini
			admin.DELETE("/tracks/:id", controllers.DeleteTrack) // <--- Tambahkan ini
			admin.POST("/series", controllers.CreateSeries)
			admin.GET("/series/:id", controllers.GetSeriesByID)
			admin.PUT("/series/:id", controllers.UpdateSeries)    // Untuk Edit
			admin.DELETE("/series/:id", controllers.DeleteSeries) // Untuk Hapus
			admin.PATCH("/series/:id/toggle", controllers.ToggleSeriesStatus)
			admin.PATCH("/series/:id/code", controllers.SetSeriesVerificationCode)

			admin.GET("/submissions/:id", controllers.GetSubmissionByID)
			admin.GET("/submissions/series/:seriesId", controllers.GetSubmissionsBySeries)
			admin.POST("/submissions/grade", controllers.GradeSubmission)

			admin.PATCH("/users/:id/role", controllers.SetUserRole)

			admin.POST("/achievement-types", controllers.CreateAchievementType)
			admin.GET("/achievement-types", controllers.GetAchievementTypes)

			admin.POST("/achievements", controllers.CreateAchievement)
			admin.GET("/achievements", controllers.GetAchievements)
			admin.PUT("/achievements/:id", controllers.UpdateAchievement)
			admin.GET("/stats/achievements", controllers.GetAchievementCount)

			admin.POST("/achievements/award", controllers.AwardAchievementToUser)
			admin.POST("/achievements/revoke", controllers.RevokeAchievementFromUser)
		}
	}
}
