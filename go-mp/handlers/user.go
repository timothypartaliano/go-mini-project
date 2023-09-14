package handlers

import (
	"mini-project/model"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// @Summary Register a new user
// @Description Register a new user with the provided email and password
// @ID register-user
// @Accept json
// @Produce json
// @Param request body model.RegisterRequestBody true "User registration request body"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to hash password" "Failed to create user" "Failed to send registration email"
// @Router /register [post]
func RegisterUserHandler(c echo.Context) error {
    var requestBody model.RegisterRequestBody
    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
    }

    newUser := model.User{
        Email:    requestBody.Email,
        Password: string(hashedPassword),
    }

    if err := db.Create(&newUser).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
    }

    if err := sendRegistrationEmail(newUser.Email); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to send registration email"})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "User registered successfully"})
}

// @Summary Login
// @Description Login with the provided email and password to obtain an authentication token
// @ID login-user
// @Accept json
// @Produce json
// @Param request body model.RegisterRequestBody true "User login request body"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Failed to generate JWT token"
// @Router /login [post]
func LoginUserHandler(c echo.Context) error {
	var requestBody model.RegisterRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
	}

	var user model.User
	if err := db.Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.UserID
	claims["user"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate JWT token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
	})
}

// @Summary Top-Up User Account
// @Description Deposit a specified amount into the user's account balance
// @ID top-up-user
// @Accept json
// @Produce json
// @Param authorization header string true "JWT authorization token"
// @Param deposit_amount body model.TopUpRequestBody true "Amount to deposit"
// @Success 200 {object} map[string]interface{} "Top-up successful"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "JWT token missing or invalid"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Failed to perform top-up" "Failed to send top-up email"
// @Router /top-up [post]
func TopUpUserHandler(c echo.Context) error {
    userEmail := c.Get("user").(string)

    var requestBody model.TopUpRequestBody
    if err := c.Bind(&requestBody); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
    }

    var user model.User
    if err := db.Where("email = ?", userEmail).First(&user).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
    }

    user.DepositAmount += requestBody.DepositAmount

    if err := db.Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to perform top-up"})
    }

    if err := sendTopUpEmail(userEmail, requestBody.DepositAmount); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to send top-up email"})
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "Top-up successful",
        "user":    user,
    })
}

func sendRegistrationEmail(userEmail string) error {
    m := gomail.NewMessage()

    m.SetHeader("From", "tim@part.com")
    m.SetHeader("To", userEmail)

    m.SetHeader("Subject", "Registration Successful")
    m.SetBody("text/plain", "Thank you for registering with our service!")

    d := gomail.NewDialer("smtp-relay.brevo.com", 587, "timothypartaliano@gmail.com", "xsmtpsib-374b5461329392e030654874722cdd0efc42a63a024bb727deb020906b40f889-zDwMRJy37QNCh1O9")

    if err := d.DialAndSend(m); err != nil {
        return err
    }

    return nil
}

func sendTopUpEmail(userEmail string, depositAmount float64) error {
    m := gomail.NewMessage()

    m.SetHeader("From", "tim@part.com")
    m.SetHeader("To", userEmail)

    m.SetHeader("Subject", "Top-Up Successful")
    emailBody := fmt.Sprintf("Your account has been topped up successfully with $%.2f.", depositAmount)
    m.SetBody("text/plain", emailBody)

    d := gomail.NewDialer("smtp-relay.brevo.com", 587, "timothypartaliano@gmail.com", "xsmtpsib-374b5461329392e030654874722cdd0efc42a63a024bb727deb020906b40f889-zDwMRJy37QNCh1O9")

    if err := d.DialAndSend(m); err != nil {
        return err
    }

    return nil
}