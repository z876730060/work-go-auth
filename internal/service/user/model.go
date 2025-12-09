package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func (u *User) TableName() string {
	return "user"
}

type UserRole struct {
	gorm.Model
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}

func (u *UserRole) TableName() string {
	return "user_role"
}

func InitUserTable(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserRole{})

	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	db.Create(&User{
		Username: "admin",
		Password: "123456",
		Fullname: "Admin",
		Email:    "admin@example.com",
		Phone:    "1234567890",
	})

	db.Create(&UserRole{
		UserID: 1,
		RoleID: 1,
	})
}
