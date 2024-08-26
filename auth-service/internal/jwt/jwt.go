package jwt

import (
	"github.com/golang-jwt/jwt"
	"go/mod/config"
	"time"
)

// Claims структура для данных, которые будут храниться в JWT токене
type Claims struct {
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip_address"`
	jwt.RegisteredClaims
}

// GenerateAccessToken создает новый JWT Access токен
func GenerateAccessToken(userID, ipAddress string) (string, error) {
	cfg := config.LoadConfig()

	// Настраиваем claims токена (включаем user_id, ip-адрес и другие стандартные данные)
	claims := Claims{
		UserID:    userID,
		IPAddress: ipAddress,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создаем токен с подписью HS512
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ValidateAccessToken проверяет корректность JWT токена и возвращает claims
func ValidateAccessToken(tokenString string) (*Claims, error) {
	cfg := config.LoadConfig()

	// Разбираем токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем, является ли токен действительным и корректным
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorClaimsInvalid)
}
