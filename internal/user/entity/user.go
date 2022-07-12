package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string `gorm:"uniqueIndex;size:75"`
	Username           string `gorm:"uniqueIndex;size:50"`
	Password           string
	Name               string
	ProfileImageUrl    string
	BackgroundImageUrl string
	Bio                string
	Profession         string
	Role               string
	IsBanned           int
	Followers          []User       `gorm:"many2many:user_followers"`
	Report             []UserReport `gorm:"foreignKey:UserID"`
	Reported           []UserReport `gorm:"foreignKey:ReportedUserID"`
	IsFollowed         int          `gorm:"-:migration;<-:false"`
	IsMine             int          `gorm:"-:migration;<-:false"`
}

type UserReport struct {
	gorm.Model
	UserID           uint
	ReportedUserID   uint
	ReportCategoryID uint
}

type BriefReport struct {
	ThreadReportsID     uint
	UserReportsID       uint
	CommentReportsID    uint
	CommunityReportsID  uint
	CreatedAt           time.Time
	ThreadID            uint
	UserID              uint
	CommentID           uint
	CommunityReportedID uint
	ReportCategory      string
	Username            string
	ProfileImageURL     string
	Type                string
}

type AdminDashboardAnalytics struct {
	UsersCount      int
	ModeratorsCount int
	ReportsCount    int
}
