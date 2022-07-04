package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Books struct {
	ID        uint    `json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

type Users struct {
	ID       uint    `json:"id"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
}

func Migrate(db *gorm.DB) error {
	fmt.Println("Migration Function Called")
	err := db.AutoMigrate(&Books{}, &Users{})
	return err
}
