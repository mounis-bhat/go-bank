package lib

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/mounis-bhat/go-bank/types"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWTToken(account *types.Account) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": account.ID,
		"first_name": account.FirstName,
		"last_name":  account.LastName,
		"balance":    account.Balance,
		"created_at": account.CreatedAt,
		"username":   account.Username,
		"roles":      strings.Join(account.Roles, ","),
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func ValidateToken(accessToken string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["exp"] == nil {
		return nil, fmt.Errorf("invalid token")
	}

	exp := int64(claims["exp"].(float64))
	if time.Now().Unix() > exp {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func GetAccountAndValidate(r *http.Request) (*types.ValidateAccountRequest, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	accessToken := strings.Split(token, " ")[1]
	claims, err := ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	validatedAccount := &types.ValidateAccountRequest{
		ID:        int(claims["account_id"].(float64)),
		FirstName: claims["first_name"].(string),
		LastName:  claims["last_name"].(string),
		Balance:   int(claims["balance"].(float64)),
		CreatedAt: claims["created_at"].(string),
		Username:  claims["username"].(string),
		Roles:     claims["roles"].(string),
		IAT:       int(claims["iat"].(float64)),
		EXP:       int(claims["exp"].(float64)),
	}

	return validatedAccount, nil
}

func HashAndSaltPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePasswords(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}

	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}

	if !strings.ContainsAny(password, "0123456789") {
		return false
	}

	if !strings.ContainsAny(password, "!@#$%^&*()") {
		return false
	}

	return true
}
