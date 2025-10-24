package dto

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin member"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateEventInput struct {
	Name string `json:"name" binding:"required"`
}

type JoinEventInput struct {
	EventCode string `json:"event_code" binding:"required"`
}

type CreateSubmissionInput struct {
	Description string `json:"description" binding:"required"`
	FileURL     string `json:"file_url" binding:"required"`
	EventID     uint   `json:"event_id" binding:"required"`
}

type CreateScoreInput struct {
	SubmissionID uint   `json:"submission_id" binding:"required"`
	Value        int    `json:"value" binding:"required,min=0,max=100"`
	Feedback     string `json:"feedback"`
}

type LeaderboardResult struct {
	UserID       uint    `json:"user_id"`
	Name         string  `json:"name"`
	AverageScore float64 `json:"average_score"`
}
