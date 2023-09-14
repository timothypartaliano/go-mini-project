package handlers

import (
	"mini-project/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

// @Summary Create Equipment
// @Description Create a new equipment item
// @ID create-equipment
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Param request body model.CreateEquipmentRequestBody true "Equipment details"
// @Success 200 {string} string "Equipment created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to create equipment"
// @Router /equipment [post]
func CreateEquipmentHandler(c echo.Context) error {
	var requestBody model.CreateEquipmentRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	newEquipment := model.Equipment{
		Name:         requestBody.Name,
		Availability: requestBody.Availability,
		RentalCosts:  requestBody.RentalCosts,
		Category:     requestBody.Category,
	}

	if err := db.Create(&newEquipment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create equipment"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Equipment created successfully"})
}

// @Summary Get All Equipment
// @Description Retrieve a list of all available equipment
// @ID get-all-equipment
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Success 200 {array} model.Equipment "List of equipment"
// @Failure 401 {object} map[string]string "JWT token missing or invalid"
// @Failure 500 {object} map[string]string "Failed to retrieve equipment"
// @Router /equipment [get]
func GetAllEquipmentHandler(c echo.Context) error {
	var equipment []model.Equipment

	if err := db.Find(&equipment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve equipment"})
	}

	return c.JSON(http.StatusOK, equipment)
}

// @Summary Update Equipment
// @Description Update an existing equipment item by ID
// @ID update-equipment
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Param id path string true "Equipment ID"
// @Param request body model.UpdateEquipmentRequestBody true "Updated equipment details"
// @Success 200 {object} map[string]interface{} "Equipment updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Equipment not found"
// @Failure 500 {object} map[string]string "Failed to update equipment"
// @Router /equipment/{id} [put]
func UpdateEquipmentHandler(c echo.Context) error {
	var requestBody model.UpdateEquipmentRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	equipmentID := c.Param("id")

	var existingEquipment model.Equipment
	if err := db.First(&existingEquipment, equipmentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Equipment not found"})
	}

	if requestBody.Name != "" {
		existingEquipment.Name = requestBody.Name
	}
	existingEquipment.Availability = requestBody.Availability
	existingEquipment.RentalCosts = requestBody.RentalCosts
	if requestBody.Category != "" {
		existingEquipment.Category = requestBody.Category
	}

	if err := db.Save(&existingEquipment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update equipment"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Equipment updated successfully",
		"equipment": existingEquipment,
	})
}

// @Summary Delete Equipment
// @Description Delete an existing equipment item by ID
// @ID delete-equipment
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Param id path string true "Equipment ID"
// @Success 200 {object} map[string]string "Equipment deleted successfully"
// @Failure 404 {object} map[string]string "Equipment not found"
// @Failure 500 {object} map[string]string "Failed to delete equipment"
// @Router /equipment/{id} [delete]
func DeleteEquipmentHandler(c echo.Context) error {
	equipmentID := c.Param("id")

	var existingEquipment model.Equipment
	if err := db.First(&existingEquipment, equipmentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Equipment not found"})
	}

	if err := db.Delete(&existingEquipment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete equipment"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Equipment deleted successfully"})
}