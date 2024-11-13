package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"social_api/schemas"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"social_api/db"
	"social_api/libs"
	"social_api/models"
)

var pool *pgxpool.Pool
var poolOnce sync.Once

type Claims struct {
	UserID string `json:"sub"`
	jwt.StandardClaims
}

var ErrTokenExpired = errors.New("Token has expired")

func (c Claims) Valid() error {
	if !c.VerifyExpiresAt(time.Now().Unix(), true) {
		return ErrTokenExpired
	}
	return nil
}

func IsUserTaken(email, username string) bool {
	query := "SELECT COUNT(*) FROM user_profile WHERE email = $1 OR username = $2"
	var count int

	err := pool.QueryRow(context.Background(), query, email, username).Scan(&count)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return true
	}

	return count > 0
}

func SaveUser(user models.User) (sql.NullString, error) {
	pool := db.Pool

	if db.Pool == nil {
		return sql.NullString{}, errors.New("Database pool is nil")
	}

	query := `
        INSERT INTO user_profile (username, firstname, lastname, email, password, picture)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	var newID sql.NullString
	err := pool.QueryRow(context.Background(), query, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Picture).Scan(&newID)
	if err != nil {
		cantRegisterError := errors.New("The user could not register. Please contact support or try again.")
		errorJSON, _ := json.Marshal(err)

		var errorMap map[string]interface{}
		if err := json.Unmarshal(errorJSON, &errorMap); err == nil {
			errorCode, ok := errorMap["Code"].(string)
			if ok {
				if errorCode == "42P05" || errorCode == "23505" {
					errConstraintName, ok := errorMap["Message"].(string)
					if ok {
						if strings.Contains(errConstraintName, "email") {
							return sql.NullString{}, errors.New("The user already exists. Please try another email.")
						} else if strings.Contains(errConstraintName, "username") {
							return sql.NullString{}, errors.New("The user already exists. Please try another username.")
						} else {
							return sql.NullString{}, errors.New("The user already exists. Please try another username or email.")
						}
					}
				} else {
					return sql.NullString{}, cantRegisterError
				}
			} else {
				return sql.NullString{}, cantRegisterError
			}
		}

		return sql.NullString{}, cantRegisterError
	}

	return newID, nil
}

func FindUserByEmailOrUsername(email, username string) (*models.User, error) {
	pool := db.Pool

	var user models.User
	query := "SELECT Id, Username, Firstname, Lastname, Email, Password, Picture FROM user_profile WHERE Email = $1 OR Username = $2"
	row := pool.QueryRow(context.Background(), query, email, username)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Picture)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func FindUserById(userID string) (*models.User, error) {
	pool := db.Pool

	var user models.User
	query := "SELECT ID, Username, FirstName, LastName, Email, Password, Picture, Description FROM user_profile WHERE ID = $1"

	var picture, description sql.NullString

	row := pool.QueryRow(context.Background(), query, userID)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Picture, &user.Description)

	if err != nil {
		return nil, err
	}

	if picture.Valid {
		user.Picture = picture
	} else {
		user.Picture = sql.NullString{String: "", Valid: false}
	}

	if description.Valid {
		user.Description = description
	} else {
		user.Description = sql.NullString{String: "", Valid: false}
	}

	return &user, nil
}

func ValidateRequiredFields(c *fiber.Ctx, requestBody interface{}, requiredFields []string) ([]string, bool) {
	missingFields := make([]string, 0)

	val := reflect.ValueOf(requestBody)

	for _, field := range requiredFields {
		fieldValue := val.FieldByName(field).String()

		if fieldValue == "" {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		HandleError(c, ErrMissingFields, http.StatusBadRequest)
		return missingFields, false
	}

	return missingFields, true
}

func UserWithoutPassword(user models.User, ID string) map[string]interface{} {
	return map[string]interface{}{
		"ID":        ID,
		"username":  user.Username,
		"email":     user.Email,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
		"picture":   user.Picture,
	}
}

func UserWithoutPasswordAndEmail(user models.User, ID string) map[string]interface{} {
	return map[string]interface{}{
		"ID":       ID,
		"username": user.Username,
		"picture":  user.Picture,
	}
}

func ExtractUserIDFromToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &libs.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil {
		return "", err
	}

	userID := token.Claims.(*libs.CustomClaims).UserID
	if userID == "" {
		return "", errors.New("UserID is empty")
	}

	return userID, nil
}

func UpdateUsername(userID string, newUsername string) error {
	pool := db.Pool

	query := "UPDATE user_profile SET Username = $1 WHERE Id = $2"
	_, err := pool.Exec(context.Background(), query, newUsername, userID)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePassword(userID string, newPassword string) error {
	pool := db.Pool

	query := "UPDATE user_profile SET password = $1 WHERE id = $2"
	_, err := pool.Exec(context.Background(), query, newPassword, userID)
	if err != nil {
		fmt.Println("Error updating password:", err)
		return err
	}

	return nil
}

func HasExtraFields(requestBody schemas.LoginRequest, allowedFields []string) bool {
	fieldCount := 0

	for _, field := range allowedFields {
		switch field {
		case "Email":
			if requestBody.Email != "" {
				fieldCount++
			}
		case "Username":
			if requestBody.Username != "" {
				fieldCount++
			}
		case "Password":
			if requestBody.Password != "" {
				fieldCount++
			}
		}
	}

	return fieldCount != len(allowedFields)
}

func FollowUser(follower_id string, followed_id string) error {
	pool := db.Pool

	query := "INSERT INTO followers (follower_id, following_id) VALUES ($1, $2)"
	_, err := pool.Exec(context.Background(), query, follower_id, followed_id)
	if err != nil {
		return err
	}

	return nil
}

func UnFollowUser(follower_id string, followed_id string) error {
	pool := db.Pool

	query := "DELETE FROM followers WHERE follower_id = $1 AND following_id = $2;"

	_, err := pool.Exec(context.Background(), query, follower_id, followed_id)
	if err != nil {
		return err
	}

	return nil
}

func IsFollowing(followerID string, followedID string) (bool, error) {
	pool := db.Pool

	query := "SELECT COUNT(*) FROM followers WHERE follower_id = $1 AND following_id = $2"
	var count int
	err := pool.QueryRow(context.Background(), query, followerID, followedID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func IsNotFollowing(followerID string, followedID string) (bool, error) {
	pool := db.Pool

	query := "SELECT COUNT(*) FROM followers WHERE follower_id = $1 AND following_id = $2"
	var count int
	err := pool.QueryRow(context.Background(), query, followerID, followedID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func ParseToken(tokenString string) (string, error) {
	var jwtKey = os.Getenv("SECRET_JWT")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("signature is invalid")
		}
		return "", errors.New("could not parse token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

func GetFollowers(uuid string) ([]models.UserRelevantInfo, error) {
	pool := db.Pool

	// Consulta para obtener solo los ID y username de los seguidores
	query := `
        SELECT u.id, u.username
        FROM followers f
        JOIN user_profile u ON f.follower_id = u.id
        WHERE f.following_id = $1
    `

	rows, err := pool.Query(context.Background(), query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.UserRelevantInfo
	for rows.Next() {
		var user models.UserRelevantInfo
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}

	if len(followers) == 0 {
		return nil, errors.New("No followers found")
	}

	return followers, nil
}
