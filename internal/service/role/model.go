package role

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name    string `json:"name" gorm:"unique not null"`
	Key     string `json:"key" gorm:"unique not null"`
	Comment string `json:"comment" gorm:"not null"`
}

type RoleMenu struct {
	gorm.Model
	Rid     uint   `json:"rid" gorm:"index"`
	MenuKey string `json:"menuKey" gorm:"index"`
}

type RoleTree struct {
	Title    string      `json:"title"`
	Key      string      `json:"key"`
	Children []*RoleTree `json:"children"`
}

func (RoleMenu) TableName() string {
	return "role_menu"
}

func (Role) TableName() string {
	return "role"
}

func InitRoleTable(db *gorm.DB) {
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&RoleMenu{})

	var count int64
	db.Model(&Role{}).Count(&count)
	if count > 0 {
		return
	}

	db.Create(&Role{
		Name:    "admin",
		Key:     "admin",
		Comment: "平台内最大权限",
	})
	db.Create(&Role{
		Name: "user",
		Key:  "user",
	})
}
