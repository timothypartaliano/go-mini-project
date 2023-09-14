package model

type Equipment struct {
    EquipmentID   uint    `gorm:"primaryKey"`
    Name          string  `gorm:"not null"`
    Availability  bool    `gorm:"not null"`
    RentalCosts   float64 `gorm:"not null"`
    Category      string  `gorm:"not null"`
}

type CreateEquipmentRequestBody struct {
    Name         string  `json:"name"`
    Availability bool    `json:"availability"`
    RentalCosts  float64 `json:"rental_costs"`
    Category     string  `json:"category"`
}

type UpdateEquipmentRequestBody struct {
    Name         string  `json:"name"`
    Availability bool    `json:"availability"`
    RentalCosts  float64 `json:"rental_costs"`
    Category     string  `json:"category"`
}