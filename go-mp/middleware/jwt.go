package middleware

import (
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v5"
    "os"
    "net/http"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        tokenString := c.Request().Header.Get("Authorization")
        if tokenString == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token credentials"})
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
            return secretKey, nil
        })

        if err != nil || !token.Valid {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token credentials"})
        }

        userClaim, ok := token.Claims.(jwt.MapClaims)["user"].(string)
        if !ok {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token credentials"})
        }

        c.Set("user", userClaim)

        return next(c)
    }
}