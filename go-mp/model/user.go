package model

type User struct {
    UserID        uint    `gorm:"primaryKey"`
    Email         string  `gorm:"not null"`
    Password      string  `gorm:"not null"`
    DepositAmount float64
}

type RegisterRequestBody struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type TopUpRequestBody struct {
	DepositAmount float64 `json:"deposit_amount"`
}