package models

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"type:enum('admin','member');default:'member'"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return
}

type Event struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Name        string  `json:"name"`
	EventCode   string  `json:"event_code" gorm:"unique;not null"`
	CreatedByID uint    `json:"created_by_id"`
	CreatedBy   User    `json:"created_by" gorm:"foreignKey:CreatedByID"`
	Members     []*User `json:"members,omitempty" gorm:"many2many:event_members;"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	e.EventCode = generateRandomCode(6)
	return
}

type Submission struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Description string `json:"description"`
	FileURL     string `json:"file_url"`
	UserID      uint   `json:"user_id"`
	EventID     uint   `json:"event_id"`

	User  User  `json:"user" gorm:"foreignKey:UserID"`
	Event Event `json:"event" gorm:"foreignKey:EventID"`
	Score Score `json:"score,omitempty" gorm:"foreignKey:SubmissionID"`
}

type Score struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Value        int    `json:"value"`
	Feedback     string `json:"feedback"`
	AdminID      uint   `json:"admin_id"`
	SubmissionID uint   `json:"submission_id"`

	Admin User `json:"admin" gorm:"foreignKey:AdminID"`
}

func generateRandomCode(n int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
