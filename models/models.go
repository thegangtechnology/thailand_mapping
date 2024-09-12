package models

type User struct {
	BaseModel `yaml:",inline"`
	UserDTO   `yaml:",inline"`
}

type Product struct {
	BaseModel  `yaml:",inline"`
	ProductDTO `yaml:",inline"`
}

type CreateProductInput struct {
	ProductDTO `yaml:",inline"`
}

type ProductReport struct {
	BaseModel        `yaml:",inline"`
	ProductReportDTO `yaml:",inline"`
}

type CreateProductReportInput struct {
	ProductReportDTO `yaml:",inline"`
}

type Message struct {
	Code    int64
	Message string
}

type Permissions struct {
	List *[]Permission `json:"list,omitempty"`
}

type Permission struct {
	Action    TransactionAction `gorm:"type:transaction_action_enum;default:create" json:"action"`
	IsGranted bool              `json:"isGranted"`
}

type Roles struct {
	Roles *[]Role `json:"roles,omitempty"`
}
