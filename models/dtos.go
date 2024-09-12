package models

import (
	"time"

	"gorm.io/gorm"
)

type ReportStatus string
type TransactionAction string
type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

const (
	Draft     ReportStatus = "draft"
	Completed ReportStatus = "completed"
)

const (
	CreateTx TransactionAction = "create"
	UpdateTx TransactionAction = "update"
	DeleteTx TransactionAction = "delete"
	ReadTx   TransactionAction = "read"
)

type BaseModel struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	CreatedAt   *time.Time      `json:"createdAt,omitempty"`
	CreatedBy   *User           `json:"createdBy,omitempty"`
	CreatedByID *uint           `json:"createdByID,omitempty"`
	DeletedAt   *gorm.DeletedAt `json:"deletedAt,omitempty"`
	UpdatedAt   *time.Time      `json:"updatedAt,omitempty"`
	UpdatedBy   *User           `json:"updatedBy,omitempty"`
	UpdatedByID *uint           `json:"updatedByID,omitempty"`
}

type UserDTO struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type ProductReportDTO struct {
	Description string   `json:"description"`
	Product     *Product `json:"product,omitempty"`
	ProductID   *uint    `json:"productID,omitempty"`
}

type ProductDTO struct {
	Price    uint `json:"price"`
	Quantity uint `json:"quantity"`
}
