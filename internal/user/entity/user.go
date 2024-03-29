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
	EmailVerifiedAt    time.Time `gorm:"default:null"`
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

type VerificationEmail struct {
	ID        uint `gorm:"primaryKey"`
	Email     string
	OTPCode   string
	ExpiredAt time.Time
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

type ReportedThread struct {
	ID                      uint
	ThreadTitle             string
	ThreadBody              string
	ThreadImageURL          string
	ThreadCreatedAt         time.Time
	LikesCount              int
	ReportedUsername        string
	ReportedProfileImageURL string
	ReportedUserProfession  string
	ReportCategory          string
	ReportCreatedAt         time.Time
	Username                string
	ProfileImageURL         string
}

type ReportedComment struct {
	ID                      uint
	CommentBody             string
	LikesCount              int
	CommentCreatedAt        time.Time
	ReportedUsername        string
	ReportedProfileImageURL string
	ReportCategory          string
	ReportCreatedAt         time.Time
	Username                string
	ProfileImageURL         string
}

type ReportedCommunity struct {
	ID                          uint
	CommunityName               string
	CommunityImageURL           string
	CommunityBackgroundImageURL string
	ReportCategory              string
	ReportCreatedAt             time.Time
	Username                    string
	ProfileImageURL             string
}

type ReportedUser struct {
	ID                          uint
	ReportedUserUsername        string
	ReportedUserName            string
	ReportedUserProfession      string
	ReporteduserBio             string
	ReportedUserProfileImageURL string
	ReportedUserBackgroundURL   string
	ReportingUserUsername       string
	ReportingUserName           string
	FollowersCount              int
	FollowingCount              int
}
