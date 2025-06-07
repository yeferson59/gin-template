// Package models contiene las definiciones de los modelos de datos de la aplicación.
package models

import (
	"time"

	"gorm.io/gorm"
)

// User representa el modelo de usuario para autenticación y ejemplo.
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName permite personalizar el nombre de la tabla si se desea.
// Por defecto será "users".
func (User) TableName() string {
	return "users"
}
