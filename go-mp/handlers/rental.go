package handlers

import (
	"mini-project/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary Create Rental History
// @Description Create a new rental history record
// @ID create-rental-history
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Param request body model.CreateRentalHistoryRequestBody true "Request body containing rental history information"
// @Success 200 {object} map[string]interface{} "Rental history record created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "User or equipment not found"
// @Failure 409 {object} map[string]string "Equipment is not available for rent"
// @Failure 402 {object} map[string]string "Insufficient deposit amount"
// @Failure 500 {object} map[string]string "Failed to create rental history"
// @Router /rental [post]
func CreateRentalHistoryHandler(c echo.Context) error {
	var requestBody model.CreateRentalHistoryRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	var user model.User
	if err := db.First(&user, requestBody.UserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	var equipment model.Equipment
	if err := db.First(&equipment, requestBody.EquipmentID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Equipment not found"})
	}

	if !equipment.Availability {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Equipment is not available for rent"})
	}

	if user.DepositAmount < equipment.RentalCosts {
		return c.JSON(http.StatusPaymentRequired, map[string]string{"message": "Insufficient deposit amount"})
	}

	user.DepositAmount -= equipment.RentalCosts

	equipment.Availability = false

	newRentalHistory := model.RentalHistory{
		UserID:       requestBody.UserID,
		EquipmentID:  requestBody.EquipmentID,
		RentalDate:   requestBody.RentalDate,
		ReturnDate:   requestBody.ReturnDate,
		RentalStatus: requestBody.RentalStatus,
	}

	tx := db.Begin()

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create rental history"})
	}
	if err := tx.Save(&equipment).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create rental history"})
	}
	if err := tx.Create(&newRentalHistory).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create rental history"})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "Equipment rented successfully",
		"user_deposit_now": user.DepositAmount,
	})
}

// @Summary Get All Rental History
// @Description Get a list of all rental history records
// @ID get-all-rental-history
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Success 200 {array} model.RentalHistory "List of rental history records"
// @Failure 500 {object} map[string]string "Failed to retrieve rental history"
// @Router /rental [get]
func GetAllRentalHistoryHandler(c echo.Context) error {
	var rentalHistory []model.RentalHistory

	if err := db.Find(&rentalHistory).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve rental history"})
	}

	return c.JSON(http.StatusOK, rentalHistory)
}

// @Summary Update Rental History
// @Description Update an existing rental history record
// @ID update-rental-history
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Param id path int true "Rental history ID to be updated"
// @Param request body model.UpdateRentalHistoryRequestBody true "Request body containing updated rental history information"
// @Success 200 {object} map[string]interface{} "Rental history updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 404 {object} map[string]string "Rental history not found"
// @Failure 500 {object} map[string]string "Failed to update rental history"
// @Router /rental/{id} [put]
func UpdateRentalHistoryHandler(c echo.Context) error {
	var requestBody model.UpdateRentalHistoryRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	rentalHistoryID := c.Param("id")

	var existingRentalHistory model.RentalHistory
	if err := db.First(&existingRentalHistory, rentalHistoryID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Rental history not found"})
	}

	existingRentalHistory.UserID = requestBody.UserID
	existingRentalHistory.EquipmentID = requestBody.EquipmentID
	existingRentalHistory.RentalDate = requestBody.RentalDate
	existingRentalHistory.ReturnDate = requestBody.ReturnDate
	existingRentalHistory.RentalStatus = requestBody.RentalStatus

	if err := db.Save(&existingRentalHistory).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update rental history"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Rental history updated successfully",
		"data":    existingRentalHistory,
	})
}

// @Summary Delete Rental History
// @Description Delete a rental history record by ID
// @ID delete-rental-history
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Param id path int true "Rental history ID to be deleted"
// @Success 200 {object} map[string]string "Rental history deleted successfully"
// @Failure 404 {object} map[string]string "Rental history not found"
// @Failure 500 {object} map[string]string "Failed to delete rental history"
// @Router /rental/{id} [delete]
func DeleteRentalHistoryHandler(c echo.Context) error {
	rentalHistoryID := c.Param("id")

	var existingRentalHistory model.RentalHistory
	if err := db.First(&existingRentalHistory, rentalHistoryID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Rental history not found"})
	}

	if err := db.Delete(&existingRentalHistory).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete rental history"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Rental history deleted successfully"})
}