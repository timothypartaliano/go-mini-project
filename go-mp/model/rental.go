package model

type RentalHistory struct {
    RentalHistoryID uint   `gorm:"primaryKey"`
    UserID          uint   `gorm:"not null"`
    EquipmentID     uint   `gorm:"not null"`
    RentalDate      string `gorm:"not null"`
    ReturnDate      string
    RentalStatus    string `gorm:"not null"`
}

type CreateRentalHistoryRequestBody struct {
    UserID       uint   `json:"user_id"`
    EquipmentID  uint   `json:"equipment_id"`
    RentalDate   string `json:"rental_date"`
    ReturnDate   string `json:"return_date"`
    RentalStatus string `json:"rental_status"`
}

type UpdateRentalHistoryRequestBody struct {
    UserID       uint   `json:"user_id"`
    EquipmentID  uint   `json:"equipment_id"`
    RentalDate   string `json:"rental_date"`
    ReturnDate   string `json:"return_date"`
    RentalStatus string `json:"rental_status"`
}