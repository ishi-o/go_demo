package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	api "github.com/ishi-o/go_demo/hertz_demo/biz/model/api"
)

type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	Username     string         `gorm:"uniqueIndex;size:64;not null"`
	Email        string         `gorm:"uniqueIndex;size:128;not null"`
	Phone        string         `gorm:"size:20"`
	Age          int32          `gorm:"default:0"`
	Gender       int32          `gorm:"default:0"`
	Status       int32          `gorm:"default:1;index"`
	PasswordHash string         `gorm:"size:256"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Token        string         `gorm:"-"`
	Permissions  []string       `gorm:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	if u.PasswordHash == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) IsActive() bool {
	return u.Status == int32(api.UserStatus_USER_STATUS_ACTIVE)
}

func (u *User) CanLogin() bool {
	return u.IsActive()
}

func (u *User) ToDTO() *api.User {
	return &api.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Age:       u.Age,
		Gender:    api.Gender(u.Gender),
		Status:    api.UserStatus(u.Status),
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}

func UserFromCreateRequest(req *api.CreateUserRequest) *User {
	return &User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Age:      req.Age,
		Gender:   int32(req.Gender),
		Status:   int32(api.UserStatus_USER_STATUS_ACTIVE),
	}
}

func (u *User) FromUpdateRequest(req *api.UpdateUserRequest) {
	if req.Username != "" {
		u.Username = req.Username
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.Phone != "" {
		u.Phone = req.Phone
	}
	if req.Age > 0 {
		u.Age = req.Age
	}
	if req.Gender != api.Gender_GENDER_UNKNOWN_UNSPECIFIED {
		u.Gender = int32(req.Gender)
	}
	if req.Status != api.UserStatus_USER_STATUS_UNKNOWN_UNSPECIFIED {
		u.Status = int32(req.Status)
	}
	u.UpdatedAt = time.Now()
}
