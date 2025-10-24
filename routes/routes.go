package routes

import (
	"github.com/Zain0205/gdgoc-subbmission-be-go/controllers"
	"github.com/Zain0205/gdgoc-subbmission-be-go/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	api := r.Group("/api/v1")
	{
		// 1. Auth (Public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 2. Member Routes (Requires Member Role)
		member := api.Group("/member")
		member.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("member"))
		{
			member.POST("/event/join", controllers.JoinEvent)
			member.POST("/submission", controllers.CreateSubmission)
		}

		// 3. Admin Routes (Requires Admin Role)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			admin.POST("/event", controllers.CreateEvent)
			admin.POST("/score", controllers.CreateScore)
			admin.GET("/submissions/event/:eventId", controllers.GetSubmissionsByEvent)
		}

		// 4. Shared Routes (Requires Auth, any role)
		shared := api.Group("/")
		shared.Use(middleware.AuthMiddleware())
		{
			shared.GET("/leaderboard/:eventId", controllers.GetLeaderboard)
		}
	}

	return r
}
